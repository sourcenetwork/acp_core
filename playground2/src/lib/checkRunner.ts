// Pure async helpers for running one Check per actor with bounded concurrency.

export type ActorStatus =
  | { state: "loading" }
  | { state: "allowed" }
  | { state: "denied" }
  | { state: "error"; message: string };

/**
 * Runs `fn` over `items`, at most `limit` at a time, reporting each outcome as it finishes.
 * Stops starting new work once `signal` aborts; outcomes after an abort are not reported.
 */
export async function runLimited<T>(
  items: T[],
  limit: number,
  fn: (item: T, signal: AbortSignal) => Promise<boolean>,
  onResult: (item: T, status: ActorStatus) => void,
  signal: AbortSignal,
): Promise<void> {
  let next = 0;
  const worker = async () => {
    while (!signal.aborted && next < items.length) {
      const item = items[next++];
      onResult(item, { state: "loading" });
      let status: ActorStatus;
      try {
        status = (await fn(item, signal)) ? { state: "allowed" } : { state: "denied" };
      } catch (e) {
        status = { state: "error", message: e instanceof Error ? e.message : String(e) };
      }
      if (signal.aborted) return;
      onResult(item, status);
    }
  };
  await Promise.all(Array.from({ length: Math.max(1, Math.min(limit, items.length)) }, worker));
}

/** The backend's error messages end with `; attrs={...}; kind=...` debugging detail: keep the readable head. */
export function shortError(message: string): string {
  const i = message.indexOf("; attrs=");
  return (i < 0 ? message : message.slice(0, i)).trim();
}

export interface CheckSummary {
  allowed: string[];
  denied: string[];
  pending: number;
  errors: { actor: string; message: string }[];
}

export function summarize(actors: string[], statuses: Record<string, ActorStatus | undefined>): CheckSummary {
  const s: CheckSummary = { allowed: [], denied: [], pending: 0, errors: [] };
  for (const a of actors) {
    const st = statuses[a];
    if (!st || st.state === "loading") s.pending++;
    else if (st.state === "allowed") s.allowed.push(a);
    else if (st.state === "denied") s.denied.push(a);
    else s.errors.push({ actor: a, message: st.message });
  }
  return s;
}

/** `did:user:alice` -> `alice` for compact summaries. */
export const shortActor = (did: string) => (did.startsWith("did:user:") ? did.slice("did:user:".length) : did);
