import { sveltekit } from "@sveltejs/kit/vite";
import { defineConfig } from "vite";

const apiProxy = {
  target: "http://localhost:8080",
  changeOrigin: true,
};

export default defineConfig({
  plugins: [sveltekit()],
  server: {
    proxy: {
      "/state": apiProxy,
      "/add": apiProxy,
      "/remove": apiProxy,
      "/move": apiProxy,
    },
  },
  preview: {
    proxy: {
      "/state": apiProxy,
      "/add": apiProxy,
      "/remove": apiProxy,
      "/move": apiProxy,
    },
    allowedHosts: true, // N.B. this is only safe is access to dev server is restricted by other means
  },
});
