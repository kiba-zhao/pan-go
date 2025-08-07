import type { Storage } from "./Storage";

export type StorageFile = {
  id: string;
  storageId: Storage["id"];
  parentPath: string;
} & Omit<Storage, "id" | "enabled">;
