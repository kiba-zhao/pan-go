import type { ExtFSRemoteItem } from "./api";
import { searchExtFSRemoteItems } from "./api";
import { ExtFSBrowseFilePath } from "../ExtFSBrowseFile/Route";
import type { ExtFSItemRecord } from "./Item";
import { ExtFSItem, ExtFSItemOpen, ExtFSItems, useExtFSItem } from "./Item";

import { newExtFSState as newExtFSStateWithRemoteFile } from "./RemoteFile";
import type { ExtFSSingleState, ExtFSState } from "./State";
import { useExtFS } from "./State";

import { useQuery } from "@tanstack/react-query";
import CloudIcon from "@mui/icons-material/Cloud";
import FolderIcon from "@mui/icons-material/Folder";
import InsertDriveFileIcon from "@mui/icons-material/InsertDriveFile";

import { useMemo } from "react";

export const ExtFSRemoteMode = "R";
export const ExtFSRemoteQueryKey = ["extfs-remote-items"];
const ExtFSRemoteState = {
  mode: ExtFSRemoteMode,
  queryKeyList: [ExtFSRemoteQueryKey],
};

type NewExtFSStateOpts = Pick<ExtFSRemoteItem, "peerId" | "name">;
export function newExtFSState(
  extfs: ExtFSState,
  opts: NewExtFSStateOpts
): ExtFSState {
  const { parentItems } = extfs;
  const state = {
    ...ExtFSRemoteState,
    peerId: opts.peerId,
  };
  return {
    ...state,
    parentItems: [...parentItems, { name: opts.name, state }],
  } as ExtFSState;
}

export type ExtFSRemoteSingleState = {
  peerId: string;
} & ExtFSSingleState;
export const RemoteItems = () => {
  const [{ parentItems, ...state }, _] = useExtFS();
  const { peerId } = state as ExtFSRemoteSingleState;

  const {
    data: items,
    isFetching,
    error,
  } = useQuery({
    queryKey: [...ExtFSRemoteQueryKey, { peerId }],
    queryFn: async () => await searchExtFSRemoteItems(peerId),
    enabled: state.mode === ExtFSRemoteMode,
  });
  return (
    <ExtFSItems
      items={items || []}
      isFetching={isFetching}
      error={error || void 0}
    >
      <RemoteItem />
    </ExtFSItems>
  );
};

export const RemoteItem = () => {
  const { style, item }: ExtFSItemRecord<ExtFSRemoteItem> = useExtFSItem();
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
      const state = newExtFSStateWithRemoteFile(extfs, {
        peerId: item.peerId,
        itemId: item.itemId,
        name: item.name,
      });
      setExtFS(state);
      return;
    }
  };

  const openUrl = useMemo(() => {
    const searchParams = new URLSearchParams({
      peerId: item.peerId,
      itemId: item.itemId.toString(),
    });
    return `${ExtFSBrowseFilePath}?${searchParams.toString()}`;
  }, [item.peerId, item.itemId]);

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
