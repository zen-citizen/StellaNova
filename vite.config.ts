import { cloudflare } from "@cloudflare/vite-plugin";
import tailwindcss from "@tailwindcss/vite";
import react from "@vitejs/plugin-react";
import { defineConfig } from "vite";

// https://vite.dev/config/
export default defineConfig({
  build: {
    emptyOutDir: true,
    minify: "terser",
    rollupOptions: {
      output: {
        manualChunks: (id: string) =>
          id.includes("node_modules") ? "vendor" : null,
      },
    },
  },
  plugins: [react(), tailwindcss(), cloudflare()],
});
