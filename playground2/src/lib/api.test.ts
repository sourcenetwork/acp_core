import { afterEach, describe, expect, it, vi } from "vitest";
import { generate, setFileBacked } from "./api";

const req = { systemModel: { policy: "p", constraints: [] } };
afterEach(() => {
  vi.unstubAllGlobals();
  setFileBacked(false);
});

describe("api errors", () => {
  it("hints at localhost:8090 in local mode and at acp in file-backed mode", async () => {
    vi.stubGlobal("fetch", async () => { throw new TypeError("failed"); });
    await expect(generate(req)).rejects.toMatchObject({ kind: "network", message: expect.stringContaining("localhost:8090") });
    setFileBacked(true);
    await expect(generate(req)).rejects.toMatchObject({ kind: "network", message: expect.stringContaining("acp playground") });
  });
  it("shows the proxy's 502 message as a network error", async () => {
    const msg = "can't reach the acp service at http://x: connection refused";
    vi.stubGlobal("fetch", async () => new Response(JSON.stringify({ code: "unavailable", message: msg }), { status: 502 }));
    await expect(generate(req)).rejects.toMatchObject({ kind: "network", message: msg, code: "unavailable" });
  });
  it("other errors stay backend errors", async () => {
    vi.stubGlobal("fetch", async () => new Response(JSON.stringify({ code: "invalid_argument", message: "bad" }), { status: 400 }));
    await expect(generate(req)).rejects.toMatchObject({ kind: "backend", message: "bad" });
  });
});
