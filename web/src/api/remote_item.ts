import { withPath, withQuery } from "fetch-utils";
import { fetchMany, fetchOne } from "./base";
import type { ExtFSNodeItem } from "./node_item";
import type { ExtFSRemoteNode } from "./remote_node";

export interface ExtFSRemoteItemAPI {
  searchExtFSRemoteItems(
    condition: ExtFSRemoteItemSearchCondition
  ): Promise<ExtFSRemoteItem[]>;
  selectExtFSRemoteItem(id: ExtFSRemoteItem["id"]): Promise<ExtFSRemoteItem>;
}

export type ExtFSRemoteItemSearchCondition = {
  peerId: ExtFSRemoteNode["peerId"];
};
export type ExtFSRemoteItem = {
  id: string;
  peerId: ExtFSRemoteNode["peerId"];
  itemId: ExtFSNodeItem["id"];
} & Omit<ExtFSNodeItem, "id" | "enabled">;

export async function searchExtFSRemoteItems(
  condition: ExtFSRemoteItemSearchCondition
): Promise<ExtFSRemoteItem[]> {
  const [_, remoteItems] = await fetchMany(
    withPath("extfs/remote-items", "merge"),
    withQuery(condition, "merge")
  );
  return remoteItems;
}

export async function selectExtFSRemoteItem(
  id: ExtFSRemoteItem["id"]
): Promise<ExtFSRemoteItem> {
  return await fetchOne(withPath(`extfs/remote-items/${id}`, "merge"));
}
