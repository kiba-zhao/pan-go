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
  ExtFSSearchItem,
  ExtFSSearchItemFields,
  ExtFSSearchItemSearchCondition,
  ExtFSSearchFile,
  ExtFSSearchFileSearchCondition,
  ExtFSNodeRefer,
  ExtFSRemoteRefer,
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
  ExtFSSearchItemFields,
  ExtFSSearchFile,
  ExtFSNodeRefer,
  ExtFSRemoteRefer,
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
}

interface ExtFSRemoteItemAPI {
  searchExtFSRemoteItems(
    condition: ExtFSRemoteItemSearchCondition
  ): Promise<ExtFSRemoteItem[]>;
}

interface ExtFSRemoteFileAPI {
  searchExtFSRemoteFiles(
    condition: ExtFSRemoteFileSearchCondition
  ): Promise<ExtFSRemoteFile[]>;
}

interface ExtFSSearchItemAPI {
  searchExtFSSearchItems(
    condition: ExtFSSearchItemSearchCondition
  ): Promise<ExtFSSearchItem[]>;
  deleteExtFSSearchItem(id: ExtFSSearchItem["id"]): Promise<void>;
  saveExtFSSearchItem(fields: ExtFSSearchItemFields): Promise<ExtFSSearchItem>;
}

interface ExtFSSearchFileAPI {
  searchExtFSSearchFiles(
    condition: ExtFSSearchFileSearchCondition
  ): Promise<ExtFSSearchFile[]>;
}

interface ExtFSNodeReferAPI {
  selectExtFSNodeRefer(id: ExtFSNodeRefer["id"]): Promise<ExtFSNodeRefer>;
}

interface ExtFSRemoteReferAPI {
  selectExtFSRemoteRefer(id: ExtFSRemoteRefer["id"]): Promise<ExtFSRemoteRefer>;
}

export type API = AppSettingsAPI &
  DiskFileAPI &
  ExtFSRemoteNodeAPI &
  ExtFSNodeItemAPI &
  ExtFSNodeFileAPI &
  ExtFSRemoteItemAPI &
  ExtFSRemoteFileAPI &
  ExtFSNodeReferAPI &
  ExtFSRemoteReferAPI &
  ExtFSSearchItemAPI &
  ExtFSSearchFileAPI;

const APIContext = createContext<API | null>(null);

export const useAPI = () => useContext(APIContext);

export const APIProvider = ({ children }: { children: ReactNode }) => (
  <APIContext.Provider value={api}>{children}</APIContext.Provider>
);
