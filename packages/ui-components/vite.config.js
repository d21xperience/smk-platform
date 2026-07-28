import vue from "@vitejs/plugin-vue";
import { resolve } from "path";
import { defineConfig } from "vite";

export default defineConfig({
  plugins: [vue()],
  build: {
    lib: {
      entry: resolve(__dirname, "src/index.js"),
      name: "UiComponents",
      formats: ["es"],
      fileName: "index",
    },
    rollupOptions: {
      external: ["vue", "quasar", "vue-router"],
      output: {
        globals: {
          vue: "Vue",
          quasar: "Quasar",
          "vue-router": "VueRouter",
        },
      },
    },
  },
});
