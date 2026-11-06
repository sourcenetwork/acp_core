import { create } from "zustand";
import { createJSONStorage, persist, type StateStorage } from "zustand/middleware";
import { ApiError, generate, setFileBacked, verify } from "@/lib/api";
import { clampConfig, DEFAULT_CONFIG } from "@/lib/config";
import { parseConstraints, type ConstraintError } from "@/lib/constraints";
import { parseRequirements, type RequirementError } from "@/lib/requirements";
import { clampWorldIndex } from "@/lib/worlds";
import {
  CONSTRAINTS_FILE, CONTRACTS_FILE, connectEvents, detectMode, fetchWorkspace, mergePrefs, mergeWorlds, parseWorldFiles, reselect,
  type Connection, type FileProblem, type Mode, type WorkspaceFiles,
} from "@/lib/workspaceFiles";
import { SAMPLE_CONSTRAINTS, SAMPLE_POLICY, SAMPLE_REQUIREMENTS } from "@/lib/sample";
import type { GenerateConfig, VerifyResponse, World } from "@/lib/types";

export type GenerateError =
  | { kind: "constraints"; errors: ConstraintError[]; file?: string }
  | { kind: "requirements"; errors: RequirementError[]; file?: string }
  | { kind: "backend" | "network" | "protocol"; message: string; code?: string };

export type MainTab = "policy" | "constraints" | "requirements" | "worlds";

/** A Verify outcome with the inputs it was computed from, to tell when it has gone stale. Ephemeral. */
export interface VerifyRun {
  response: VerifyResponse;
  scope: number;
  inputs: { policy: string; constraints: string; requirements: string };
}

export const DEFAULT_VERIFY_SCOPE = 3;

// What file-backed mode still persists: the rest is the workspace's own (see mergePrefs).
const PREF_KEYS = ["config", "verifyScope", "autoVerify"];

interface State {
  // How the workspace is backed: detected once at startup (initWorkspace). Never persisted.
  mode: Mode;
  /** File-backed mode: the directory, the version of the files held, why the last refresh failed, the event stream. */
  dir: string | null;
  workspaceVersion: string | null;
  workspaceError: string | null;
  connection: Connection;
  /** File-backed mode: saved worlds whose file doesn't parse, each reported on its own. */
  worldProblems: FileProblem[];
  /** Verify again when the files change (file-backed mode). A preference. */
  autoVerify: boolean;
  setAutoVerify: (on: boolean) => void;
  initWorkspace: () => Promise<void>;
  applyFiles: (files: WorkspaceFiles) => void;

  // Workspace: one system model + a list of worlds (saved | ephemeral).
  policy: string;
  constraints: string;
  worlds: World[];
  config: GenerateConfig;
  requirements: string;
  verifyScope: number;

  generating: boolean;
  error: GenerateError | null;
  verifying: boolean;
  verifyRun: VerifyRun | null;

  // Transient UI state (not persisted).
  activeTab: MainTab;
  selectedWorld: number;
  setActiveTab: (t: MainTab) => void;
  setSelectedWorld: (i: number) => void;

  setPolicy: (v: string) => void;
  setConstraints: (v: string) => void;
  setRequirements: (v: string) => void;
  setVerifyScope: (n: number) => void;
  setConfig: (c: Partial<GenerateConfig>) => void;
  verify: () => Promise<void>;
  cancelVerify: () => void;
  generate: () => Promise<void>;
}

let controller: AbortController | null = null;
let verifyController: AbortController | null = null;
let nextId = 0;
let stopEvents: (() => void) | null = null;
const fileOf = (mode: Mode, file: string) => (mode === "files" ? file : undefined);

// In file-backed mode only the preferences are written: the local-mode workspace in storage is kept as it was.
const storage: StateStorage = {
  getItem: (name) => localStorage.getItem(name),
  setItem: (name, value) => {
    const files = useWorkspace.getState().mode === "files";
    localStorage.setItem(name, files ? mergePrefs(localStorage.getItem(name), value, PREF_KEYS) : value);
  },
  removeItem: (name) => localStorage.removeItem(name),
};

