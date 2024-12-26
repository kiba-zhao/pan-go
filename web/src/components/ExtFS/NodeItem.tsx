import type { ExtFSNodeItem, ExtFSSearchFile } from "../../API";
import { useAPI } from "../../API";
import { ExtFSNodeItemRoutePath } from "../ExtFSNodeItem";
import type { ExtFSItemRecord } from "./Item";
import {
  ExtFSItem,
  ExtFSItems,
  ExtFSItemSettings,
  ExtFSItemTag,
  useExtFSItem,
} from "./Item";
import { More, MoreHelpItem, MoreNewItem } from "./More";
import { newExtFSState as newExtFSStateWithNodeFile } from "./NodeFile";
import type { ExtFSState } from "./State";
import { useExtFS } from "./State";

import { useQuery } from "@tanstack/react-query";
import { useMemo } from "react";

import FolderIcon from "@mui/icons-material/Folder";
import InsertDriveFileIcon from "@mui/icons-material/InsertDriveFile";

const ExtFSNodeItemTagRoutePath = "/extfs/node-item-tags";
export const ExtFSNodeMode = "N";
const ExtFSNodeQueryKey = ["extfs-node-items"];

const ExtFSNodeState = {
  mode: ExtFSNodeMode,
  queryKeyList: [ExtFSNodeQueryKey],
};

export function newExtFSState(extfs: ExtFSState, name: string): ExtFSState {
  const { parentItems } = extfs;
  return {
    ...ExtFSNodeState,
    parentItems: [...parentItems, { name: name, state: ExtFSNodeState }],
  } as ExtFSState;
}

export function newItemSettingsUrl(id: ExtFSNodeItem["id"]): string {
  return `${ExtFSNodeItemRoutePath}/${id}`;
}

type NewExtFSStateReferOpts = Pick<ExtFSSearchFile, "name" | "referId">;
export function newExtFSStateWithRefer(
  extfs: ExtFSState,
  opts: NewExtFSStateReferOpts
): ExtFSState {
  return newExtFSState(extfs, opts.name);
}

export const NodeItems = () => {
  const [extfs, _] = useExtFS();

  const api = useAPI();
  const { data: items, isFetching } = useQuery({
    queryKey: ExtFSNodeQueryKey,
    queryFn: async () => await api?.selectAllExtFSNodeItems(),
    enabled: extfs.mode === ExtFSNodeMode,
  });
  return (
    <ExtFSItems items={items || []} isFetching={isFetching}>
      <NodeItem />
    </ExtFSItems>
  );
};

export const NodeItem = () => {
  const { style, item }: ExtFSItemRecord<ExtFSNodeItem> = useExtFSItem();
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
      const state = newExtFSStateWithNodeFile(extfs, {
        itemId: item.id,
        name: item.name,
      });
      setExtFS(state);
      return;
    }
  };

  const settingsUrl = useMemo(() => newItemSettingsUrl(item.id), [item.id]);
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
        to={`${ExtFSNodeItemTagRoutePath}/${item.id}`}
        disabled={!item.available}
        quantity={item.tagQuantity}
        pendingQuantity={item.pendingTagQuantity}
      />
      <ExtFSItemSettings to={settingsUrl} />
    </ExtFSItem>
  );
};

export const NodeMore = () => {
  return (
    <More>
      <NodeNewMore />
      <MoreHelpItem />
    </More>
  );
};

export const NodeNewMore = () => (
  <MoreNewItem to={`${ExtFSNodeItemRoutePath}/create`} />
);
