import { withPath, withQuery } from "fetch-utils";
import { fetchMany, fetchOne } from "./base";
import type { ExtFSNodeItem } from "./node_item";

export interface ExtFSNodeFileAPI {
  searchExtFSNodeFiles(
    itemId: ExtFSNodeFile["itemId"],
    condition: ExtFSNodeFileSearchCondition
  ): Promise<ExtFSNodeFile[]>;
  selectExtFSNodeFile(
    itemId:ExtFSNodeFile["itemId"],
    filePath:ExtFSNodeFile["filePath"]
  ): Promise<ExtFSNodeFile>;
}

export type ExtFSNodeFileSearchCondition = Partial<Pick<ExtFSNodeFile,"parentPath">>
export type ExtFSNodeFile = {
  itemId: ExtFSNodeItem["id"];
  parentPath: string;
} & Omit<ExtFSNodeItem, "id" | "enabled">;
export async function searchExtFSNodeFiles(
  itemId: ExtFSNodeFile["itemId"],
  {parentPath,...opts}: ExtFSNodeFileSearchCondition): Promise<ExtFSNodeFile[]> {
  const [_, nodeItems] = await fetchMany(
    withPath(`extfs/node-items/${itemId}/_files`, "merge"),
    withQuery(
      { parentPath: parentPath || "", ...opts },
      "merge"
    )
  );
  return nodeItems;
}
export async function selectExtFSNodeFile(
  itemId:ExtFSNodeFile["itemId"],
  filePath:ExtFSNodeFile["filePath"]
): Promise<ExtFSNodeFile> {
  return await fetchOne(withPath(`extfs/node-items/${itemId}/_files/${filePath}`, "merge"));
}
