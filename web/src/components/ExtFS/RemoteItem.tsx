import type { ExtFSRemoteItem } from "../../api";
import { useAPI } from "../API";
import type { ExtFSItemRecord } from "./Item";
import { ExtFSItem, ExtFSItems, useExtFSItem, ExtFSItemOpen } from "./Item";
import { More, MoreHelpItem } from "./More";
import { newExtFSState as newExtFSStateWithRemoteFile } from "./RemoteFile";
import type { ExtFSSingleState, ExtFSState } from "./State";
import { useExtFS } from "./State";

import { useQuery } from "@tanstack/react-query";

import CloudIcon from "@mui/icons-material/Cloud";
import FolderIcon from "@mui/icons-material/Folder";
import InsertDriveFileIcon from "@mui/icons-material/InsertDriveFile";

import { useMemo } from "react";

export const ExtFSRemoteMode = "R";
const ExtFSRemoteQueryKey = ["extfs-remote-items"];
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

  const api = useAPI();
  const { data: items, isFetching } = useQuery({
    queryKey: [...ExtFSRemoteQueryKey, { peerId }],
    queryFn: async () => await api?.searchExtFSRemoteItems({ peerId }),
    enabled: state.mode === ExtFSRemoteMode,
  });
  return (
    <ExtFSItems items={items || []} isFetching={isFetching}>
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

  // TODO: implement
  const openUrl = useMemo(() => ``, [item.peerId, item.itemId]);
  //

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
      <ExtFSItemOpen to={openUrl} disabled={!item.available} />
    </ExtFSItem>
  );
};

export const RemoteMore = () => {
  return (
    <More>
      <MoreHelpItem />
    </More>
  );
};
