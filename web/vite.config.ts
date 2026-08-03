import path from "path";
import react from "@vitejs/plugin-react";
import tailwindcss from "@tailwindcss/vite";
import { UserConfig, defineConfig, loadEnv } from "vite";

// https://vitejs.dev/config/
export default ({ mode }: UserConfig) => {
  const env = loadEnv(mode || "development", process.cwd());
  return defineConfig({
    plugins: [react(), tailwindcss()],
    server: {
      proxy: {
        [`/${env.VITE_API_PATH}`]: env.VITE_API_PROXY_URL,
      },
    },
    resolve: {
      alias: {
        "@": path.resolve(__dirname, "./src"),
      },
    },
  });
};
