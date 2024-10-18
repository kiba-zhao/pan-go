import type { ExtFSItemRecord } from "./Item";
import { ExtFSItem, ExtFSItems, ExtFSItemTag, useExtFSItem } from "./Item";
import { More, MoreHelpItem } from "./More";
import type { ExtFSSingleState } from "./State";
import { useExtFS } from "./State";

import { useQuery } from "@tanstack/react-query";
import type { ExtFSRemoteFileItem } from "../../API";
import { useAPI } from "../../API";

import { useMemo } from "react";

import FolderIcon from "@mui/icons-material/Folder";
import InsertDriveFileIcon from "@mui/icons-material/InsertDriveFile";

const ExtFSRemoteFileItemTagRoutePath = "/extfs/remote-file-item-tags";
const ExtFSRemoteFileMode = "RF";
const ExtFSRemoteFileQueryKey = ["extfs-remote-file-items"];
export const ExtFSRemoteFileState = {
  mode: ExtFSRemoteFileMode,
  queryKeyList: [ExtFSRemoteFileQueryKey],
};

export type ExtFSRemoteFileSingleState = {
  nodeId: string;
  itemId: number;
  parentPath?: string;
} & ExtFSSingleState;
export const RemoteFileItems = () => {
  const [{ parentItems, ...state }, _] = useExtFS();
  const { nodeId, itemId, parentPath } = state as ExtFSRemoteFileSingleState;

  const api = useAPI();
  const { data: items, isFetching } = useQuery({
    queryKey: [...ExtFSRemoteFileQueryKey, { nodeId, itemId, parentPath }],
    queryFn: async () =>
      await api?.searchExtFSRemoteFileItems({ nodeId, itemId, parentPath }),
    enabled: !!api && state.mode === ExtFSRemoteFileMode,
  });

  return (
    <ExtFSItems items={items || []} isFetching={isFetching}>
      <RemoteFileItem />
    </ExtFSItems>
  );
};

export const RemoteFileItem = () => {
  const { style, data }: ExtFSItemRecord<ExtFSRemoteFileItem> = useExtFSItem();
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
      const { parentItems, ...state } = extfs;
      const fileState = {
        ...state,
        parentPath: data.filePath,
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
        to={`${ExtFSRemoteFileItemTagRoutePath}/${data.id}`}
        disabled={!data.available}
        quantity={data.tagQuantity}
        pendingQuantity={data.pendingTagQuantity}
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
