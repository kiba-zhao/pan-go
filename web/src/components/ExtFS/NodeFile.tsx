import { ExtFSNodeItemRoutePath } from "../ExtFSNodeItem";
import type { ExtFSItemRecord } from "./Item";
import {
  ExtFSItem,
  ExtFSItems,
  ExtFSItemSettings,
  ExtFSItemTag,
  useExtFSItem,
} from "./Item";
import { More, MoreHelpItem, MoreSettingsItem } from "./More";
import type { ExtFSSingleState } from "./State";
import { useExtFS } from "./State";

import type { ExtFSNodeFile } from "../../API";
import { useAPI } from "../../API";

import { useQuery } from "@tanstack/react-query";

import FolderIcon from "@mui/icons-material/Folder";
import InsertDriveFileIcon from "@mui/icons-material/InsertDriveFile";

import { useMemo } from "react";

const ExtFSNodeFileRoutePath = "/extfs/node-files";
const ExtFSNodeFileTagRoutePath = "/extfs/node-file-tags";

const ExtFSNodeFileMode = "F";
const ExtFSNodeFileQueryKey = ["extfs-node-files"];
export const ExtFSNodeFileState = {
  mode: ExtFSNodeFileMode,
  queryKeyList: [ExtFSNodeFileQueryKey],
};

export type ExtFSNodeFileSingleState = {
  itemId: number;
  parentPath?: string;
} & ExtFSSingleState;
export const NodeFiles = () => {
  const [{ parentItems, ...state }, _] = useExtFS();
  const { itemId, parentPath } = state as ExtFSNodeFileSingleState;
  const api = useAPI();
  const { data: items, isFetching } = useQuery({
    queryKey: [...ExtFSNodeFileQueryKey, { itemId, parentPath }],
    queryFn: async () =>
      await api?.searchExtFSNodeFiles({ itemId, parentPath }),
    enabled: state.mode === ExtFSNodeFileMode,
  });

  return (
    <ExtFSItems items={items || []} isFetching={isFetching}>
      <NodeFile />
    </ExtFSItems>
  );
};

export const NodeFile = () => {
  const { style, data }: ExtFSItemRecord<ExtFSNodeFile> = useExtFSItem();

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
        to={`${ExtFSNodeFileTagRoutePath}/${data.id}`}
        disabled={!data.available}
        quantity={data.tagQuantity}
        pendingQuantity={data.pendingTagQuantity}
      />
      <ExtFSItemSettings to={`${ExtFSNodeFileRoutePath}/${data.id}`} />
    </ExtFSItem>
  );
};

export const NodeFileMore = () => {
  return (
    <More>
      <NodeFileSettingsMore />
      <MoreHelpItem />
    </More>
  );
};

const NodeFileSettingsMore = () => {
  const [{ parentItems, ...state }, _] = useExtFS();
  const { itemId, parentPath } = state as ExtFSNodeFileSingleState;
  if (parentPath) return void 0;
  return <MoreSettingsItem to={`${ExtFSNodeItemRoutePath}/${itemId}`} />;
};
