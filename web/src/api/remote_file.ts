import { withPath, withQuery } from "fetch-utils";
import { fetchMany, fetchOne } from "./base";
import type { ExtFSRemoteItem } from "./remote_item";

export interface ExtFSRemoteFileAPI {
  searchExtFSRemoteFiles(
    peerId: ExtFSRemoteFile["peerId"],
    itemId: ExtFSRemoteFile["itemId"],
    condition: ExtFSRemoteFileSearchCondition
  ): Promise<ExtFSRemoteFile[]>;
  selectExtFSRemoteFile(
    peerId: ExtFSRemoteFile["peerId"],
    itemId: ExtFSRemoteFile["itemId"],
    filePath: ExtFSRemoteFile["filePath"]
  ): Promise<ExtFSRemoteFile>;
}

export type ExtFSRemoteFileSearchCondition = Partial<Pick<ExtFSRemoteFile,"parentPath">>

export type ExtFSRemoteFile = {
  parentPath: string;
} & Omit<ExtFSRemoteItem, "id">;

export async function searchExtFSRemoteFiles(
  peerId: ExtFSRemoteFile["peerId"],
  itemId: ExtFSRemoteFile["itemId"],
  {parentPath,...opts}: ExtFSRemoteFileSearchCondition
): Promise<ExtFSRemoteFile[]> {
  const [_, remotefiles] = await fetchMany(
    withPath(`extfs/remotes/${peerId}/remote-items/${itemId}/_files`, "merge"),
    withQuery(
      { parentPath: parentPath || "", ...opts },
      "merge"
    )
  );

  return remotefiles;
}

export async function selectExtFSRemoteFile(
  peerId: ExtFSRemoteFile["peerId"],
  itemId: ExtFSRemoteFile["itemId"],
  filePath: ExtFSRemoteFile["filePath"]
): Promise<ExtFSRemoteFile> {
  return await fetchOne(withPath(`extfs/remotes/${peerId}/remote-items/${itemId}/_files/${filePath}`, "merge"));
}
