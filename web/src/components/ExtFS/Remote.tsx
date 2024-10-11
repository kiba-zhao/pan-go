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
  nodeId: string;
} & ExtFSSingleState;
export const RemoteItems = () => {
  const [{ parentItems, ...state }, _] = useExtFS();
  const { nodeId } = state as ExtFSRemoteSingleState;

  const api = useAPI();
  const { data: items, isFetching } = useQuery({
    queryKey: [...ExtFSRemoteQueryKey, { nodeId }],
    queryFn: async () => await api?.searchExtFSRemoteItems({ nodeId }),
    enabled: state.mode === ExtFSRemoteMode,
  });
  return (
    <ExtFSItems items={items || []} isFetching={isFetching}>
      <RemoteItem />
    </ExtFSItems>
  );
};

export const RemoteItem = () => {
  const { style, data }: ExtFSItemRecord<ExtFSRemoteItem> = useExtFSItem();
  const avatarIcon = useMemo(() => {
    if (data?.fileType === "D")
      return <FolderIcon color={data.available ? "primary" : "disabled"} />;
    if (data?.fileType === "F")
      return (
        <InsertDriveFileIcon color={data.available ? "action" : "disabled"} />
      );
  }, [data?.fileType]);

  const [extfs, setExtFS] = useExtFS();
  const handleClick = () => {
    if (!data.available) return;
    if (data.fileType === "D") {
      const { parentItems } = extfs;
      const fileState = {
        ...ExtFSRemoteFileState,
        nodeId: data.nodeId,
        itemId: data.remoteItemId,
      };
      setExtFS({
        ...fileState,
        parentItems: [...parentItems, { name: data.name, state: fileState }],
      });
      return;
    }
  };

  return (
    <ExtFSItem
      style={style}
      primary={data.name}
      secondary={data.updatedAt}
      avatarIcon={avatarIcon}
      onClick={handleClick}
      disabled={!data.available}
    >
      <ExtFSItemTag
        to={`${ExtFSRemoteItemTagRoutePath}/${data.id}`}
        disabled={!data.available}
        quantity={data.tagQuantity}
        pendingQuantity={data.pendingTagQuantity}
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
