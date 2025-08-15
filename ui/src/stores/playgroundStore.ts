import { STORE_KEY, STORE_VERSION } from "@/constants";
import { getPlaygroundWorkerProxy } from "@/playgroundWorkerProxy";
import { PolicyCatalogue } from "@/types/proto-js/sourcenetwork/acp_core/catalogue";
import { countErrors } from "@/utils/errorUtils";
import { theoremResultPassing } from "@/utils/mapTheoremResultMarkers";
import {
  ExplainCheckRequest,
  ExplainCheckResponse,
  NewSandboxRequest,
  PlaygroundService,
} from "@acp/playground";
import { SandboxData, SandboxDataErrors, SandboxTemplate } from "@acp/sandbox";
import { AnnotatedPolicyTheoremResult } from "@acp/theorem";
import type { IRange } from "monaco-editor";
import { v4 as uuidv4 } from "uuid";
import { create } from "zustand";
import {
  createJSONStorage,
  persist,
  subscribeWithSelector,
} from "zustand/middleware";
import {
  getActiveSandboxHandle,
  getLastActiveSandbox,
} from "../lib/playgroundStoreUtils";

type ExplainCheckFormData = {
  resourceId: string;
  resourceType: string;
  permission: string;
  actorId: string;
};

const initialVerifyTheoremsState: Pick<
  PlaygroundState,
  "verifyTheoremsError" | "verifyTheoremsResult" | "verifyTheoremsStatus"
> = {
  verifyTheoremsError: undefined,
  verifyTheoremsResult: undefined,
  verifyTheoremsStatus: "pending",
};

const initalSetStateState: Pick<
  PlaygroundState,
  "setStateDataErrors" | "setStateError" | "setStateDataErrorCount"
> = {
  setStateDataErrorCount: 0,
  setStateDataErrors: undefined,
  setStateError: undefined,
};

const initialExplainCheckState: Pick<
  PlaygroundState,
  "explainCheckResult" | "explainCheckStatus" | "explainCheckError"
> = {
  explainCheckResult: undefined,
  explainCheckStatus: "pending",
  explainCheckError: undefined,
};

const initialCatalogueState: Pick<
  PlaygroundState,
  "catalogue" | "catalogueStatus" | "catalogueError"
> = {
  catalogue: undefined,
  catalogueStatus: "pending",
  catalogueError: undefined,
};

const initialPersistedPlaygroundData: Pick<
  PlaygroundState,
  "sandboxes" | "lastActiveId"
> = {
  lastActiveId: null,
  sandboxes: [],
};

export const initialStates = {
  persistedPlaygroundData: initialPersistedPlaygroundData,
  setStateState: initalSetStateState,
  verifyTheoremsState: initialVerifyTheoremsState,
  explainCheckState: initialExplainCheckState,
  catalogueState: initialCatalogueState,
};

export const blankSandboxData = (name: string = "new-sandbox"): SandboxData => {
  return {
    policyDefinition: `name: '${name}' \n`,
    policyTheorem: `Authorizations {\n\n}\n\nDelegations {\n}`,
    relationships: "",
  };
};

