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
  ExtFSSearchFile,
  ExtFSSearchFileSearchCondition,
  ExtFSSearchItem,
  ExtFSSearchItemFields,
  ExtFSSearchItemSearchCondition,
  ExtFSSearchFileSearchResults,
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
  ExtFSSearchFile,
  ExtFSSearchItem,
  ExtFSSearchItemFields,
  ExtFSSearchFileSearchResults,
};

interface AppSettingsAPI {
  selectAllAppSettings(): Promise<AppSettings>;
  saveAppSettings(settings: AppSettingsFields): Promise<AppSettings>;
}

interface DiskFileAPI {
  searchDiskFiles(
    condition: DiskFileSearchCondition
  ): Promise<[number, DiskFile[]]>;
}

interface ExtFSRemoteNodeAPI {
  selectAllExtFSRemoteNodes(): Promise<ExtFSRemoteNode[]>;
}

interface ExtFSNodeItemAPI {
  selectAllExtFSNodeItems(): Promise<ExtFSNodeItem[]>;
  saveExtFSNodeItem(
    fields: ExtFSNodeItemFields,
    id?: ExtFSNodeItem["id"]
  ): Promise<ExtFSNodeItem>;
  selectExtFSNodeItem(id: ExtFSNodeItem["id"]): Promise<ExtFSNodeItem>;
  deleteExtFSNodeItem(id: ExtFSNodeItem["id"]): Promise<void>;
}

interface ExtFSNodeFileAPI {
  searchExtFSNodeFiles(
    condition: ExtFSNodeFileSearchCondition
  ): Promise<ExtFSNodeFile[]>;
  selectExtFSNodeFile(id: ExtFSNodeFile["id"]): Promise<ExtFSNodeFile>;
}

interface ExtFSRemoteItemAPI {
  searchExtFSRemoteItems(
    condition: ExtFSRemoteItemSearchCondition
  ): Promise<ExtFSRemoteItem[]>;
  selectExtFSRemoteItem(id: ExtFSRemoteItem["id"]): Promise<ExtFSRemoteItem>;
}

interface ExtFSRemoteFileAPI {
  searchExtFSRemoteFiles(
    condition: ExtFSRemoteFileSearchCondition
  ): Promise<ExtFSRemoteFile[]>;
  selectExtFSRemoteFile(id: ExtFSRemoteFile["id"]): Promise<ExtFSRemoteFile>;
}

interface ExtFSSearchItemAPI {
  searchExtFSSearchItems(
    condition: ExtFSSearchItemSearchCondition
  ): Promise<ExtFSSearchItem[]>;
  deleteExtFSSearchItem(id: ExtFSSearchItem["id"]): Promise<void>;
  saveExtFSSearchItem(
    fields: ExtFSSearchItemFields,
    id?: ExtFSSearchItem["id"]
  ): Promise<ExtFSSearchItem>;
}

interface ExtFSSearchFileAPI {
  searchExtFSSearchFileResults(
    condition: ExtFSSearchFileSearchCondition
  ): Promise<ExtFSSearchFileSearchResults>;
}

export type API = AppSettingsAPI &
  DiskFileAPI &
  ExtFSRemoteNodeAPI &
  ExtFSNodeItemAPI &
  ExtFSNodeFileAPI &
  ExtFSRemoteItemAPI &
  ExtFSRemoteFileAPI &
  ExtFSSearchItemAPI &
  ExtFSSearchFileAPI;

const APIContext = createContext<API | null>(null);

export const useAPI = () => useContext(APIContext);

export const APIProvider = ({ children }: { children: ReactNode }) => (
  <APIContext.Provider value={api}>{children}</APIContext.Provider>
);
