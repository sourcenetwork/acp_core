import { PlaygroundService } from "@/types/proto-js/sourcenetwork/acp_core/playground.js";
import "../lib/wasm_exec.js";

declare global {
  const Go: {
    new (): {
      run: (inst: WebAssembly.Instance) => Promise<void>;
      importObject: WebAssembly.Imports;
    };
  };

  const AcpPlayground: { new: () => Promise<PlaygroundService> };
}

let playground: PlaygroundService | null = null;
let initialized = false;

const importMeta = import.meta.env;
const WASM_PATH = importMeta.VITE_WASM_PATH ?? "/playground.wasm";

type MessageReq = { id: string; method: string; payload?: unknown };
type MessageRes =
  | { id: string; success: true; result: unknown }
  | {
      id: string;
      success: false;
      error: { message: string; name: string; stack?: string };
    };

async function loadAndInitialize() {
  if (initialized) return;

  const go = new Go();
  const result = await WebAssembly.instantiateStreaming(
    fetch(WASM_PATH),
    go.importObject
  );

  go.run(result.instance);
  playground = await self.AcpPlayground.new();
  initialized = true;
}

const clonePayloadForPost = (r: unknown) => {
  try {
    return structuredClone(r);
  } catch {
    if (r == null) return null;
    if (typeof r === "function") return null;
    try {
      return JSON.parse(JSON.stringify(r));
    } catch {
      return null;
    }
  }
};

self.addEventListener("message", async (ev: MessageEvent<MessageReq>) => {
  const { id, method, payload } = ev.data;

  try {
    if (method === "__init__") {
      await loadAndInitialize();

      self.postMessage({
        id,
        success: true,
        result: { initialized: true },
      } as MessageRes);
      return;
    }

    if (!initialized || !playground) await loadAndInitialize();

    if (!playground) throw new Error(`Playground failed to initialize`);

    const playgroundFn = (playground as any)[method];

    if (typeof playgroundFn !== "function")
      throw new Error(`Unknown method: ${method}`);

    const result = await playgroundFn(payload);

    self.postMessage({
      id,
      success: true,
      result: clonePayloadForPost(result),
    } as MessageRes);
  } catch (e) {
    const err = e as Error;

    self.postMessage({
      id,
      success: false,
      error: { message: err.message, name: err.name, stack: err.stack },
    } as MessageRes);
  }
});
