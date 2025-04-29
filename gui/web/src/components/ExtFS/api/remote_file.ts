/**
 * ExtFSRemoteFile API Definition File
 */
import { withPath, withQuery } from "fetch-utils";
import { fetchMany, fetchOne } from "../../../utils/api";
import type { ExtFSRemoteItem } from "./remote_item";

export type ExtFSRemoteFileSearchCondition = Partial<
  Pick<ExtFSRemoteFile, "parentPath">
>;

export type ExtFSRemoteFile = {
  parentPath: string;
} & Omit<ExtFSRemoteItem, "id">;

/**
 * Search remote files with condition.
 * @param peerId The peer id.
 * @param itemId The remote item id.
 * @param condition The condition to search remote files.
 * @returns A list of remote files.
 */
export async function searchExtFSRemoteFiles(
  peerId: ExtFSRemoteFile["peerId"],
  itemId: ExtFSRemoteFile["itemId"],
  { parentPath, ...opts }: ExtFSRemoteFileSearchCondition
): Promise<ExtFSRemoteFile[]> {
  const [_, remotefiles] = await fetchMany(
    withPath(`extfs/remotes/${peerId}/remote-items/${itemId}/_files`, "merge"),
    withQuery({ parentPath: parentPath || "", ...opts }, "merge")
  );

  return remotefiles;
}

/**
 * Select a remote file with id and file path.
 * @param peerId The peer id.
 * @param itemId The remote item id.
 * @param filePath The file path.
 * @returns The selected remote file.
 */
export async function selectExtFSRemoteFile(
  peerId: ExtFSRemoteFile["peerId"],
  itemId: ExtFSRemoteFile["itemId"],
  filePath: ExtFSRemoteFile["filePath"]
): Promise<ExtFSRemoteFile> {
  return await fetchOne(
    withPath(
      `extfs/remotes/${peerId}/remote-items/${itemId}/_files/${filePath}`,
      "merge"
    )
  );
}
