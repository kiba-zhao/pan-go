import { defineConfig, loadEnv, UserConfig } from "vite";
import react from "@vitejs/plugin-react";

// https://vitejs.dev/config/
export default ({ mode }: UserConfig) => {
  let env = loadEnv(mode || "development", process.cwd());
  return defineConfig({
    plugins: [react()],
    base: "./",
    server: {
      proxy: {
        [env.VITE_PROXY_PATH]: env.VITE_PROXY_URL,
      },
    },
    optimizeDeps: {
      include: ["@mui/material/Tooltip"],
    },
  });
};
