import type { ExtFSNodeItem, ExtFSSearchFile } from "./api";
import { selectAllExtFSNodeItems } from "./api";
import { generateEditPath, generateShowPath } from "../Route/utils";
import { ExtFSBrowseFilePath } from "../ExtFSBrowseFile/Route";
import { ExtFSNodeItemPath } from "../ExtFSNodeItem/Route";
import type { ExtFSItemRecord } from "./Item";
import {
  ExtFSItem,
  ExtFSItemOpen,
  ExtFSItems,
  ExtFSItemSettings,
  useExtFSItem,
} from "./Item";
import { newExtFSState as newExtFSStateWithNodeFile } from "./NodeFile";
import type { ExtFSState } from "./State";
import { useExtFS } from "./State";

import { useQuery } from "@tanstack/react-query";
import { useMemo } from "react";

import FolderIcon from "@mui/icons-material/Folder";
import InsertDriveFileIcon from "@mui/icons-material/InsertDriveFile";

import { useNavigate } from "../Route/Router";

export const ExtFSNodeMode = "N";
export const ExtFSNodeQueryKey = ["extfs-node-items"];

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

type NewExtFSStateReferOpts = Pick<ExtFSSearchFile, "name">;
export function newExtFSStateWithRefer(
  extfs: ExtFSState,
  opts: NewExtFSStateReferOpts
): ExtFSState {
  return newExtFSState(extfs, opts.name);
}

export const NodeItems = () => {
  const [extfs, _] = useExtFS();

  const {
    data: items,
    isFetching,
    error,
  } = useQuery({
    queryKey: ExtFSNodeQueryKey,
    queryFn: async () => await selectAllExtFSNodeItems(),
    enabled: extfs.mode === ExtFSNodeMode,
  });
  return (
    <ExtFSItems
      items={items || []}
      isFetching={isFetching}
      error={error || void 0}
    >
      <NodeItem />
    </ExtFSItems>
  );
};

export const NodeItem = () => {
  const { style, item }: ExtFSItemRecord<ExtFSNodeItem> = useExtFSItem();
  const navigate = useNavigate();
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

    navigate(generateShowPath(ExtFSNodeItemPath, item.id));
  };

  const settingsUrl = useMemo(
    () => generateEditPath(ExtFSNodeItemPath, item.id),
    [item.id]
  );

  const openUrl = useMemo(() => {
    const searchParams = new URLSearchParams({
      itemId: item.id.toString(),
    });
    return `${ExtFSBrowseFilePath}?${searchParams.toString()}`;
  }, [item.id]);

  const openHidden = useMemo(() => {
    return item.fileType !== "F";
  }, [item.fileType]);

  return (
    <ExtFSItem
      style={style}
      primary={item.name}
      secondary={item.updatedAt}
      avatarIcon={avatarIcon}
      onClick={handleClick}
      disabled={!item.available}
    >
      <ExtFSItemOpen
        to={openUrl}
        disabled={!item.available}
        hidden={openHidden}
      />
      <ExtFSItemSettings to={settingsUrl} />
    </ExtFSItem>
  );
};
