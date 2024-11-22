import { STORE_KEY, STORE_VERSION, WASM_PATH } from "@/utils/constants";
import { countErrors } from "@/utils/errorUtils";
import { loadPlaygroundWasm } from "@/utils/loadWasm";
import { theoremResultPassing } from "@/utils/mapTheoremResultMarkers";
import { NewSandboxRequest, PlaygroundService } from "@acp/playground";
import { SandboxData, SandboxDataErrors, SandboxTemplate } from "@acp/sandbox";
import { AnnotatedPolicyTheoremResult } from "@acp/theorem";
import { create } from "zustand";
import {
  createJSONStorage,
  persist,
  subscribeWithSelector,
} from "zustand/middleware";
import { initialStates } from "./initialStoreState";
import {
  getActiveSandboxHandle,
  getLastActiveSandbox,
} from "./playgroundUtils";

export const blankSandbox = {
  policyDefinition: `name: ""`,
  policyTheorem: `Authorizations {\n\n}\n\nDelegations {\n}`,
  relationships: "",
};

export interface PersistedSandboxData {
  createdAt?: Date;
  updatedAt?: Date;
  id: string;
  name: string;
  description: string;
  data: SandboxData;
}

export interface PlaygroundState {
  /* Pesisted State */
  lastActiveId: string | null;
  sandboxes: PersistedSandboxData[];

  /* Ephemeral State */
  idHandleMap: Record<string, number>;
  activeHandle?: number;

  /* Playground Data */
  playgroundStatus: "uninitialized" | "loading" | "ready" | "error";
  playgroundError?: string;
  playground?: PlaygroundService | null;
  setStateDataErrors?: SandboxDataErrors;
  setStateDataErrorCount: number;
  setStateError?: string;
  verifyTheoremsStatus: "pending" | "loading" | "passed" | "error";
  verifyTheoremsResult?: AnnotatedPolicyTheoremResult;
  verifyTheoremsError?: string;
  sandboxTemplates: SandboxTemplate[] | null;

  /* Playground Actions */
  initPlayground: () => Promise<void>;
  newPlaygroundSandbox: (
    args: NewSandboxRequest
  ) => Promise<number | undefined>;
  setPlaygroundState: (args: Partial<SandboxData>) => Promise<void>;
  verifyTheorems: () => Promise<void>;
  loadTemplate: (template: SandboxTemplate) => Promise<void>;

  /* Pesisted State Actions */
  findSandboxById: (id?: string | null) => {
    active: PersistedSandboxData | undefined;
    handle: number | null;
  };
  setActiveSandbox: (id: string) => Promise<void>;
  newSandbox: (sandbox: Partial<PersistedSandboxData>) => string;
  deleteStoredSandbox: (id: string) => void;
  updateStoredSandbox: (
    data: Partial<PersistedSandboxData>,
    id?: string
  ) => void;
  mapIdToHandle: (id: string, handle: number) => void;
}

/**
 * Playground Store
 */

