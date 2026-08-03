/**
 * vite environment variables
 *
 * @see https://vite.dev/config/shared-options.html#define
 */

/// <reference types="vite/client" />

interface ImportMetaEnv {
  readonly VITE_API_PATH?: string;
  readonly VITE_APP_NAME?: string;
  readonly VITE_APP_VERSION?: string;
  readonly VITE_FAKE_DATA?: boolean;
}

interface ImportMeta {
  readonly env: ImportMetaEnv;
}
