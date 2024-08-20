import type { ReactNode } from "react";
import { createContext, useContext } from "react";
import type {
  AppSettings,
  AppSettingsFields,
  DiskFile,
  DiskFileSearchCondition,
} from "./api";
import * as api from "./api";

export type { AppSettings, AppSettingsFields, DiskFile };
export interface API {
  getAppSettings(): Promise<AppSettings>;
  saveAppSettings(settings: AppSettingsFields): Promise<AppSettings>;
  searchDiskFiles(
    condition: DiskFileSearchCondition
  ): Promise<[number, DiskFile[]]>;
}

const APIContext = createContext<API | null>(null);

export const useAPI = () => useContext(APIContext);

export const APIProvider = ({ children }: { children: ReactNode }) => (
  <APIContext.Provider value={api}>{children}</APIContext.Provider>
);
