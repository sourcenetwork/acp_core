import { LocatedMessage } from "@/types/proto-js/sourcenetwork/acp_core/parser_message";
import * as monaco from "monaco-editor";

export function mapLocatedMessageMarkers(
  messages: LocatedMessage[]
): monaco.editor.IMarkerData[] {
  return messages.map((msg) => ({
    message: msg.message,
    startLineNumber: Number(msg.interval?.start?.line ?? 0),
    endLineNumber: Number(msg.interval?.end?.line ?? 0),
    startColumn: Number(msg.interval?.start?.column ?? 0),
    endColumn: Number(msg.interval?.end?.column ?? 0),
    severity: monaco.MarkerSeverity.Error,
  }));
}
