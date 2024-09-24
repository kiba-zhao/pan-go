import type { ReactNode } from "react";
import { createContext, useContext } from "react";
import type {
  AppSettings,
  AppSettingsFields,
  DiskFile,
  DiskFileSearchCondition,
  ExtFSFileItem,
  ExtFSFileItemSearchCondition,
  ExtFSNodeItem,
  ExtFSNodeItemFields,
  ExtFSRemoteFileItem,
  ExtFSRemoteFileItemSearchCondition,
  ExtFSRemoteItem,
  ExtFSRemoteItemSearchCondition,
  ExtFSRemoteNode,
} from "./api";
import * as api from "./api";

export type {
  AppSettings,
  AppSettingsFields,
  DiskFile,
  ExtFSFileItem,
  ExtFSNodeItem,
  ExtFSNodeItemFields,
  ExtFSRemoteFileItem,
  ExtFSRemoteItem,
  ExtFSRemoteNode,
};
export interface API {
  selectAllAppSettings(): Promise<AppSettings>;
  saveAppSettings(settings: AppSettingsFields): Promise<AppSettings>;
  searchDiskFiles(
    condition: DiskFileSearchCondition
  ): Promise<[number, DiskFile[]]>;
  selectAllExtFSRemoteNodes(): Promise<ExtFSRemoteNode[]>;
  selectAllExtFSNodeItems(): Promise<ExtFSNodeItem[]>;
  saveExtFSNodeItem(
    fields: ExtFSNodeItemFields,
    id?: ExtFSNodeItem["id"]
  ): Promise<ExtFSNodeItem>;
  selectExtFSNodeItem(id: ExtFSNodeItem["id"]): Promise<ExtFSNodeItem>;
  deleteExtFSNodeItem(id: ExtFSNodeItem["id"]): Promise<void>;
  searchExtFSFileItems(
    condition: ExtFSFileItemSearchCondition
  ): Promise<ExtFSFileItem[]>;
  searchExtFSRemoteItems(
    condition: ExtFSRemoteItemSearchCondition
  ): Promise<ExtFSRemoteItem[]>;
  searchExtFSRemoteFileItems(
    condition: ExtFSRemoteFileItemSearchCondition
  ): Promise<ExtFSRemoteFileItem[]>;
}

const APIContext = createContext<API | null>(null);

export const useAPI = () => useContext(APIContext);

export const APIProvider = ({ children }: { children: ReactNode }) => (
  <APIContext.Provider value={api}>{children}</APIContext.Provider>
);