export const usePlaygroundStore = create<PlaygroundState>()(
  subscribeWithSelector(
    persist(
      (set, get) => {
        return {
          playgroundStatus: "uninitialized",
          sandboxTemplates: null,
          idHandleMap: {},

          // Inital state
          ...initialStates.persistedPlaygroundData,
          ...initialStates.setStateState,
          ...initialStates.verifyTheoremsState,

          initPlayground: async () => {
            try {
              const {
                playgroundStatus,
                setPlaygroundState,
                newPlaygroundSandbox,
                newSandbox,
              } = get();

              const { active } = getLastActiveSandbox();

              if (playgroundStatus !== "uninitialized") return;

              set({ playgroundStatus: "loading" });

              // Start Loading the playground wasm module
              await loadPlaygroundWasm(WASM_PATH);

              const playground = await window.AcpPlayground.new();
              const sampleResult = await playground.GetSampleSandboxes({});

              // Mark the playground ready
              set({
                playgroundStatus: "ready",
                playground: playground,
                sandboxTemplates: sampleResult?.samples,
              });

              const sandboxToLoad = active ?? sampleResult?.samples[0];

              const sandboxInput = {
                name: sandboxToLoad?.name ?? "default",
                description: sandboxToLoad?.description ?? "",
                data: sandboxToLoad?.data ?? blankSandbox,
              };

              // If there is no loaded sandboxes
              const sandboxId = active?.id || newSandbox(sandboxInput);

              const handle = await newPlaygroundSandbox({
                name: sandboxInput?.name,
                description: sandboxInput?.description,
              });

              if (!handle) {
                // TODO
                throw new Error("Failed to lookup sandbox handle");
              }

              set((state) => ({
                idHandleMap: { ...state.idHandleMap, [sandboxId]: handle },
              }));

              void setPlaygroundState(sandboxInput.data);
            } catch (error) {
              set({
                playgroundStatus: "error",
                playgroundError: (error as Error)?.message,
              });
            }
          },

          newPlaygroundSandbox: async (args) => {
            const { playground } = get();
            if (!playground) throw new Error("Playground is not initialized");

            // Initialize a new sandbox
            const sandbox = await playground.NewSandbox(args);
            const handle = sandbox.record?.handle;

            set({ activeHandle: handle });

            return handle;
          },

          setPlaygroundState: async (updates) => {
            try {
              const { updateStoredSandbox } = get();
              const { playground, handle } = getActiveSandboxHandle();
              const { active } = getLastActiveSandbox();

              // Reset playground state data
              set(initialStates.setStateState);

              const newState = {
                ...(active?.data ?? blankSandbox),
                ...updates,
              };

              updateStoredSandbox({ data: newState });

              const stateResult = await playground?.SetState({
                handle,
                data: newState,
              });

              set({
                setStateDataErrors: stateResult?.errors,
                setStateDataErrorCount: countErrors(stateResult?.errors),
              });

              set(initialStates.verifyTheoremsState);

              // Persist state data into storage
            } catch (error) {
              console.error(error, "Failed to run playground SetState");
              set({ setStateError: (error as Error)?.message });
            }
          },

          verifyTheorems: async () => {
            try {
              const { playground, handle } = getActiveSandboxHandle();

              set(initialStates.verifyTheoremsState);

              const { result } = await playground.VerifyTheorems({ handle });

              set({
                verifyTheoremsStatus: theoremResultPassing(result)
                  ? "passed"
                  : "error",
                verifyTheoremsResult: result,
              });
            } catch (error) {
              console.error(error, "Failed to run playground VerifyTheorems");
              set({ verifyTheoremsError: (error as Error)?.message });
            }
          },

          loadTemplate: async (template) => {
            try {
              const { setPlaygroundState, updateStoredSandbox } = get();
              const { active } = getLastActiveSandbox();

              if (!template.data) throw new Error("Template data not found");
              await setPlaygroundState(template.data);
              updateStoredSandbox({ data: template.data }, active?.id);
            } catch (error) {
              console.error(error, "Failed to load sample template data");
            }
          },

          newSandbox: (input) => {
            const sandboxId = self.crypto.randomUUID();
            const newSandbox: PersistedSandboxData = {
              id: sandboxId,
              name: `name: ""`,
              description: ``,
              data: blankSandbox,
              ...input,
            };

            set((state) => ({
              sandboxes: [newSandbox, ...state.sandboxes],
              lastActiveId: sandboxId,
            }));

            return `${sandboxId}`;
          },

          deleteStoredSandbox: (id) => {
            set((state) => {
              const isActive = state.lastActiveId === id;
              const sandboxes = state.sandboxes.filter((s) => s.id !== id);

              return {
                sandboxes: state.sandboxes.filter((s) => s.id !== id),
                lastActiveId: isActive ? sandboxes[0]?.id : state.lastActiveId,
              };
            });
          },

          updateStoredSandbox: (data, id) => {
            const { lastActiveId } = get();
            const sandboxId = lastActiveId ?? id;

            set((state) => {
              const sandboxes = state.sandboxes.map((sandbox) => {
                return sandbox.id === sandboxId
                  ? { ...sandbox, ...data }
                  : sandbox;
              });

              return { sandboxes };
            });
          },

          setActiveSandbox: async (id: string) => {
            try {
              const {
                newPlaygroundSandbox,
                setPlaygroundState,
                mapIdToHandle,
                findSandboxById,
              } = get();

              const { active, handle: currentSandboxHandle } =
                findSandboxById(id);

              //
              const newSandboxHandle =
                !currentSandboxHandle &&
                (await newPlaygroundSandbox({
                  name: active?.name ?? "",
                  description: active?.description ?? "",
                }));

              // Use existing handle or create a new sandbox
              const sandboxHandle = currentSandboxHandle ?? newSandboxHandle;

              if (newSandboxHandle) mapIdToHandle(id, newSandboxHandle);

              if (!sandboxHandle || !active) {
                // TODO
                throw new Error("Failed to find handle for sandbox");
              }

              set(() => ({
                lastActiveId: id,
                activeHandle: sandboxHandle,
              }));

              setPlaygroundState(active.data);
            } catch (error) {
              console.error(error, "Failed to set sandbox active");
            }
          },

          findSandboxById: (id) => {
            const { sandboxes, idHandleMap } = get();
            return {
              active: sandboxes?.find((s) => s.id === id),
              handle: id ? idHandleMap[id] : null,
            };
          },

          mapIdToHandle: (id: string, handle: number) => {
            set((state) => ({
              idHandleMap: { ...state.idHandleMap, [id]: handle },
            }));
          },
        };
      },
      {
        name: `${STORE_KEY}-2-${STORE_VERSION}`,
        storage: createJSONStorage(() => localStorage), // (optional) by default the 'localStorage' is used
        partialize: (state) => ({
          sandboxes: state.sandboxes,
          lastActiveId: state.lastActiveId,
        }),
      }
    )
  )
);
