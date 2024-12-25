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
import type { ExtFSSingleState, ExtFSState } from "./State";
import { useExtFS } from "./State";

import type { ExtFSNodeFile, ExtFSSearchFile, API } from "../../API";
import { useAPI } from "../../API";

import { useQuery } from "@tanstack/react-query";

import FolderIcon from "@mui/icons-material/Folder";
import InsertDriveFileIcon from "@mui/icons-material/InsertDriveFile";

import { useMemo } from "react";

const ExtFSNodeFileRoutePath = "/extfs/node-files";
const ExtFSNodeFileTagRoutePath = "/extfs/node-file-tags";

export const ExtFSNodeFileMode = "NF";
const ExtFSNodeFileQueryKey = ["extfs-node-files"];
const ExtFSNodeFileState = {
  mode: ExtFSNodeFileMode,
  queryKeyList: [ExtFSNodeFileQueryKey],
};

type NewExtFSStateOpts = {
  filePath?: ExtFSNodeFile["filePath"];
} & Pick<ExtFSNodeFile, "itemId" | "name">;
export function newExtFSState(
  extfs: ExtFSState,
  opts: NewExtFSStateOpts
): ExtFSState {
  const { parentItems } = extfs;
  const fileState = {
    ...ExtFSNodeFileState,
    itemId: opts.itemId,
    parentPath: opts.filePath,
  };
  return {
    ...fileState,
    parentItems: [...parentItems, { name: opts.name, state: fileState }],
  } as ExtFSState;
}

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
  const { style, item }: ExtFSItemRecord<ExtFSNodeFile> = useExtFSItem();

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
    >
      <ExtFSItemTag
        to={`${ExtFSNodeFileTagRoutePath}/${item.id}`}
        disabled={!item.available}
        quantity={item.tagQuantity}
        pendingQuantity={item.pendingTagQuantity}
      />
      <ExtFSItemSettings to={`${ExtFSNodeFileRoutePath}/${item.id}`} />
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
