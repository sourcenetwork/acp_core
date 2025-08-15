/// <reference types="vite/client" />

interface ImportMetaEnv {
  readonly VITE_WASM_PATH: string;
  readonly VITE_STORE_KEY: string;
  readonly VITE_STORE_VERSION: string;
  readonly VITE_SHARE_URL: string;
  readonly VITE_BACKEND_API: string;
  readonly VITE_BROADCAST_CHANNEL_NAME: string;
  // more env variables...
}

interface ImportMeta {
  readonly env: ImportMetaEnv;
}
