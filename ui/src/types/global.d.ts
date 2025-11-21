import { PlaygroundService } from "./proto-js/sourcenetwork/acp_core/playground";

declare global {
  export interface Window {
    Go: {
      // eslint-disable-next-line
      new (): {
        run: (inst: WebAssembly.Instance) => Promise<void>;
        importObject: WebAssembly.Imports;
      };
    };
    AcpPlayground: { new: () => Promise<PlaygroundService> };
  }
}

export {};
