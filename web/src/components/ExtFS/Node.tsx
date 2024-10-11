import type { ExtFSNodeItem } from "../../API";
import { useAPI } from "../../API";
import { ExtFSNodeItemRoutePath } from "../ExtFSNodeItem";
import { ExtFSFileState } from "./File";
import type { ExtFSItemRecord } from "./Item";
import {
  ExtFSItem,
  ExtFSItems,
  ExtFSItemSettings,
  ExtFSItemTag,
  useExtFSItem,
} from "./Item";
import { More, MoreHelpItem, MoreNewItem } from "./More";
import { useExtFS } from "./State";

import { useQuery } from "@tanstack/react-query";
import { useMemo } from "react";

import FolderIcon from "@mui/icons-material/Folder";
import InsertDriveFileIcon from "@mui/icons-material/InsertDriveFile";

const ExtFSNodeItemTagRoutePath = "/extfs/node-item-tags";
const ExtFSNodeMode = "N";
const ExtFSNodeQueryKey = ["extfs-node-items"];
export const ExtFSNodeState = {
  mode: ExtFSNodeMode,
  queryKeyList: [ExtFSNodeQueryKey],
};

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
  const { style, data }: ExtFSItemRecord<ExtFSNodeItem> = useExtFSItem();
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
        ...ExtFSFileState,
        itemId: data.id,
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
        to={`${ExtFSNodeItemTagRoutePath}/${data.id}`}
        disabled={!data.available}
        quantity={data.tagQuantity}
        pendingQuantity={data.pendingTagQuantity}
      />
      <ExtFSItemSettings to={`${ExtFSNodeItemRoutePath}/${data.id}`} />
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