export const blankSandboxTemplate = (name = "New Sandbox") => {
  return {
    name,
    description: "",
    data: blankSandboxData(name),
  };
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
  playgroundSyncing: boolean; // Whether the playground is syncing with the server
  playgroundStatus: "uninitialized" | "loading" | "ready" | "error";
  playgroundError?: string;
  playground?: PlaygroundService | null;
  setStateDataErrors?: SandboxDataErrors;
  setStateDataErrorCount: number;
  setStateError?: string;
  verifyTheoremsStatus: "pending" | "loading" | "passed" | "failed" | "error";
  verifyTheoremsResult?: AnnotatedPolicyTheoremResult;
  verifyTheoremsError?: string;
  sandboxTemplates: SandboxTemplate[] | null;
  sandboxStateStatus: "unset" | "set" | "error";
  editorSelections: Record<string, IRange | null>;
  explainCheckResult?: ExplainCheckResponse;
  explainCheckStatus: "pending" | "loading" | "passed" | "error";
  explainCheckError?: string | null;
  catalogue?: PolicyCatalogue;
  catalogueStatus: "pending" | "loading" | "passed" | "error";
  catalogueError?: string;
  explainCheckFormData?: Partial<ExplainCheckFormData> | null;

  /* Playground Actions */
  initPlayground: () => Promise<void>;
  newPlaygroundSandbox: (
    args: NewSandboxRequest
  ) => Promise<number | undefined>;
  setPlaygroundState: (args: Partial<SandboxData>) => Promise<void>;
  verifyTheorems: () => Promise<void>;
  loadTemplate: (template: SandboxTemplate) => Promise<void>;
  explainCheck: (input: ExplainCheckRequest) => Promise<void>;
  getCatalogue: (handle: number) => Promise<void>;
  setExplainCheckFormData: (data: Partial<ExplainCheckFormData>) => void;

  /* Pesisted State Actions */
  findSandboxById: (id?: string | null) => {
    active: PersistedSandboxData | undefined;
    handle: number | null;
  };
  setActiveSandbox: (id: string) => Promise<void>;
  newSandbox: (sandbox: Partial<PersistedSandboxData>) => Promise<string>;
  newEmptySandbox: () => Promise<string>;
  deleteStoredSandbox: (id: string) => void;
  updateStoredSandbox: (
    data: Partial<PersistedSandboxData>,
    id?: string | null
  ) => void;
  updateActiveStoredSandbox: (data: Partial<PersistedSandboxData>) => void;
  mapIdToHandle: (id: string, handle: number) => void;
  handleSandboxSyncChangeReceived: (
    sandbox: PersistedSandboxData[]
  ) => Promise<void>;
  setEditorSelection: (id: string, selection: IRange | null) => void;
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
          playgroundSyncing: false,
          sandboxStateStatus: "unset",
          sandboxTemplates: null,
          idHandleMap: {},
          editorSelections: {},

          // Inital state
          ...initialStates.persistedPlaygroundData,
          ...initialStates.setStateState,
          ...initialStates.verifyTheoremsState,
          ...initialStates.explainCheckState,
          ...initialStates.catalogueState,

          // Form State
          explainCheckFormData: null,

          initPlayground: async () => {
            try {
              const { playgroundStatus, newSandbox, setActiveSandbox } = get();

              const { active } = getLastActiveSandbox();
              const activeSandboxId = active?.id;

              if (playgroundStatus !== "uninitialized") return;

              set({ playgroundStatus: "loading" });

              // Start Loading the playground wasm module
              // await loadPlaygroundWasm(WASM_PATH);
              // const playground = await window.AcpPlayground.new();

              const playgroundWorker = await getPlaygroundWorkerProxy();
              const playground = await playgroundWorker.getPlaygroundProxy();
              const sampleResult = await playground.GetSampleSandboxes({});
              const firstTemplate = sampleResult?.samples[0];

              // Mark the playground ready
              set({
                playgroundStatus: "ready",
                playground: playground,
                sandboxTemplates: sampleResult?.samples,
              });

              if (activeSandboxId != null) {
                await setActiveSandbox(activeSandboxId);
                return;
              }

              await newSandbox({
                name: firstTemplate?.name ?? "New Sandbox",
                description: firstTemplate?.description ?? "",
                data: firstTemplate?.data ?? blankSandboxData(),
              });
            } catch (error) {
              console.error(error, "Failed to initialize playground");
              set({
                playgroundStatus: "error",
                playgroundError: (error as Error)?.message,
              });
            }
          },

          newPlaygroundSandbox: async (args) => {
            const { playground } = get();
            if (!playground) throw new Error("Playground is not initialized");

            const sandbox = await playground.NewSandbox(args);

            return sandbox.record?.handle;
          },

          setPlaygroundState: async (updates) => {
            try {
              const { updateActiveStoredSandbox, getCatalogue } = get();
              const { playground, handle } = getActiveSandboxHandle();
              const { active } = getLastActiveSandbox();

              // Reset playground state data
              set(initialStates.setStateState);

              set({ playgroundSyncing: true });

              const newState = {
                ...(active?.data ?? blankSandboxData()),
                ...updates,
              };

              updateActiveStoredSandbox({ data: newState });

              const setStateInput = { handle, data: newState };
              const stateResult = await playground?.SetState(setStateInput);

              set({
                sandboxStateStatus: stateResult.ok ? "set" : "error",
                setStateDataErrors: stateResult?.errors,
                setStateDataErrorCount: countErrors(stateResult?.errors),
                // Reset theorems state as invalidated by the new state
                ...initialStates.verifyTheoremsState,
              });

              // Get the catalogue for the sandbox
              void getCatalogue(handle);
            } catch (error) {
              set({
                sandboxStateStatus: "error",
                setStateError: (error as Error)?.message,
              });
            } finally {
              set({ playgroundSyncing: false });
            }
          },

          verifyTheorems: async () => {
            try {
              const { playground, handle } = getActiveSandboxHandle();

              // Reset state for theorems
              set(initialStates.verifyTheoremsState);
              set({ verifyTheoremsStatus: "loading" });

              const { result } = await playground.VerifyTheorems({ handle });

              set({
                verifyTheoremsStatus: theoremResultPassing(result)
                  ? "passed"
                  : "failed",
                verifyTheoremsResult: result,
              });
            } catch (error) {
              set({
                verifyTheoremsStatus: "error",
                verifyTheoremsError: (error as Error)?.message,
              });
            }
          },

          loadTemplate: async (template) => {
            try {
              const { setPlaygroundState } = get();

              if (!template.data) throw new Error("Template data not found");

              await setPlaygroundState(template.data);
            } catch (error) {
              console.error(error, "Failed to load sample template data");
            }
          },

          newSandbox: async (input, setActive = true) => {
            const { setActiveSandbox } = get();
            const sandboxId = uuidv4();

            const newSandbox: PersistedSandboxData = {
              id: sandboxId,
              ...blankSandboxTemplate(input.name),
              ...input,
            };

            set((state) => ({ sandboxes: [newSandbox, ...state.sandboxes] }));

            if (setActive === true) {
              try {
                await setActiveSandbox(sandboxId);
              } catch (error) {
                console.error(error, "Failed to set sandbox active"); // TODO
              }
            }

            return `${sandboxId}`;
          },

          newEmptySandbox: async () => {
            return await get().newSandbox(blankSandboxTemplate());
          },

          deleteStoredSandbox: (id) => {
            const { setActiveSandbox, newEmptySandbox } = get();

            set((state) => {
              const isActive = state.lastActiveId === id;
              const sandboxes = state.sandboxes.filter((s) => s.id !== id);
              const idHandleMap = { ...state.idHandleMap };
              delete idHandleMap[id];

              return {
                sandboxes: sandboxes,
                idHandleMap: idHandleMap,
                lastActiveId: isActive ? sandboxes[0]?.id : state.lastActiveId,
              };
            });

            const { sandboxes, lastActiveId } = get();

            // If there is a new active id, set it active
            if (lastActiveId) void setActiveSandbox(lastActiveId);

            // If there are no sandboxes, replace with an empty one
            if (!sandboxes.length) void newEmptySandbox();
          },

          updateStoredSandbox: (data, id) => {
            const { lastActiveId } = get();
            const sandboxId = id ?? lastActiveId;

            set((state) => {
              const sandboxes = state.sandboxes.map((sandbox) => {
                return sandbox.id === sandboxId
                  ? { ...sandbox, ...data }
                  : sandbox;
              });

              return { sandboxes };
            });
          },

          updateActiveStoredSandbox: (data) => {
            const { lastActiveId, updateStoredSandbox } = get();
            updateStoredSandbox(data, lastActiveId);
          },

          setActiveSandbox: async (id: string) => {
            try {
              const {
                newPlaygroundSandbox,
                setPlaygroundState,
                mapIdToHandle,
                findSandboxById,
              } = get();

              const { active, handle: activeHandle } = findSandboxById(id);

              const sandboxHandle =
                activeHandle ??
                (await newPlaygroundSandbox({
                  name: active?.name ?? "",
                  description: active?.description ?? "",
                }));

              // TODO
              if (!sandboxHandle || !active)
                throw new Error("Failed to find handle for sandbox");

              set(() => ({
                lastActiveId: id,
                activeHandle: sandboxHandle,
              }));

              void mapIdToHandle(id, sandboxHandle);
              void setPlaygroundState(active.data);
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

          handleSandboxSyncChangeReceived: async (
            sandbox: PersistedSandboxData[]
          ) => {
            const { setPlaygroundState, lastActiveId } = get();
            const activeSandbox = sandbox.find((s) => s.id === lastActiveId);
            if (!activeSandbox) return;
            await setPlaygroundState(activeSandbox.data);
          },
          setEditorSelection: (id, selection) => {
            set((state) => ({
              editorSelections: { ...state.editorSelections, [id]: selection },
            }));
          },

          explainCheck: async (input) => {
            const { playground } = get();
            if (!playground) throw new Error("Playground is not initialized");

            try {
              set({ explainCheckStatus: "loading" });

              const result = await playground.ExplainCheck(input);

              set({
                explainCheckStatus: "passed",
                explainCheckResult: result,
                explainCheckError: null,
              });
            } catch (error) {
              console.error(error, "Failed to explain check");
              set({
                explainCheckStatus: "error",
                explainCheckError: (error as Error)?.message,
              });
            }
          },

          setExplainCheckFormData: (data) => {
            const { explainCheckFormData } = get();
            set({ explainCheckFormData: { ...explainCheckFormData, ...data } });
          },

          getCatalogue: async (handle) => {
            const { playground, activeHandle } = get();
            if (!playground) throw new Error("Playground is not initialized");

            try {
              set({ catalogueStatus: "loading" });

              const targetHandle = handle ?? activeHandle ?? 0;

              const { catalogue } = await playground.GetCatalogue({
                handle: targetHandle,
              });

              set({ catalogue, catalogueStatus: "passed" });
            } catch (error) {
              console.error(error, "Failed to get catalogue");
              set({
                catalogueStatus: "error",
                catalogueError: (error as Error)?.message,
              });
            }
          },
        };
      },
      {
        name: `${STORE_KEY}-2-${STORE_VERSION}`,
        storage: createJSONStorage(() => localStorage), // (optional) by default the 'localStorage' is used
        partialize: (state) => ({
          sandboxes: state.sandboxes,
          lastActiveId: state.lastActiveId,
          explainCheckFormData: state.explainCheckFormData,
        }),
      }
    )
  )
);
