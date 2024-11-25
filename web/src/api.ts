import {
  simple,
  withJSONBody,
  withMethod,
  withPath,
  withQuery,
} from "fetch-utils";

import jsonServerProvider from "ra-data-json-server";

const ROOT_PATH = `${import.meta.env.BASE_URL}${
  import.meta.env.VITE_API_PATH || ""
}`;

const { fetchOne, fetchMany } = simple(ROOT_PATH);

// export const dataProvider = simpleDataProvider({ fetchOne, fetchMany });
export const dataProvider = jsonServerProvider(ROOT_PATH);

export type AppSettings = {
  name: string;
  rootPath: string;
  peerId: string;
  webAddress: string[];
  peerAddress: string[];
  broadcastAddress: string[];
  publicAddress: string[];
  guardEnabled: boolean;
  guardAccess: boolean;
};

export type AppSettingsFields = Partial<
  Omit<AppSettings, "peerId" | "rootPath">
>;

export async function selectAllAppSettings(): Promise<AppSettings> {
  return await fetchOne(withPath("app/settings", "merge"));
}

export async function saveAppSettings(
  fields: AppSettingsFields
): Promise<AppSettings> {
  const settings_ = await fetchOne(
    withPath("app/settings", "merge"),
    withMethod("PATCH"),
    withJSONBody(fields)
  );
  return settings_;
}

export type DiskFile = {
  id: string;
  name: string;
  filePath: string;
  parentPath: string;
  fileType: string;
  updatedAt: string;
};

export type DiskFileSearchCondition = {
  parentPath?: string;
  filePath?: string;
  fileType?: string;
};

export async function searchDiskFiles(
  condition: DiskFileSearchCondition
): Promise<[number, DiskFile[]]> {
  return await fetchMany(
    withPath("app/disk-files", "merge"),
    withQuery(condition)
  );
}

export type ExtFSRemoteNode = {
  id: string;
  peerId: string;
  name: string;
  available: boolean;
  createdAt: string;
  updatedAt: string;
  tagQuantity: number;
  pendingTagQuantity: number;
};

export async function selectAllExtFSRemoteNodes(): Promise<ExtFSRemoteNode[]> {
  const [_, remoteNodes] = await fetchMany(
    withPath("extfs/remote-nodes", "merge")
  );
  return remoteNodes;
}

export type ExtFSNodeItem = {
  id: number;
  name: string;
  filePath: string;
  fileType: "F" | "D";
  size: number;
  enabled: boolean;
  available: boolean;
  createdAt: string;
  updatedAt: string;
  deletedAt: string;
  tagQuantity: number;
  pendingTagQuantity: number;
};

export type ExtFSNodeItemFields = Omit<
  ExtFSNodeItem,
  | "id"
  | "createdAt"
  | "updatedAt"
  | "deletedAt"
  | "available"
  | "fileType"
  | "size"
  | "tagQuantity"
  | "pendingTagQuantity"
>;

export async function selectAllExtFSNodeItems(): Promise<ExtFSNodeItem[]> {
  const [_, nodeItems] = await fetchMany(withPath("extfs/node-items", "merge"));
  return nodeItems;
}

export async function saveExtFSNodeItem(
  fields: ExtFSNodeItemFields,
  id?: ExtFSNodeItem["id"]
): Promise<ExtFSNodeItem> {
  return await fetchOne(
    withPath(`extfs/node-items${id ? `/${id}` : ""}`, "merge"),
    withMethod(id ? "PATCH" : "POST"),
    withJSONBody(fields)
  );
}

export async function selectExtFSNodeItem(
  id: ExtFSNodeItem["id"]
): Promise<ExtFSNodeItem> {
  return await fetchOne(withPath(`extfs/node-items/${id}`, "merge"));
}

export async function deleteExtFSNodeItem(id: ExtFSNodeItem["id"]) {
  return await fetchOne(
    withPath(`extfs/node-items/${id}`, "merge"),
    withMethod("DELETE")
  );
}

export type ExtFSNodeFileSearchCondition = {
  itemId: ExtFSNodeFile["itemId"];
  parentPath?: string;
};
export type ExtFSNodeFile = {
  id: string;
  itemId: ExtFSNodeItem["id"];
  name: string;
  filePath: string;
  parentPath: string;
  fileType: "F" | "D";
  size: number;
  available: boolean;
  createdAt: string;
  updatedAt: string;
  tagQuantity: number;
  pendingTagQuantity: number;
};
export async function searchExtFSNodeFiles({
  itemId,
  parentPath,
  ...opts
}: ExtFSNodeFileSearchCondition): Promise<ExtFSNodeFile[]> {
  const [_, nodeItems] = await fetchMany(
    withPath("extfs/node-files", "merge"),
    withQuery(
      { itemId: itemId.toString(), parentPath: parentPath || "/", ...opts },
      "merge"
    )
  );
  return nodeItems;
}

export type ExtFSRemoteItemSearchCondition = {
  peerId: ExtFSRemoteNode["peerId"];
};
export type ExtFSRemoteItem = {
  id: string;
  peerId: ExtFSRemoteNode["peerId"];
  itemId: number;
  name: string;
  fileType: "F" | "D";
  size: number;
  available: boolean;
  createdAt: string;
  updatedAt: string;
  tagQuantity: number;
  pendingTagQuantity: number;
};
export async function searchExtFSRemoteItems(
  condition: ExtFSRemoteItemSearchCondition
): Promise<ExtFSRemoteItem[]> {
  const [_, remoteItems] = await fetchMany(
    withPath("extfs/remote-items", "merge"),
    withQuery(condition, "merge")
  );
  return remoteItems;
}

export type ExtFSRemoteFileSearchCondition = {
  peerId: ExtFSRemoteNode["peerId"];
  itemId: ExtFSRemoteItem["itemId"];
  parentPath?: string;
};

export type ExtFSRemoteFile = {
  id: string;
  peerId: ExtFSRemoteNode["peerId"];
  itemId: ExtFSRemoteItem["itemId"];
  name: string;
  filePath: string;
  parentPath: string;
  fileType: "F" | "D";
  size: number;
  available: boolean;
  createdAt: string;
  updatedAt: string;
  tagQuantity: number;
  pendingTagQuantity: number;
};

export async function searchExtFSRemoteFiles({
  itemId,
  parentPath,
  ...opts
}: ExtFSRemoteFileSearchCondition): Promise<ExtFSRemoteFile[]> {
  const [_, remotefiles] = await fetchMany(
    withPath("extfs/remote-files", "merge"),
    withQuery(
      { itemId: itemId.toString(), parentPath: parentPath || "/", ...opts },
      "merge"
    )
  );

  return remotefiles;
}
