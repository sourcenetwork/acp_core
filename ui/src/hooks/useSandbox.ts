import { usePlaygroundStore } from "@/lib/playgroundStore";

/* Get the active sandbox or optionally a sandbox by id */
export const useSandbox = (id?: string | null) => {
  const [sandboxes, lastActiveId] = usePlaygroundStore((state) => [
    state.sandboxes,
    state.lastActiveId,
  ]);

  return sandboxes.find((s) => s.id === (id ?? lastActiveId));
};
