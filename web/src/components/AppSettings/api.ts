/**
 * AppSettings API Definition File
 */
import { withJSONBody, withMethod, withPath } from "fetch-utils";

// import { fetchOne } from "../../lib/fetch";

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

/**
 * Fetches the current app settings from the server.
 *
 * @returns the current app settings.
 */
export async function selectAllAppSettings(): Promise<AppSettings> {
  return await fetchOne(withPath("app/settings", "merge"));
}

/**
 * Updates the app settings on the server with the provided fields.
 *
 * @param fields - Partial app settings fields to be updated.
 * @returns A promise that resolves to the updated app settings.
 */
export async function saveAppSettings(
  fields: AppSettingsFields,
): Promise<AppSettings> {
  const settings_ = await fetchOne(
    withPath("app/settings", "merge"),
    withMethod("PATCH"),
    withJSONBody(fields),
  );
  return settings_;
}
