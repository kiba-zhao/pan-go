import type { ExtFSRemoteItem } from "../../API";
import { useAPI } from "../../API";
import type { ExtFSItemRecord } from "./Item";
import { ExtFSItem, ExtFSItems, ExtFSItemTag, useExtFSItem } from "./Item";
import { More, MoreHelpItem } from "./More";
import { ExtFSRemoteFileState } from "./RemoteFile";
import type { ExtFSSingleState } from "./State";
import { useExtFS } from "./State";

import { useQuery } from "@tanstack/react-query";

import FolderIcon from "@mui/icons-material/Folder";
import InsertDriveFileIcon from "@mui/icons-material/InsertDriveFile";

import { useMemo } from "react";

const ExtFSRemoteItemTagRoutePath = "/extfs/remote-item-tags";
const ExtFSRemoteMode = "R";
const ExtFSRemoteQueryKey = ["extfs-remote-items"];
export const ExtFSRemoteState = {
  mode: ExtFSRemoteMode,
  queryKeyList: [ExtFSRemoteQueryKey],
};

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
      const { parentItems } = extfs;
      const fileState = {
        ...ExtFSRemoteFileState,
        peerId: item.peerId,
        itemId: item.itemId,
      };
      setExtFS({
        ...fileState,
        parentItems: [...parentItems, { name: item.name, state: fileState }],
      });
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
    >
      <ExtFSItemTag
        to={`${ExtFSRemoteItemTagRoutePath}/${item.id}`}
        disabled={!item.available}
        quantity={item.tagQuantity}
        pendingQuantity={item.pendingTagQuantity}
      />
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
