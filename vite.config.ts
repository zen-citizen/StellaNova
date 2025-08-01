import { defineConfig } from "vite";
import react from "@vitejs/plugin-react";
import tailwindcss from "@tailwindcss/vite";

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
  plugins: [react(), tailwindcss()],
});
