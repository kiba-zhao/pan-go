import type { StorageFile } from "./StorageFile";

export type StorageSearchFile = {
  id: number;
  score: number;
  tokens: string[];
} & Omit<StorageFile, "id">;
