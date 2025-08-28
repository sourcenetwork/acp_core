import react from "@vitejs/plugin-react";
import { defineConfig, loadEnv } from "vite";
import tsconfigPaths from "vite-tsconfig-paths";

// https://vitejs.dev/config/
export default defineConfig(({ mode }) => {
  // @ts-expect-error
  const env = loadEnv(mode, process.cwd(), "");

  return {
    plugins: [tsconfigPaths(), react()],
    server: {
      proxy: {
        "/api": {
          target: env.VITE_BACKEND_API || "http://localhost:8081",
          changeOrigin: true,
        },
      },
    },
  };
});
