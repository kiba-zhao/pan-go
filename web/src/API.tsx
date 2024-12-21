import type { ReactNode } from "react";
import { createContext, useContext } from "react";
import type {
  AppSettings,
  AppSettingsFields,
  DiskFile,
  DiskFileSearchCondition,
  ExtFSNodeFile,
  ExtFSNodeFileSearchCondition,
  ExtFSNodeItem,
  ExtFSNodeItemFields,
  ExtFSRemoteFile,
  ExtFSRemoteFileSearchCondition,
  ExtFSRemoteItem,
  ExtFSRemoteItemSearchCondition,
  ExtFSRemoteNode,
  ExtFSSearchItemSearchCondition,
  ExtFSSearchItem,
} from "./api";
import * as api from "./api";

export type {
  AppSettings,
  AppSettingsFields,
  DiskFile,
  ExtFSNodeFile,
  ExtFSNodeItem,
  ExtFSNodeItemFields,
  ExtFSRemoteFile,
  ExtFSRemoteItem,
  ExtFSRemoteNode,
  ExtFSSearchItem,
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
  searchExtFSNodeFiles(
    condition: ExtFSNodeFileSearchCondition
  ): Promise<ExtFSNodeFile[]>;
  searchExtFSRemoteItems(
    condition: ExtFSRemoteItemSearchCondition
  ): Promise<ExtFSRemoteItem[]>;
  searchExtFSRemoteFiles(
    condition: ExtFSRemoteFileSearchCondition
  ): Promise<ExtFSRemoteFile[]>;
  searchExtFSSearchItems(
    condition: ExtFSSearchItemSearchCondition
  ): Promise<ExtFSSearchItem[]>;
}

const APIContext = createContext<API | null>(null);

export const useAPI = () => useContext(APIContext);

export const APIProvider = ({ children }: { children: ReactNode }) => (
  <APIContext.Provider value={api}>{children}</APIContext.Provider>
);
