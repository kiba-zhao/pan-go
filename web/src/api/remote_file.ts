import { withPath, withQuery } from "fetch-utils";
import { fetchMany, fetchOne } from "./base";
import type { ExtFSRemoteItem } from "./remote_item";

export interface ExtFSRemoteFileAPI {
  searchExtFSRemoteFiles(
    condition: ExtFSRemoteFileSearchCondition
  ): Promise<ExtFSRemoteFile[]>;
  selectExtFSRemoteFile(id: ExtFSRemoteFile["id"]): Promise<ExtFSRemoteFile>;
}

export type ExtFSRemoteFileSearchCondition = {
  parentPath?: string;
} & Pick<ExtFSRemoteItem, "peerId" | "itemId">;

export type ExtFSRemoteFile = {
  id: string;
  parentPath: string;
} & Omit<ExtFSRemoteItem, "id">;

export async function searchExtFSRemoteFiles({
  itemId,
  parentPath,
  ...opts
}: ExtFSRemoteFileSearchCondition): Promise<ExtFSRemoteFile[]> {
  const [_, remotefiles] = await fetchMany(
    withPath("extfs/remote-files", "merge"),
    withQuery(
      { itemId: itemId.toString(), parentPath: parentPath || "/", ...opts },
      "merge"
    )
  );

  return remotefiles;
}

export async function selectExtFSRemoteFile(
  id: ExtFSRemoteFile["id"]
): Promise<ExtFSRemoteFile> {
  return await fetchOne(withPath(`extfs/remote-files/${id}`, "merge"));
}
