import { PlaygroundState } from "./playgroundStore";

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
};
