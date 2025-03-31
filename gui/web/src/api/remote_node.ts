/**
 * ExtFSRemoteNode API Definition File
 */
import { withPath } from "fetch-utils";
import { fetchMany } from "./base";

export interface ExtFSRemoteNodeAPI {
  selectAllExtFSRemoteNodes(): Promise<ExtFSRemoteNode[]>;
}

export type ExtFSRemoteNode = {
  id: string;
  peerId: string;
  name: string;
  available: boolean;
  createdAt: string;
  updatedAt: string;
};

/**
 * Select all remote nodes.
 *
 * @returns A list of remote nodes.
 */
export async function selectAllExtFSRemoteNodes(): Promise<ExtFSRemoteNode[]> {
  const [_, remoteNodes] = await fetchMany(
    withPath("extfs/remote-nodes", "merge")
  );
  return remoteNodes;
}
