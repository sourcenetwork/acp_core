const importMeta = import.meta.env;

export const WASM_PATH = importMeta.VITE_WASM_PATH ?? "playground.wasm";
export const STORE_KEY = importMeta.VITE_STORE_KEY ?? "acp-playground-storage";
export const STORE_VERSION = importMeta.VITE_STORE_VERSION ?? "0";
export const SHARE_URL = importMeta.VITE_SHARE_URL ?? "";
