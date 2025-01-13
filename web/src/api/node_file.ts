import { withPath, withQuery } from "fetch-utils";
import { fetchMany, fetchOne } from "./base";
import type { ExtFSNodeItem } from "./node_item";

export interface ExtFSNodeFileAPI {
  searchExtFSNodeFiles(
    condition: ExtFSNodeFileSearchCondition
  ): Promise<ExtFSNodeFile[]>;
  selectExtFSNodeFile(id: ExtFSNodeFile["id"]): Promise<ExtFSNodeFile>;
}

export type ExtFSNodeFileSearchCondition = {
  itemId: ExtFSNodeFile["itemId"];
  parentPath?: string;
};
export type ExtFSNodeFile = {
  id: string;
  itemId: ExtFSNodeItem["id"];
  parentPath: string;
} & Omit<ExtFSNodeItem, "id" | "enabled">;
export async function searchExtFSNodeFiles({
  itemId,
  parentPath,
  ...opts
}: ExtFSNodeFileSearchCondition): Promise<ExtFSNodeFile[]> {
  const [_, nodeItems] = await fetchMany(
    withPath("extfs/node-files", "merge"),
    withQuery(
      { itemId: itemId.toString(), parentPath: parentPath || "/", ...opts },
      "merge"
    )
  );
  return nodeItems;
}

export async function selectExtFSNodeFile(
  id: ExtFSNodeFile["id"]
): Promise<ExtFSNodeFile> {
  return await fetchOne(withPath(`extfs/node-files/${id}`, "merge"));
}
