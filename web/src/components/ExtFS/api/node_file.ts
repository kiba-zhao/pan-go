/**
 * ExtFSNodeFile API Definition File
 */
import { withPath, withQuery } from "fetch-utils";
import { fetchMany, fetchOne } from "../../../lib/fetch";
import type { ExtFSNodeItem } from "./node_item";

export type ExtFSNodeFileSearchCondition = Partial<
  Pick<ExtFSNodeFile, "parentPath">
>;
export type ExtFSNodeFile = {
  itemId: ExtFSNodeItem["id"];
  parentPath: string;
} & Omit<ExtFSNodeItem, "id" | "enabled">;

/**
 * Search node files with condition.
 * @param itemId The node item id.
 * @param condition The condition to search node files.
 * @returns A list of node files.
 */
export async function searchExtFSNodeFiles(
  itemId: ExtFSNodeFile["itemId"],
  { parentPath, ...opts }: ExtFSNodeFileSearchCondition,
): Promise<ExtFSNodeFile[]> {
  const [_, nodeItems] = await fetchMany(
    withPath(`extfs/node-items/${itemId}/_files`, "merge"),
    withQuery({ parentPath: parentPath || "", ...opts }, "merge"),
  );
  return nodeItems;
}
/**
 * Select a node file with id and file path.
 * @param itemId The node item id.
 * @param filePath The file path.
 * @returns The selected node file.
 */
export async function selectExtFSNodeFile(
  itemId: ExtFSNodeFile["itemId"],
  filePath: ExtFSNodeFile["filePath"],
): Promise<ExtFSNodeFile> {
  return await fetchOne(
    withPath(`extfs/node-items/${itemId}/_files/${filePath}`, "merge"),
  );
}
