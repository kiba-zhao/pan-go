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
