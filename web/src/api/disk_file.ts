/**
 * DiskFile API Definition File
 */
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

/**
 * Searches for disk files based on the specified search condition.
 * 
 * @param condition - The search criteria including optional parentPath, filePath, 
 *   and fileType to filter the disk files.
 * @returns A promise that resolves to a tuple containing the total number of 
 *   matching disk files and an array of DiskFile objects.
 */

export async function searchDiskFiles(
  condition: DiskFileSearchCondition
): Promise<[number, DiskFile[]]> {
  return await fetchMany(
    withPath("app/disk-files", "merge"),
    withQuery(condition)
  );
}
