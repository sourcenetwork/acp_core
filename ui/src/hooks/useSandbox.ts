import { usePlaygroundStore } from "@/stores/playgroundStore";

/* Get the active sandbox or optionally a sandbox by id */
export const useSandbox = (id?: string | null) => {
  const sandboxes = usePlaygroundStore((state) => state.sandboxes);
  const lastActiveId = usePlaygroundStore((state) => state.lastActiveId);

  return sandboxes.find((s) => s.id === (id ?? lastActiveId));
};
