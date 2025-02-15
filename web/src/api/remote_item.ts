import { withPath, withQuery } from "fetch-utils";
import { fetchMany, fetchOne } from "./base";
import type { ExtFSNodeItem } from "./node_item";
import type { ExtFSRemoteNode } from "./remote_node";

export interface ExtFSRemoteItemAPI {
  searchExtFSRemoteItems(
    peerId: ExtFSRemoteItem["peerId"],
  ): Promise<ExtFSRemoteItem[]>;
  selectExtFSRemoteItem(
    peerId: ExtFSRemoteItem["peerId"],
    itemId: ExtFSRemoteItem["itemId"]
  ): Promise<ExtFSRemoteItem>;
}

export type ExtFSRemoteItem = {
  id: string;
  peerId: ExtFSRemoteNode["peerId"];
  itemId: ExtFSNodeItem["id"];
} & Omit<ExtFSNodeItem, "id" | "enabled">;

export async function searchExtFSRemoteItems(
  peerId: ExtFSRemoteItem["peerId"],
): Promise<ExtFSRemoteItem[]> {
  const [_, remoteItems] = await fetchMany(
    withPath(`extfs/remotes/${peerId}/remote-items`, "merge"),
  );
  return remoteItems;
}

export async function selectExtFSRemoteItem(
  peerId: ExtFSRemoteItem["peerId"],
  itemId: ExtFSRemoteItem["itemId"]
): Promise<ExtFSRemoteItem> {
  return await fetchOne(withPath(`extfs/remotes/${peerId}/remote-items/${itemId}`, "merge"));
}
