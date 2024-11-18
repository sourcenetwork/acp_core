import { PlaygroundService } from "./proto-js/sourcenetwork/acp_core/playground";

declare global {
  export interface Window {
    // eslint-disable-next-line @typescript-eslint/no-explicit-any
    Go: any;
    AcpPlayground: { new: () => Promise<PlaygroundService> };
  }
}

export {};
