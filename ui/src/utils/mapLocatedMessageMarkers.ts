import { LocatedMessage } from "@acp/parser_message";
import type { editor } from "monaco-editor";

export function mapLocatedMessageMarkers(
  messages: LocatedMessage[]
): editor.IMarkerData[] {
  return messages.map((msg) => ({
    message: msg.message,
    startLineNumber: Number(msg.interval?.start?.line ?? 0),
    endLineNumber: Number(msg.interval?.end?.line ?? 0),
    startColumn: Number(msg.interval?.start?.column ?? 0),
    endColumn: Number(msg.interval?.end?.column ?? 0),
    severity: 8, // monaco.MarkerSeverity.Error
  }));
}
