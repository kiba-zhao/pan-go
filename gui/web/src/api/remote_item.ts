/**
 * ExtFSRemoteItem API Definition File
 */
import { withPath } from "fetch-utils";
import { fetchMany, fetchOne } from "./base";
import type { ExtFSNodeItem } from "./node_item";
import type { ExtFSRemoteNode } from "./remote_node";

export interface ExtFSRemoteItemAPI {
  searchExtFSRemoteItems(
    peerId: ExtFSRemoteItem["peerId"]
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

/**
 * Search remote items with peer id.
 * @param peerId The peer id.
 * @returns A list of remote items.
 */
export async function searchExtFSRemoteItems(
  peerId: ExtFSRemoteItem["peerId"]
): Promise<ExtFSRemoteItem[]> {
  const [_, remoteItems] = await fetchMany(
    withPath(`extfs/remotes/${peerId}/remote-items`, "merge")
  );
  return remoteItems;
}

/**
 * Select a remote item with id and peer id.
 * @param peerId The peer id.
 * @param itemId The remote item id.
 * @returns The selected remote item.
 */
export async function selectExtFSRemoteItem(
  peerId: ExtFSRemoteItem["peerId"],
  itemId: ExtFSRemoteItem["itemId"]
): Promise<ExtFSRemoteItem> {
  return await fetchOne(
    withPath(`extfs/remotes/${peerId}/remote-items/${itemId}`, "merge")
  );
}
