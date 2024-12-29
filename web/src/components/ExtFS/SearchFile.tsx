import {
  ExtFSItem,
  ExtFSItemRecord,
  ExtFSItems,
  ExtFSItemSettings,
  useExtFSItem,
} from "./Item";
import { More, MoreHelpItem } from "./More";
import { newExtFSState as newExtFSStateWithNodeFile } from "./NodeFile";
import { newItemSettingsUrl as newItemSettingsUrlWithNodeItem } from "./NodeItem";
import { newExtFSState as newExtFSStateWithRemoteFile } from "./RemoteFile";
import { ExtFSSingleState, ExtFSState, useExtFS } from "./State";

import { useQuery } from "@tanstack/react-query";
import type { ExtFSSearchFile, ExtFSSearchItem } from "../../API";
import { useAPI } from "../../API";

import FolderIcon from "@mui/icons-material/Folder";
import InsertDriveFileIcon from "@mui/icons-material/InsertDriveFile";
import SearchOffIcon from "@mui/icons-material/SearchOff";
import Link from "@mui/material/Link";
import MenuItem, { MenuItemOwnProps } from "@mui/material/MenuItem";

import { useMemo } from "react";

export const ExtFSSearchFileMode = "SF";
const ExtFSSearchFileQueryKey = ["extfs-search-files"];
const ExtFSSearchFileState = {
  mode: ExtFSSearchFileMode,
  queryKeyList: [ExtFSSearchFileQueryKey],
};

export type ExtFSSearchFileSingleState = {
  snapshot: ExtFSState;
} & ExtFSSingleState &
  Pick<ExtFSSearchFile, "searchId">;

export type ExtFSSearchFileStateOpts = Pick<ExtFSSearchItem, "id" | "query">;
export function newExtFSState(
  extfs: ExtFSState,
  opts: ExtFSSearchFileStateOpts
): ExtFSState {
  const { parentItems, ...state_ } = extfs;
  const { snapshot } = state_ as ExtFSSearchFileSingleState;
  const state = {
    ...ExtFSSearchFileState,
    searchId: opts.id,
    snapshot: snapshot || extfs,
  };
  return {
    ...state,
    parentItems: [{ name: opts.query, state: state }],
  } as ExtFSState;
}

function restoreExtFSState(extfs: ExtFSState): ExtFSState {
  const { parentItems, ...state } = extfs;
  if (state.mode === ExtFSSearchFileMode)
    return (state as ExtFSSearchFileSingleState).snapshot;
  const rootItem = parentItems.at(0);
  if (rootItem !== void 0 && rootItem.state.mode === ExtFSSearchFileMode)
    return (rootItem.state as ExtFSSearchFileSingleState).snapshot;
  return extfs;
}

export const SearchFiles = () => {
  const [{ parentItems, ...state }, _] = useExtFS();
  const { searchId } = state as ExtFSSearchFileSingleState;
  const api = useAPI();
  const { data: items, isFetching } = useQuery({
    queryKey: [...ExtFSSearchFileQueryKey, searchId],
    queryFn: async () => await api?.searchExtFSSearchFiles({ searchId }),
    enabled: state.mode === ExtFSSearchFileMode,
  });

  return (
    <ExtFSItems items={items || []} isFetching={isFetching}>
      <SearchFile />
    </ExtFSItems>
  );
};

export const SearchFile = () => {
  const { style, item }: ExtFSItemRecord<ExtFSSearchFile> = useExtFSItem();

  const avatarIcon = useMemo(() => {
    if (item?.fileType === "D")
      return <FolderIcon color={item.available ? "primary" : "disabled"} />;
    if (item?.fileType === "F")
      return (
        <InsertDriveFileIcon color={item.available ? "action" : "disabled"} />
      );
  }, [item?.fileType]);

  const settingsUrl = useMemo(() => {
    // TODO: redirect to settings view
    if (item === void 0) return "";
    if (item?.peerId === void 0 && item?.filePath === void 0)
      return newItemSettingsUrlWithNodeItem(item.itemId);
    return "";
  }, [item?.peerId, item?.filePath, item?.itemId]);

  const [extfs, setExtFS] = useExtFS();

  const api = useAPI();
  const handleClick = async () => {
    if (item.fileType === "D") {
      let state: ExtFSState | undefined;
      if (item?.peerId === void 0 && item?.filePath === void 0) {
        state = newExtFSStateWithNodeFile(extfs, {
          name: item.name,
          itemId: item.itemId,
        });
      } else if (item?.peerId === void 0 && item?.filePath !== void 0) {
        state = newExtFSStateWithNodeFile(extfs, item);
      } else if (item?.peerId !== void 0 && item?.filePath === void 0) {
        state = newExtFSStateWithRemoteFile(extfs, item);
      } else if (item?.peerId !== void 0 && item?.filePath !== void 0) {
        state = newExtFSStateWithRemoteFile(extfs, item);
      }

      if (state !== void 0) setExtFS(state);
    }
    // TODO: open a file
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
      <ExtFSItemSettings to={settingsUrl} />
    </ExtFSItem>
  );
};

export const SearchFileMore = () => {
  return (
    <More>
      <MoreHelpItem />
    </More>
  );
};

export const SearchNavigationBreadcrumbRoot = () => {
  const [extfs, setExtFS] = useExtFS();
  const handleClick = () => {
    const state = restoreExtFSState(extfs);
    setExtFS(state);
  };
  return (
    <Link
      underline="hover"
      sx={{ display: "flex", alignItems: "center", cursor: "pointer" }}
      color="inherit"
      onClick={handleClick}
    >
      <SearchOffIcon fontSize="inherit" />
    </Link>
  );
};

export const SearchNavigationMenuRoot = ({
  sx,
}: {
  sx: MenuItemOwnProps["sx"];
  anchorElWidth?: number;
}) => {
  const [extfs, setExtFS] = useExtFS();
  const handleClick = () => {
    const state = restoreExtFSState(extfs);
    setExtFS(state);
  };
  return (
    <MenuItem onClick={handleClick} sx={sx}>
      Search Exit
    </MenuItem>
  );
};
