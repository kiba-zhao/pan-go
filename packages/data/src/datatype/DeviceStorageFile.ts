import type { DeviceStorage } from "./DeviceStorage";

export type DeviceStorageFile = { id: string; parentPath: string } & Omit<
  DeviceStorage,
  "id"
>;
