import {
  AnnotatedAuthorizationTheoremResult,
  AnnotatedDelegationTheoremResult,
  AnnotatedPolicyTheoremResult,
  AnnotatedReachabilityTheoremResult,
} from "@acp/theorem";
import * as monaco from "monaco-editor";

type TheoremResultType =
  | AnnotatedDelegationTheoremResult
  | AnnotatedAuthorizationTheoremResult
  | AnnotatedReachabilityTheoremResult;

interface CategorizedTheoremMarkers {
  accepted: monaco.editor.IMarkerData[];
  rejected: monaco.editor.IMarkerData[];
  errors: monaco.editor.IMarkerData[];
}

function categorizeTheoremMarkers(
  theoremResult: TheoremResultType[] = []
): CategorizedTheoremMarkers {
  const accepted: monaco.editor.IMarkerData[] = [];
  const rejected: monaco.editor.IMarkerData[] = [];
  const errors: monaco.editor.IMarkerData[] = [];

  theoremResult.forEach((result) => {
    const marker = {
      message: `${result.result?.result?.status?.toString() ?? "Error"}`,
      startLineNumber: Number(result.interval?.start?.line ?? 0),
      endLineNumber: Number(result.interval?.end?.line ?? 0),
      startColumn: Number(result.interval?.start?.column ?? 0),
      endColumn: Number.MAX_SAFE_INTEGER,
      severity: monaco.MarkerSeverity.Info,
    };

    const status = String(result.result?.result?.status);

    switch (status) {
      case "Accept":
        accepted.push({ ...marker, severity: monaco.MarkerSeverity.Info });
        break;
      case "Reject":
        rejected.push({ ...marker, severity: monaco.MarkerSeverity.Error });
        break;
      default:
        errors.push({ ...marker, severity: monaco.MarkerSeverity.Error });
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
