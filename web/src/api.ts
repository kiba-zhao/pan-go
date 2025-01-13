export * from "./api/app_settings";
export * from "./api/base";
export * from "./api/disk_file";
export * from "./api/node_file";
export * from "./api/node_item";
export * from "./api/remote_file";
export * from "./api/remote_item";
export * from "./api/remote_node";
export * from "./api/remote_search_file";
export * from "./api/search_file";
export * from "./api/search_item";

import type { AppSettingsAPI } from "./api/app_settings";
import type { DiskFileAPI } from "./api/disk_file";
import type { ExtFSNodeFileAPI } from "./api/node_file";
import type { ExtFSNodeItemAPI } from "./api/node_item";
import type { ExtFSRemoteFileAPI } from "./api/remote_file";
import type { ExtFSRemoteItemAPI } from "./api/remote_item";
import type { ExtFSRemoteNodeAPI } from "./api/remote_node";
import type { ExtFSRemoteSearchFileAPI } from "./api/remote_search_file";
import type { ExtFSSearchFileAPI } from "./api/search_file";
import type { ExtFSSearchItemAPI } from "./api/search_item";

export type API = AppSettingsAPI &
  DiskFileAPI &
  ExtFSRemoteNodeAPI &
  ExtFSNodeItemAPI &
  ExtFSNodeFileAPI &
  ExtFSRemoteItemAPI &
  ExtFSRemoteFileAPI &
  ExtFSSearchItemAPI &
  ExtFSSearchFileAPI &
  ExtFSRemoteSearchFileAPI;
