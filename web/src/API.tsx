import type { ReactNode } from "react";
import { createContext, useContext } from "react";
import type {
  AppSettings,
  AppSettingsFields,
  DiskFile,
  DiskFileSearchCondition,
  ExtFSItem,
  ExtFSItemSearchCondition,
  ExtFSNodeItem,
  ExtFSNodeItemFields,
} from "./api";
import * as api from "./api";

export type {
  AppSettings,
  AppSettingsFields,
  DiskFile,
  ExtFSItem,
  ExtFSItemSearchCondition,
  ExtFSNodeItem,
  ExtFSNodeItemFields,
};
export interface API {
  getAppSettings(): Promise<AppSettings>;
  saveAppSettings(settings: AppSettingsFields): Promise<AppSettings>;
  searchDiskFiles(
    condition: DiskFileSearchCondition
  ): Promise<[number, DiskFile[]]>;
  searchExtFSItems(
    condition: ExtFSItemSearchCondition
  ): Promise<[number, ExtFSItem[]]>;
  saveExtFSNodeItem(
    fields: ExtFSNodeItemFields,
    id?: ExtFSNodeItem["id"]
  ): Promise<ExtFSNodeItem>;
  selectExtFSNodeItem(id: ExtFSNodeItem["id"]): Promise<ExtFSNodeItem>;
  deleteExtFSNodeItem(id: ExtFSNodeItem["id"]): Promise<void>;
}

const APIContext = createContext<API | null>(null);

export const useAPI = () => useContext(APIContext);

export const APIProvider = ({ children }: { children: ReactNode }) => (
  <APIContext.Provider value={api}>{children}</APIContext.Provider>
);
