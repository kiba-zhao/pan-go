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
  nodeId: string;
  webAddress: string[];
  nodeAddress: string[];
  broadcastAddress: string[];
  publicAddress: string[];
  guardEnabled: boolean;
  guardAccess: boolean;
};

export type AppSettingsFields = Partial<Omit<AppSettings, "nodeId">>;

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
  filepath: string;
  parent: string;
  fileType: string;
  updatedAt: string;
};

export type DiskFileSearchCondition = {
  parent?: string;
  filepath?: string;
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
  nodeId: string;
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
  filepath: string;
  filetype: "F" | "D";
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
  | "filetype"
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

export type ExtFSFileItemSearchCondition = {
  itemId: ExtFSFileItem["itemId"];
  parentPath?: string;
};
export type ExtFSFileItem = {
  id: string;
  itemId: ExtFSNodeItem["id"];
  name: string;
  filepath: string;
  parentPath: string;
  filetype: "F" | "D";
  size: number;
  available: boolean;
  createdAt: string;
  updatedAt: string;
  tagQuantity: number;
  pendingTagQuantity: number;
};
export async function searchExtFSFileItems({
  itemId,
  parentPath,
  ...opts
}: ExtFSFileItemSearchCondition): Promise<ExtFSFileItem[]> {
  const [_, nodeItems] = await fetchMany(
    withPath("extfs/file-items", "merge"),
    withQuery(
      { itemId: itemId.toString(), parentPath: parentPath || "/", ...opts },
      "merge"
    )
  );
  return nodeItems;
}

export type ExtFSRemoteItemSearchCondition = {
  nodeId: ExtFSRemoteNode["nodeId"];
};
export type ExtFSRemoteItem = {
  id: string;
  nodeId: ExtFSRemoteNode["nodeId"];
  remoteItemId: number;
  name: string;
  filetype: "F" | "D";
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
    withPath("extfs/remote-node-items", "merge"),
    withQuery(condition, "merge")
  );
  return remoteItems;
}

export type ExtFSRemoteFileItemSearchCondition = {
  nodeId: ExtFSRemoteNode["nodeId"];
  itemId: ExtFSRemoteItem["remoteItemId"];
  parentPath?: string;
};

export type ExtFSRemoteFileItem = {
  id: string;
  nodeId: ExtFSRemoteNode["nodeId"];
  itemId: ExtFSRemoteItem["remoteItemId"];
  name: string;
  filepath: string;
  parentPath: string;
  filetype: "F" | "D";
  size: number;
  available: boolean;
  createdAt: string;
  updatedAt: string;
  tagQuantity: number;
  pendingTagQuantity: number;
};

export async function searchExtFSRemoteFileItems({
  itemId,
  parentPath,
  ...opts
}: ExtFSRemoteFileItemSearchCondition): Promise<ExtFSRemoteFileItem[]> {
  const [_, remotefiles] = await fetchMany(
    withPath("extfs/remote-file-items", "merge"),
    withQuery(
      { itemId: itemId.toString(), parentPath: parentPath || "/", ...opts },
      "merge"
    )
  );
  return remotefiles;
}
