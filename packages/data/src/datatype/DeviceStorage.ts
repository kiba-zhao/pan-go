import type { Device } from "./Device";
import type { Storage } from "./Storage";

export type DeviceStorage = {
  id: string;
  peerId: Device["peerId"];
  storageId: Storage["id"];
} & Omit<Storage, "id" | "enabled">;
