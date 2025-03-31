import { RoutePath as ExtFSBrowseRoutePath } from "../ExtFSBrowseFile";
import type { ExtFSItemRecord } from "./Item";
import { ExtFSItem, ExtFSItemOpen, ExtFSItems, useExtFSItem } from "./Item";
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

export const ExtFSRemoteFileMode = "RF";
export const ExtFSRemoteFileQueryKey = ["extfs-remote-files"];
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
  const {
    data: items,
    isFetching,
    error,
  } = useQuery({
    queryKey: [...ExtFSRemoteFileQueryKey, { peerId, itemId, parentPath }],
    queryFn: async () =>
      await api?.searchExtFSRemoteFiles(peerId, itemId, { parentPath }),
    enabled: !!api && state.mode === ExtFSRemoteFileMode,
  });

  return (
    <ExtFSItems
      items={items || []}
      isFetching={isFetching}
      error={error || void 0}
    >
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

  const openUrl = useMemo(() => {
    const searchParams = new URLSearchParams({
      peerId: item.peerId,
      itemId: item.itemId.toString(),
      filePath: item.filePath,
    });
    return `${ExtFSBrowseRoutePath}?${searchParams.toString()}`;
  }, [item.peerId, item.itemId, item.filePath]);

  const openHidden = useMemo(() => {
    return item.fileType !== "F";
  }, [item.fileType]);

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
      <ExtFSItemOpen
        to={openUrl}
        disabled={!item.available}
        hidden={openHidden}
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
