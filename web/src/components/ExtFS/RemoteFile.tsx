import type { ExtFSItemRecord } from "./Item";
import { ExtFSItem, ExtFSItems, ExtFSItemTag, useExtFSItem } from "./Item";
import { More, MoreHelpItem } from "./More";
import type { ExtFSSingleState } from "./State";
import { useExtFS } from "./State";

import { useQuery } from "@tanstack/react-query";
import type { ExtFSRemoteFile } from "../../API";
import { useAPI } from "../../API";

import { useMemo } from "react";

import FolderIcon from "@mui/icons-material/Folder";
import InsertDriveFileIcon from "@mui/icons-material/InsertDriveFile";

const ExtFSRemoteFileTagRoutePath = "/extfs/remote-file-tags";
const ExtFSRemoteFileMode = "RF";
const ExtFSRemoteFileQueryKey = ["extfs-remote-files"];
export const ExtFSRemoteFileState = {
  mode: ExtFSRemoteFileMode,
  queryKeyList: [ExtFSRemoteFileQueryKey],
};

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
      const { parentItems, ...state } = extfs;
      const fileState = {
        ...state,
        parentPath: item.filePath,
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
