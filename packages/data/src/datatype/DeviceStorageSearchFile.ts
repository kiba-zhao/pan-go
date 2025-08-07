import type { DeviceStorageFile } from "./DeviceStorageFile";

export type DeviceStorageSearchFile = {
  id: number;
  score: number;
  tokens: string[];
} & Omit<DeviceStorageFile, "id">;
