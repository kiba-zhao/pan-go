import {
  ExtFSItems,
  ExtFSItemRecord,
  useExtFSItem,
  ExtFSItem,
  ExtFSItemSettings,
} from "./Item";
import { ExtFSSingleState, useExtFS, ExtFSState } from "./State";
import { More, MoreHelpItem } from "./More";
import { newItemSettingsUrl as newItemSettingsUrlWithNodeItem } from "./NodeItem";
import { newExtFSState as newExtFSStateWithNodeFile } from "./NodeFile";
import { newExtFSState as newExtFSStateWithRemoteFile } from "./RemoteFile";

import type { ExtFSSearchFile } from "../../API";
import { useAPI } from "../../API";
import { useQuery } from "@tanstack/react-query";

import FolderIcon from "@mui/icons-material/Folder";
import InsertDriveFileIcon from "@mui/icons-material/InsertDriveFile";
import Link from "@mui/material/Link";
import SearchOffIcon from "@mui/icons-material/SearchOff";
import MenuItem, { MenuItemOwnProps } from "@mui/material/MenuItem";

import { useMemo } from "react";

export const ExtFSSearchFileMode = "SF";
const ExtFSSearchFileQueryKey = ["extfs-search-files"];
const ExtFSSearchFileState = {
  mode: ExtFSSearchFileMode,
  queryKeyList: [ExtFSSearchFileQueryKey],
};

export type ExtFSSearchFileSingleState = {
  q: string;
  snapshot: ExtFSState;
} & ExtFSSingleState;

export function newExtFSState(extfs: ExtFSState, query: string): ExtFSState {
  const { parentItems, ...state_ } = extfs;
  const { snapshot } = state_ as ExtFSSearchFileSingleState;
  const state = {
    ...ExtFSSearchFileState,
    q: query,
    snapshot: snapshot || extfs,
  };
  return {
    ...state,
    parentItems: [{ name: query, state: state }],
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
  const { q } = state as ExtFSSearchFileSingleState;
  const api = useAPI();
  const { data: items, isFetching } = useQuery({
    queryKey: [...ExtFSSearchFileQueryKey, q],
    queryFn: async () => await api?.searchExtFSSearchFiles({ q }),
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
  }, [item?.fileType, item?.referType]);

  const settingsUrl = useMemo(() => {
    // TODO: redirect to settings view
    if (item?.referType === "NI")
      return newItemSettingsUrlWithNodeItem(Number(item.referId));
    return "";
  }, [item?.referType, item?.referId]);

  const [extfs, setExtFS] = useExtFS();

  const api = useAPI();
  const handleClick = async () => {
    if (item.fileType === "D") {
      let state: ExtFSState | undefined;
      if (item?.referType === "NI") {
        state = newExtFSStateWithNodeFile(extfs, {
          name: item.name,
          itemId: Number(item.referId),
        });
      } else if (item?.referType === "NF") {
        const nodeFile = await api?.selectExtFSNodeFile(item.referId);
        if (nodeFile !== void 0)
          state = newExtFSStateWithNodeFile(extfs, {
            ...nodeFile,
            name: item.name,
          });
      } else if (item?.referType === "RI") {
        const remoteItem = await api?.selectExtFSRemoteItem(item.referId);
        if (remoteItem !== void 0)
          state = newExtFSStateWithRemoteFile(extfs, {
            ...remoteItem,
            name: item.name,
          });
      } else if (item?.referType === "RF") {
        const remoteFile = await api?.selectExtFSRemoteFile(item.referId);
        if (remoteFile !== void 0)
          state = newExtFSStateWithRemoteFile(extfs, {
            ...remoteFile,
            name: item.name,
          });
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
