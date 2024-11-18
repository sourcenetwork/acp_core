// import { SandboxRecord } from "@/types/acpPlayground";
import {
  NewSandboxRequest,
  PlaygroundService,
} from "@/types/proto-js/sourcenetwork/acp_core/playground";
import {
  SandboxData,
  SandboxDataErrors,
  SandboxRecord,
} from "@/types/proto-js/sourcenetwork/acp_core/sandbox";
import { AnnotatedPolicyTheoremResult } from "@/types/proto-js/sourcenetwork/acp_core/theorem";
import { WASM_PATH } from "@/utils/constants";
import { loadPlaygroundWasm } from "@/utils/loadWasm";
import { theoremResultPassing } from "@/utils/mapTheoremResultMarkers";
import { create } from "zustand";
import { subscribeWithSelector } from "zustand/middleware";
import { usePlaygroundStorageStore } from "./acpStorage";

export interface PlaygroundState {
  status: "uninitialized" | "loading" | "ready" | "error";
  error?: string;
  module?: WebAssembly.Instance;
  playground?: PlaygroundService | null;
  active?: SandboxRecord;
  sandboxErrors?: SandboxDataErrors;
  sandboxErrorCount: number;
  annotatedPolicyTheoremResult?: AnnotatedPolicyTheoremResult;
  verifyTheoremsStatus: "pending" | "loading" | "passed" | "error";
  setStateError?: unknown;

  initialize: () => Promise<void>;
  newSandbox: (args: NewSandboxRequest) => Promise<void>;
  setState: (args: Partial<SandboxData>) => Promise<void>;
  verifyTheorems: () => Promise<void>;
}

export const usePlaygroundStore = create<PlaygroundState>()(
  subscribeWithSelector((set, get) => {
    const storage = usePlaygroundStorageStore.getState;
    const { updateStore } = storage();

    const getActiveSandbox = () => {
      const { playground, active } = get();
      if (!playground || !active)
        throw new Error("Playground or sandbox is not initialized");
      return { playground, active, handle: active.handle };
    };

    return {
      status: "uninitialized",
      verifyTheoremsStatus: "pending",
      sandboxErrorCount: 0,
      initialize: async () => {
        try {
          set({ status: "loading" });

          const module = await loadPlaygroundWasm(WASM_PATH);
          const playground = await window.AcpPlayground.new();

          set({
            status: "ready",
            module: module.instance,
            playground: playground,
          });

          const { setState, newSandbox } = get();

          // Initialize a default sandbox
          await newSandbox({ name: "playground-1", description: "" });
          // Set the inital state of the sandbox
          await setState({});
        } catch (error) {
          set({ status: "error", error: (error as Error)?.message });
        }
      },

      newSandbox: async (args) => {
        const { playground } = get();
        if (!playground) throw new Error("Playground is not initialized");

        // Initialize a new sandbox
        const sandbox = await playground.NewSandbox(args);
        set({ active: sandbox.record });
      },

      setState: async (updates) => {
        try {
          const { playground, handle } = getActiveSandbox();

          set({ setStateError: null });

          const state = storage().data || {};
          const newState = { ...state, ...updates };

          const stateResult = await playground?.SetState({
            handle,
            data: newState,
          });

          set({
            sandboxErrors: stateResult?.errors,
            sandboxErrorCount: Object.values(stateResult?.errors ?? {}).reduce(
              (count, err) => count + err?.length,
              0
            ),
            verifyTheoremsStatus: "pending",
            annotatedPolicyTheoremResult: undefined,
          });

          // Persist state data into storage
          updateStore(newState);
        } catch (error) {
          console.error(error, "Failed to run playground SetState");
          set({ setStateError: error });
        }
      },

      verifyTheorems: async () => {
        try {
          const { playground, handle } = getActiveSandbox();

          set({
            annotatedPolicyTheoremResult: undefined,
            verifyTheoremsStatus: "pending",
          });

          const { result } = await playground.VerifyTheorems({ handle });

          const validationPassing = result
            ? theoremResultPassing(result)
            : true;

          set({
            verifyTheoremsStatus: validationPassing ? "passed" : "error",
            annotatedPolicyTheoremResult: result,
          });
        } catch (error) {
          console.error(error, "Failed to run playground VerifyTheorems");
        }
      },
    };
  })
);
