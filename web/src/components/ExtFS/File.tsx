import type { ExtFSItemRecord } from "./Item";
import {
  ExtFSItem,
  ExtFSItems,
  ExtFSItemSettings,
  ExtFSItemTag,
  useExtFSItem,
} from "./Item";
import { More, MoreHelpItem } from "./More";
import type { ExtFSSingleState } from "./State";
import { useExtFS } from "./State";

import type { ExtFSFileItem } from "../../API";
import { useAPI } from "../../API";

import { useQuery } from "@tanstack/react-query";

import FolderIcon from "@mui/icons-material/Folder";
import InsertDriveFileIcon from "@mui/icons-material/InsertDriveFile";

import { useMemo } from "react";

const ExtFSNodeItemRoutePath = "/extfs/file-items";
const ExtFSFileItemTagRoutePath = "/extfs/file-item-tags";

const ExtFSFileMode = "F";
const ExtFSFileQueryKey = ["extfs-file-items"];
export const ExtFSFileState = {
  mode: ExtFSFileMode,
  queryKeyList: [ExtFSFileQueryKey],
};

export type ExtFSFileSingleState = {
  itemId: number;
  parentPath?: string;
} & ExtFSSingleState;
export const FileItems = () => {
  const [{ parentItems, ...state }, _] = useExtFS();
  const { itemId, parentPath } = state as ExtFSFileSingleState;
  const api = useAPI();
  const { data: items, isFetching } = useQuery({
    queryKey: [...ExtFSFileQueryKey, { itemId, parentPath }],
    queryFn: async () =>
      await api?.searchExtFSFileItems({ itemId, parentPath }),
    enabled: state.mode === ExtFSFileMode,
  });

  return (
    <ExtFSItems items={items || []} isFetching={isFetching}>
      <FileItem />
    </ExtFSItems>
  );
};

export const FileItem = () => {
  const { style, data }: ExtFSItemRecord<ExtFSFileItem> = useExtFSItem();

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
        to={`${ExtFSFileItemTagRoutePath}/${data.id}`}
        disabled={!data.available}
        quantity={data.tagQuantity}
        pendingQuantity={data.pendingTagQuantity}
      />
      <ExtFSItemSettings to={`${ExtFSNodeItemRoutePath}/${data.id}`} />
    </ExtFSItem>
  );
};

export const FileMore = () => {
  return (
    <More>
      <MoreHelpItem />
    </More>
  );
};
