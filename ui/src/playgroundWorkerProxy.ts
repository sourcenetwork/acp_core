import { PlaygroundService } from "@/types/proto-js/sourcenetwork/acp_core/playground";

type MessageRes =
  | { id: string; success: true; result: unknown }
  | {
      id: string;
      success: false;
      error: { message: string; name: string; stack?: string };
    };

const clonePayloadForPost = (v: unknown) => {
  try {
    return structuredClone(v);
  } catch {
    if (v == null) return null;
    if (typeof v === "function") return null;
    try {
      return JSON.parse(JSON.stringify(v));
    } catch {
      return null;
    }
  }
};

export class PlaygroundWorkerClient {
  private id = 0;
  private worker: Worker | null = null;
  private pendingMessages = new Map<
    string,
    {
      resolve: (value: unknown) => void;
      reject: (e: any) => void;
    }
  >();
  private proxy: PlaygroundService | null = null;

  async initialize(): Promise<void> {
    if (this.worker) return;

    this.worker = new Worker(
      new URL("./workers/playground.worker.ts", import.meta.url),
      { type: "module" }
    );

    this.worker.addEventListener("message", (ev: MessageEvent<MessageRes>) => {
      const { id, success: ok } = ev.data;
      const p = this.pendingMessages.get(id);
      if (!p) return;

      this.pendingMessages.delete(id);
      if (ok) p.resolve(ev.data.result);
      else p.reject(ev.data.error);
    });

    await this.call("__init__", null);
  }

  private call<T = unknown>(method: string, payload: unknown): Promise<T> {
    if (!this.worker) throw new Error("Worker not initialized");

    const id = `w_${++this.id}`;
    const safePayload = clonePayloadForPost(payload);

    return new Promise<T>((resolve, reject) => {
      this.pendingMessages.set(id, {
        resolve: (v: unknown) => resolve(v as T),
        reject: (e: any) => reject(e),
      });

      this.worker!.postMessage({ id, method, payload: safePayload });

      setTimeout(() => {
        if (this.pendingMessages.has(id)) {
          this.pendingMessages.delete(id);
          reject(new Error(`Worker call timeout: ${method}`));
        }
      }, 30000);
    });
  }

  async getPlaygroundProxy(): Promise<PlaygroundService> {
    if (!this.worker) await this.initialize();
    if (!this.proxy) {
      this.proxy = new Proxy(
        {},
        {
          get: (_t, prop) => {
            if (prop === "then") return undefined; // prevent thenable detection
            return (payload: unknown) => this.call(String(prop), payload);
          },
        }
      ) as PlaygroundService;
    }
    return this.proxy;
  }

  terminate(): void {
    this.pendingMessages.forEach(({ reject }) =>
      reject(new Error("Worker terminated"))
    );
    this.pendingMessages.clear();
    this.worker?.terminate();
    this.worker = null;
    this.proxy = null;
  }
}

let playgroundWorker: PlaygroundWorkerClient | null = null;

export async function getPlaygroundWorkerProxy(): Promise<PlaygroundWorkerClient> {
  if (playgroundWorker) return playgroundWorker;

  playgroundWorker = new PlaygroundWorkerClient();
  await playgroundWorker.initialize();

  return playgroundWorker;
}
