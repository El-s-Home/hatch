import { defineConfig } from "vite";
import react from "@vitejs/plugin-react";
import tailwindcss from "@tailwindcss/vite";

// Build output is written straight into the Go API package so it can be
// embedded into the single binary via embed.FS.
export default defineConfig({
  plugins: [react(), tailwindcss()],
  build: {
    outDir: "../api/internal/web/dist",
    emptyOutDir: true,
  },
  server: {
    proxy: {
      "/v1": "http://localhost:8080",
      "/e": "http://localhost:8080",
    },
  },
});
