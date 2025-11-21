import { usePlaygroundStore } from "@/stores/playgroundStore";

/* Get the active sandbox or optionally a sandbox by id */
export const useSandbox = (id?: string | null) => {
  return usePlaygroundStore((state) => {
    const targetId = id ?? state.lastActiveId;
    return state.sandboxes.find((s) => s.id === targetId);
  });
};
