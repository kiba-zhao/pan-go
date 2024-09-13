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

export async function getAppSettings(): Promise<AppSettings> {
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

export type ExtFSItem = {
  id: string;
  fileType: "S" | "D" | "F" | "N" | "RD" | "RF" | "RN";
  name: string;
  updatedAt: string;
  tagQuantity: number;
  pendingTagQuantity: number;
  disabled?: boolean;
  parentId?: string;
  linkId?: string;
};

export type ExtFSItemSearchCondition = Pick<ExtFSItem, "parentId"> & {
  fileType?: Array<ExtFSItem["fileType"]> | ExtFSItem["fileType"];
};

export async function searchExtFSItems(
  condition: ExtFSItemSearchCondition
): Promise<[number, ExtFSItem[]]> {
  const { fileType, ...params } = condition;
  const query = new URLSearchParams(params);
  if (Array.isArray(fileType)) {
    fileType.reduce((q, t) => {
      q.append("fileType", t);
      return q;
    }, query);
  } else if (fileType !== void 0) {
    query.append("fileType", fileType);
  }
  return await fetchMany(withPath("extfs/items", "merge"), withQuery(query));
}

export type ExtFSNodeItem = {
  id: number;
  name: string;
  filepath: string;
  enabled: boolean;
  available: boolean;
  createdAt: string;
  updatedAt: string;
  deletedAt: string;
};

export type ExtFSNodeItemFields = Omit<
  ExtFSNodeItem,
  "id" | "createdAt" | "updatedAt" | "deletedAt" | "available"
>;

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
