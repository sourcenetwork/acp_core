import {
  AnnotatedDelegationTheoremResult,
  AnnotatedPolicyTheoremResult,
  AnnotatedReachabilityTheoremResult,
  ResultStatus,
} from "@/types/proto-js/sourcenetwork/acp_core/theorem";
import * as monaco from "monaco-editor";

type TheoremResultType =
  | AnnotatedDelegationTheoremResult
  | AnnotatedDelegationTheoremResult
  | AnnotatedReachabilityTheoremResult;

function mapTheoremMarkers(
  theoremResult: TheoremResultType[]
): monaco.editor.IMarkerData[] {
  return theoremResult
    .filter(
      ({ result }) =>
        result?.result?.status !== "Accept" // FIXME
    )
    .map((result) => {
      return {
        message: `${result.result?.result?.status?.toString() ?? "Error"}`,
        startLineNumber: Number(result.interval?.start?.line ?? 0),
        endLineNumber: Number(result.interval?.end?.line ?? 0),
        startColumn: Number(result.interval?.start?.column ?? 0),
        endColumn: Number(result.interval?.end?.column ?? 0),
        severity: monaco.MarkerSeverity.Error,
      };
    });
}

export function mapTheoremResultMarkers(result?: AnnotatedPolicyTheoremResult) {
  const { delegationTheoremsResult, authorizationTheoremsResult } =
    result ?? {};

  const authMarkers = mapTheoremMarkers(authorizationTheoremsResult ?? []);
  const delegationMarkers = mapTheoremMarkers(delegationTheoremsResult ?? []);

  return [...authMarkers, ...delegationMarkers];
}

export function theoremResultPassing(result: AnnotatedPolicyTheoremResult) {
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
      (set) =>
        resultStatusFromJSON(set.result?.result?.status) !== ResultStatus.Accept
    )
  );

  return !hasErrors;
}
