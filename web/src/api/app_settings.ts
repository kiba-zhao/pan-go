import { withJSONBody, withMethod, withPath } from "fetch-utils";

import { fetchOne } from "./base";

export interface AppSettingsAPI {
  selectAllAppSettings(): Promise<AppSettings>;
  saveAppSettings(settings: AppSettingsFields): Promise<AppSettings>;
}

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
