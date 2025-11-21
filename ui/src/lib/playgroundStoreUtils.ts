import { usePlaygroundStore } from "../stores/playgroundStore";

export const getLastActiveSandbox = () => {
  const { findSandboxById, lastActiveId } = usePlaygroundStore.getState();
  return findSandboxById(lastActiveId);
};

export const getActiveSandboxHandle = () => {
  const { playground, activeHandle } = usePlaygroundStore.getState();
  if (!playground || !activeHandle)
    throw new Error("Playground or sandbox is not initialized");
  return { playground, handle: activeHandle };
};
