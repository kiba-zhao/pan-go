/**
 * vite environment variables
 * 
 * @see https://vite.dev/config/shared-options.html#define
 */

/// <reference types="vite/client" />

interface ImportMetaEnv {
  readonly VITE_PROXY_URL?: string;
  readonly VITE_API_PATH?: string;
  readonly VITE_APP_NAME?: string;
  readonly VITE_I18N_API_ERROR?: string;
}

interface ImportMeta {
  readonly env: ImportMetaEnv;
}
