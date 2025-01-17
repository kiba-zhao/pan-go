import type { ExtFSItemRecord } from "./Item";
import { ExtFSItem, ExtFSItems, ExtFSItemTag, useExtFSItem } from "./Item";
import { More, MoreHelpItem } from "./More";
import type { ExtFSSingleState, ExtFSState } from "./State";
import { useExtFS } from "./State";

import { useQuery } from "@tanstack/react-query";
import type { ExtFSRemoteFile } from "../../api";
import { useAPI } from "../API";

import { useMemo } from "react";

import CloudIcon from "@mui/icons-material/Cloud";
import FolderIcon from "@mui/icons-material/Folder";
import InsertDriveFileIcon from "@mui/icons-material/InsertDriveFile";

const ExtFSRemoteFileTagRoutePath = "/extfs/remote-file-tags";
export const ExtFSRemoteFileMode = "RF";
const ExtFSRemoteFileQueryKey = ["extfs-remote-files"];
const ExtFSRemoteFileState = {
  mode: ExtFSRemoteFileMode,
  queryKeyList: [ExtFSRemoteFileQueryKey],
};

type NewExtFSStateOpts = {
  filePath?: ExtFSRemoteFile["filePath"];
} & Pick<ExtFSRemoteFile, "peerId" | "itemId" | "name">;
export function newExtFSState(
  extfs: ExtFSState,
  opts: NewExtFSStateOpts
): ExtFSState {
  const { parentItems } = extfs;
  const fileState = {
    ...ExtFSRemoteFileState,
    peerId: opts.peerId,
    itemId: opts.itemId,
    parentPath: opts.filePath,
  };
  return {
    ...fileState,
    parentItems: [...parentItems, { name: opts.name, state: fileState }],
  } as ExtFSState;
}

export type ExtFSRemoteFileSingleState = {
  peerId: string;
  itemId: number;
  parentPath?: string;
} & ExtFSSingleState;
export const RemoteFiles = () => {
  const [{ parentItems, ...state }, _] = useExtFS();
  const { peerId, itemId, parentPath } = state as ExtFSRemoteFileSingleState;

  const api = useAPI();
  const { data: items, isFetching } = useQuery({
    queryKey: [...ExtFSRemoteFileQueryKey, { peerId, itemId, parentPath }],
    queryFn: async () =>
      await api?.searchExtFSRemoteFiles({ peerId, itemId, parentPath }),
    enabled: !!api && state.mode === ExtFSRemoteFileMode,
  });

  return (
    <ExtFSItems items={items || []} isFetching={isFetching}>
      <RemoteFile />
    </ExtFSItems>
  );
};

export const RemoteFile = () => {
  const { style, item }: ExtFSItemRecord<ExtFSRemoteFile> = useExtFSItem();
  const avatarIcon = useMemo(() => {
    if (item?.fileType === "D")
      return <FolderIcon color={item.available ? "primary" : "disabled"} />;
    if (item?.fileType === "F")
      return (
        <InsertDriveFileIcon color={item.available ? "action" : "disabled"} />
      );
  }, [item?.fileType]);

  const [extfs, setExtFS] = useExtFS();
  const handleClick = () => {
    if (!item.available) return;
    if (item.fileType === "D") {
      const state = newExtFSState(extfs, item);
      setExtFS(state);
      return;
    }
  };
  return (
    <ExtFSItem
      style={style}
      primary={item.name}
      secondary={item.updatedAt}
      avatarIcon={avatarIcon}
      onClick={handleClick}
      disabled={!item.available}
      extIcon={<CloudIcon fontSize="small" />}
    >
      <ExtFSItemTag
        to={`${ExtFSRemoteFileTagRoutePath}/${item.id}`}
        disabled={!item.available}
        quantity={item.tagQuantity}
        pendingQuantity={item.pendingTagQuantity}
      />
    </ExtFSItem>
  );
};
export const RemoteFileMore = () => {
  return (
    <More>
      <MoreHelpItem />
    </More>
  );
};
