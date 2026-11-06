import tailwindcss from "@tailwindcss/vite";
import react from "@vitejs/plugin-react";
import path from "node:path";
import { defineConfig } from "vitest/config";

const backend = process.env.PLAYGROUND2_BACKEND ?? "http://localhost:8090";

export default defineConfig({
  plugins: [tailwindcss(), react()],
  resolve: { alias: { "@": path.resolve(import.meta.dirname, "src") } },
  server: {
    proxy: {
      "/sourcenetwork.acp_core": { target: backend, changeOrigin: true },
      // `acp playground` serves the workspace here (see docs/playground2-acp-handoff.md); a 404 means local mode.
      "/api": { target: backend, changeOrigin: true },
    },
  },
  test: { environment: "node", include: ["src/**/*.test.ts"] },
});
