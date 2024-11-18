import { SandboxData } from "@/types/proto-js/sourcenetwork/acp_core/sandbox";
import { create } from "zustand";
import {
  createJSONStorage,
  persist,
  subscribeWithSelector,
} from "zustand/middleware";
import { samples } from "./acpSamples";

const STORE_KEY = "acp-playground-storage";
const STORE_VERSION = "0";

export interface PlaygroundStorageState {
  data: SandboxData;
  updateStore: (payload: Partial<SandboxData>) => void;
  loadSample: (sampleId: string) => void;
}

export const usePlaygroundStorageStore = create<PlaygroundStorageState>()(
  subscribeWithSelector(
    persist(
      (set) => {
        return {
          data: {
            policyDefinition: "",
            policyTheorem: "",
            relationships: "",
          },
          updateStore: (updates) => {
            set(({ data }) => ({ data: { ...data, ...updates } }));
          },
          loadSample: (sampleId) => {
            const sample = samples.get(sampleId);
            if (sample) set(() => ({ data: sample.contents }));
          },
        };
      },
      {
        name: `${STORE_KEY}-${STORE_VERSION}`,
        storage: createJSONStorage(() => localStorage), // (optional) by default the 'localStorage' is used
        partialize: (state) => ({ data: state.data }),
      }
    )
  )
);