export const useWorkspace = create<State>()(
  persist(
    (set, get) => ({
      mode: "pending",
      dir: null,
      workspaceVersion: null,
      workspaceError: null,
      connection: "connecting",
      worldProblems: [],
      autoVerify: false,
      policy: SAMPLE_POLICY,
      constraints: SAMPLE_CONSTRAINTS,
      worlds: [],
      config: DEFAULT_CONFIG,
      requirements: SAMPLE_REQUIREMENTS,
      verifyScope: DEFAULT_VERIFY_SCOPE,
      generating: false,
      error: null,
      verifying: false,
      verifyRun: null,
      activeTab: "policy",
      selectedWorld: 0,
      setActiveTab: (activeTab) => set({ activeTab }),
      setSelectedWorld: (i) => set({ selectedWorld: clampWorldIndex(i, get().worlds.length) }),

      // The files are the source of truth in file-backed mode: edits are ignored.
      setPolicy: (policy) => get().mode !== "files" && set({ policy }),
      setConstraints: (constraints) => get().mode !== "files" && set({ constraints }),
      setRequirements: (requirements) => get().mode !== "files" && set({ requirements }),
      setAutoVerify: (autoVerify) => {
        set({ autoVerify });
        const { verifyRun: run, policy, constraints, requirements } = get();
        const stale = run !== null && (run.inputs.policy !== policy || run.inputs.constraints !== constraints || run.inputs.requirements !== requirements);
        if (autoVerify && stale) void get().verify();
      },

      initWorkspace: async () => {
        if (get().mode !== "pending") return;
        const d = await detectMode();
        if (d.mode === "local") {
          set({ mode: "local" });
          return;
        }
        setFileBacked(true);
        if (d.files) get().applyFiles(d.files);
        // unreadable from the start: nothing of ours to show (not the local workspace) until the files are readable
        if (!d.files) set({ mode: "files", workspaceError: d.error, policy: "", constraints: "", requirements: "", worlds: [], selectedWorld: 0 });
        stopEvents?.();
        stopEvents = connectEvents({
          connection: (connection) => set({ connection }),
          version: () => get().workspaceVersion,
          refetch: async () => {
            const r = await fetchWorkspace();
            // keep the last good state on an error: the next change retries
            // (an unreachable acp is the connection indicator's business, not a workspace error)
            if ("error" in r) {
              if (r.reachable) set({ workspaceError: r.error });
            } else get().applyFiles(r.files);
            return !("error" in r) || r.reachable;
          },
        });
      },

      // Replaces the workspace's files in the store; everything parsed from them is derived from the new text.
      applyFiles: (files) => {
        const prev = get();
        if (prev.workspaceVersion === files.version && prev.mode === "files") {
          if (prev.workspaceError !== null) set({ workspaceError: null }); // readable again, same files as before
          return;
        }
        const { worlds: saved, problems } = parseWorldFiles(files.worlds);
        const worlds = mergeWorlds(prev.worlds, saved);
        const changed = prev.policy !== files.policy || prev.constraints !== files.constraints || prev.requirements !== files.contracts;
        set({
          mode: "files",
          dir: files.dir,
          workspaceVersion: files.version,
          workspaceError: null,
          worldProblems: problems,
          policy: files.policy,
          constraints: files.constraints,
          requirements: files.contracts,
          worlds,
          selectedWorld: reselect(prev.worlds, prev.selectedWorld, worlds),
        });
        // a stale error about the old text is no longer true (parse errors are derived from the files)
        if (changed && (prev.error?.kind === "constraints" || prev.error?.kind === "requirements")) set({ error: null });
        if (changed && prev.autoVerify && prev.mode === "files") void get().verify();
      },
      setVerifyScope: (verifyScope) => set({ verifyScope }),
      setConfig: (c) => set({ config: { ...get().config, ...c } }),

      cancelVerify: () => {
        verifyController?.abort();
        verifyController = null;
        set({ verifying: false });
      },

      // Verify the current Policy, constraints and requirements. A newer run aborts the previous one, whose response
      // is then ignored. The scope is not capped here: the backend's error is shown.
      verify: async () => {
        const { policy, constraints, requirements, verifyScope } = get();
        const pc = parseConstraints(constraints);
        if (!pc.ok) {
          set({ error: { kind: "constraints", errors: pc.errors, file: fileOf(get().mode, CONSTRAINTS_FILE) } });
          return;
        }
        const pr = parseRequirements(requirements);
        if (!pr.ok) {
          set({ error: { kind: "requirements", errors: pr.errors, file: fileOf(get().mode, CONTRACTS_FILE) } });
          return;
        }
        verifyController?.abort();
        const c = (verifyController = new AbortController());
        const scope = Number.isFinite(verifyScope) ? Math.max(1, Math.round(verifyScope)) : DEFAULT_VERIFY_SCOPE;
        set({ verifying: true, error: null, verifyScope: scope });
        try {
          const response = await verify(
            { systemModel: { policy, constraints: pc.constraints, requirements: pr.requirements }, config: { scope } },
            c.signal,
          );
          if (c.signal.aborted) return;
          set({ verifying: false, verifyRun: { response, scope: response.scope ?? scope, inputs: { policy, constraints, requirements } } });
        } catch (e) {
          if (c.signal.aborted || (e instanceof DOMException && e.name === "AbortError")) return;
          const err = e instanceof ApiError ? e : new ApiError("network", String(e));
          set({ verifying: false, error: { kind: err.kind, message: err.message, code: err.code } });
        }
      },

      generate: async () => {
        const { policy, constraints, config } = get();
        const parsed = parseConstraints(constraints);
        if (!parsed.ok) {
          set({ error: { kind: "constraints", errors: parsed.errors, file: fileOf(get().mode, CONSTRAINTS_FILE) } });
          return;
        }
        controller?.abort();
        const c = (controller = new AbortController());
        const cfg = clampConfig(config);
        set({ generating: true, error: null, config: cfg });
        try {
          const res = await generate({ systemModel: { policy, constraints: parsed.constraints }, config: cfg }, c.signal);
          // Generated worlds are ephemeral: a new Generate replaces them (saved worlds are kept).
          const fresh: World[] = (res.worlds ?? []).map((w) => ({
            id: `w${++nextId}`,
            kind: "ephemeral",
            relationships: w.relationships ?? "",
          }));
          set({
            worlds: [...get().worlds.filter((w) => w.kind === "saved"), ...fresh],
            generating: false,
            selectedWorld: 0,
            activeTab: "worlds",
          });
        } catch (e) {
          if (e instanceof DOMException && e.name === "AbortError") return;
          const err = e instanceof ApiError ? e : new ApiError("network", String(e));
          set({ generating: false, error: { kind: err.kind, message: err.message, code: err.code } });
        }
      },
    }),
    {
      name: "playground2.workspace",
      version: 1,
      storage: createJSONStorage(() => storage),
      // Ephemeral worlds and transient UI state are not persisted.
      partialize: (s) => ({
        policy: s.policy,
        constraints: s.constraints,
        config: s.config,
        requirements: s.requirements,
        verifyScope: s.verifyScope,
        autoVerify: s.autoVerify,
        worlds: s.worlds.filter((w) => w.kind === "saved"),
      }),
    },
  ),
);
