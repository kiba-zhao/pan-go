/// <reference types="vite/client" />

interface ImportMetaEnv {
    readonly VITE_PROXY_URL?: string
    readonly VITE_PROXY_PATH?: string
  }
  
  interface ImportMeta {
    readonly env: ImportMetaEnv
  }