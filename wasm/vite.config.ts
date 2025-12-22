import { defineConfig } from "rolldown-vite";

export default defineConfig({
  root: __dirname,
  server: {
    port: 5173,
  },
  build: {
    outDir: "dist",
  },
  bundler: "rolldown",
});
