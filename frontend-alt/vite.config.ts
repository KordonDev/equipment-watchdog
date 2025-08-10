import { defineConfig } from "vite";
import { svelte } from "@sveltejs/vite-plugin-svelte";
import dns from "dns";

dns.setDefaultResultOrder("verbatim");

// https://vitejs.dev/config/
export default defineConfig({
  plugins: [svelte()],
});
