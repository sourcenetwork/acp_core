import {
  AnnotatedAuthorizationTheoremResult,
  AnnotatedDelegationTheoremResult,
  AnnotatedPolicyTheoremResult,
  AnnotatedReachabilityTheoremResult,
} from "@acp/theorem";
import type { editor } from "monaco-editor";

type TheoremResultType =
  | AnnotatedDelegationTheoremResult
  | AnnotatedAuthorizationTheoremResult
  | AnnotatedReachabilityTheoremResult;

interface CategorizedTheoremMarkers {
  accepted: editor.IMarkerData[];
  rejected: editor.IMarkerData[];
  errors: editor.IMarkerData[];
}

function categorizeTheoremMarkers(
  theoremResult: TheoremResultType[] = []
): CategorizedTheoremMarkers {
  const accepted: editor.IMarkerData[] = [];
  const rejected: editor.IMarkerData[] = [];
  const errors: editor.IMarkerData[] = [];

  theoremResult.forEach((result) => {
    const marker = {
      message: `${result.result?.result?.status?.toString() ?? "Error"}`,
      startLineNumber: Number(result.interval?.start?.line ?? 0),
      endLineNumber: Number(result.interval?.end?.line ?? 0),
      startColumn: Number(result.interval?.start?.column ?? 0),
      endColumn: Number.MAX_SAFE_INTEGER,
      severity: 2, //monaco.MarkerSeverity.Info
    };

    const status = String(result.result?.result?.status);

    switch (status) {
      case "Accept":
        accepted.push({ ...marker, severity: 2 });
        break;
      case "Reject":
        rejected.push({ ...marker, severity: 8 });
        break;
      default:
        errors.push({ ...marker, severity: 8 });
        break;
    }
  });

  return { accepted, rejected, errors };
}

export function mapTheoremResultMarkers(result?: AnnotatedPolicyTheoremResult) {
  if (!result) return;
  const { delegationTheoremsResult, authorizationTheoremsResult } = result;

  const authMarkers = categorizeTheoremMarkers(authorizationTheoremsResult);
  const delegationMarkers = categorizeTheoremMarkers(delegationTheoremsResult);

  return { authMarkers, delegationMarkers };
}

export function theoremResultPassing(result?: AnnotatedPolicyTheoremResult) {
  if (!result) return false;

  const {
    delegationTheoremsResult,
    authorizationTheoremsResult,
    reachabilityTheoremsResult,
  } = result;

  const hasErrors = [
    delegationTheoremsResult,
    authorizationTheoremsResult,
    reachabilityTheoremsResult,
  ]?.some((type) =>
    type?.some(
      (set) => (set?.result?.result?.status as unknown as string) !== "Accept" // FIXME
    )
  );

  return !hasErrors;
}
