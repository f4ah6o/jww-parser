import { defineConfig } from "rolldown-vite";
import path from "node:path";

export default defineConfig({
  root: __dirname,
  base: "./",
  server: {
    port: 5173,
  },
  build: {
    outDir: "dist",
  },
  resolve: {
    alias: [{ find: /^three$/, replacement: path.resolve(__dirname, "src/three-compat.ts") }],
  },
  bundler: "rolldown",
});
