import react from "@vitejs/plugin-react";
import { UserConfig, defineConfig, loadEnv } from "vite";

// https://vitejs.dev/config/
export default ({ mode }: UserConfig) => {
  const env = loadEnv(mode || "development", process.cwd());
  return defineConfig({
    plugins: [react()],
    server: {
      proxy: {
        [`/${env.VITE_API_PATH}`]: env.VITE_PROXY_URL,
      },
    },
    optimizeDeps: {
      include: ["@mui/material/Tooltip"],
    },
  });
};
