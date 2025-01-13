import { withPath, withQuery } from "fetch-utils";
import { fetchMany } from "./base";

export interface DiskFileAPI {
  searchDiskFiles(
    condition: DiskFileSearchCondition
  ): Promise<[number, DiskFile[]]>;
}

export type DiskFile = {
  id: string;
  name: string;
  filePath: string;
  parentPath: string;
  fileType: string;
  updatedAt: string;
};

export type DiskFileSearchCondition = {
  parentPath?: string;
  filePath?: string;
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
