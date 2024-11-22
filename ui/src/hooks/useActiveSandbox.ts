import { usePlaygroundStore } from "@/lib/playgroundStore";

export const useActiveSandbox = () => {
  const [sandboxes, lastActiveId] = usePlaygroundStore((state) => [
    state.sandboxes,
    state.lastActiveId,
  ]);
  return sandboxes.find((s) => s.id === lastActiveId);
};
