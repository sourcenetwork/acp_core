import type { CheckRequest, CheckResponse, ExplainRequest, ExplainResponse, GenerateRequest, GenerateResponse, VerifyRequest, VerifyResponse } from "./types";

const SERVICE = "/sourcenetwork.acp_core.Playground2";
export const GENERATE_PATH = `${SERVICE}/Generate`;
export const CHECK_PATH = `${SERVICE}/Check`;
export const VERIFY_PATH = `${SERVICE}/Verify`;
export const EXPLAIN_PATH = `${SERVICE}/Explain`;

export type ApiErrorKind = "backend" | "network" | "protocol";

export class ApiError extends Error {
  kind: ApiErrorKind;
  code?: string;
  status?: number;
  constructor(kind: ApiErrorKind, message: string, opts: { code?: string; status?: number } = {}) {
    super(message);
    this.kind = kind;
    this.code = opts.code;
    this.status = opts.status;
  }
}

// Inside `acp playground` the backend is acp, which proxies the service: say so in errors.
let fileBacked = false;
export const setFileBacked = (on: boolean) => {
  fileBacked = on;
};

async function post<T>(path: string, req: unknown, signal?: AbortSignal): Promise<T> {
  let res: Response;
  try {
    res = await fetch(path, {
      method: "POST",
      headers: { "Content-Type": "application/json" },
      body: JSON.stringify(req),
      signal,
    });
  } catch (e) {
    if (e instanceof DOMException && e.name === "AbortError") throw e;
    throw new ApiError(
      "network",
      fileBacked
        ? `Could not reach acp playground (${(e as Error).message}). Is it still running?`
        : `Could not reach the backend (${(e as Error).message}). Is it running on localhost:8090?`,
    );
  }

  const raw = await res.text();
  let body: unknown;
  try {
    body = raw === "" && res.ok ? {} : JSON.parse(raw); // an empty error body is a proxy with nothing behind it
  } catch {
    // Dev proxy answers 5xx with plain text when the backend is down.
    const kind = res.ok ? "protocol" : "network";
    throw new ApiError(kind, res.ok ? "Backend returned invalid JSON" : `Backend unavailable (HTTP ${res.status}): ${raw.slice(0, 200) || res.statusText}`, { status: res.status });
  }

  if (!res.ok) {
    const b = body as { code?: string; message?: string };
    // acp's proxy answers 502 {"code":"unavailable","message":"can't reach the acp service at ..."}: show it as is
    if (res.status === 502 && b.code === "unavailable") throw new ApiError("network", b.message ?? "the acp service is unavailable", { code: b.code, status: res.status });
    throw new ApiError("backend", b.message ?? `HTTP ${res.status}`, { code: b.code, status: res.status });
  }
  return body as T;
}

export const generate = (req: GenerateRequest, signal?: AbortSignal) => post<GenerateResponse>(GENERATE_PATH, req, signal);

/** `allowed` is omitted by the wire format when false: use `res.allowed === true`. */
export const check = (req: CheckRequest, signal?: AbortSignal) => post<CheckResponse>(CHECK_PATH, req, signal);

export const explain = (req: ExplainRequest, signal?: AbortSignal) => post<ExplainResponse>(EXPLAIN_PATH, req, signal);

export const verify = (req: VerifyRequest, signal?: AbortSignal) => post<VerifyResponse>(VERIFY_PATH, req, signal);
