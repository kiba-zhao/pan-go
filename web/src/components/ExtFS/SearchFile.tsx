import { ListItems } from "../List/Item";
import {
  ExtFSItem,
  ExtFSItemRecord,
  ExtFSItemSettings,
  useExtFSItem,
} from "./Item";
import { More, MoreHelpItem } from "./More";
import { newExtFSState as newExtFSStateWithNodeFile } from "./NodeFile";
import { newItemSettingsUrl as newItemSettingsUrlWithNodeItem } from "./NodeItem";
import { newExtFSState as newExtFSStateWithRemoteFile } from "./RemoteFile";
import { ExtFSSingleState, ExtFSState, useExtFS } from "./State";
import { newSearchFileStore } from "./SearchFileStore";
import type { SearchFileStore } from "./SearchFileStore";

import type { ExtFSSearchFile, ExtFSSearchItem } from "../../api";
import { useAPI } from "../API";
import type { API } from "../../api";

import FolderIcon from "@mui/icons-material/Folder";
import InsertDriveFileIcon from "@mui/icons-material/InsertDriveFile";
import SearchOffIcon from "@mui/icons-material/SearchOff";
import LinearProgress from "@mui/material/LinearProgress";
import Link from "@mui/material/Link";
import MenuItem, { MenuItemOwnProps } from "@mui/material/MenuItem";

import { Fragment, useEffect, useMemo, useSyncExternalStore } from "react";

export const ExtFSSearchFileMode = "SF";
const ExtFSSearchFileQueryKey = ["extfs-search-files"];
const ExtFSSearchFileState = {
  mode: ExtFSSearchFileMode,
  queryKeyList: [ExtFSSearchFileQueryKey],
};

export type ExtFSSearchFileSingleState = {
  snapshot: ExtFSState;
  store: SearchFileStore;
} & ExtFSSingleState;

export type ExtFSSearchFileStateOpts = { api: API } & Pick<
  ExtFSSearchItem,
  "query"
>;
export function newExtFSState(
  extfs: ExtFSState,
  opts: ExtFSSearchFileStateOpts
): ExtFSState {
  const { parentItems, ...state_ } = extfs;
  const { snapshot } = state_ as ExtFSSearchFileSingleState;
  const state = {
    ...ExtFSSearchFileState,
    snapshot: snapshot || extfs,
  };
  return {
    ...state,
    store: newSearchFileStore(opts.query, opts.api),
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

type ExtFSSearchFileData = { peerId: string } & ExtFSSearchFile;

export const SearchFiles = () => {
  const [{ parentItems, ...state }, _] = useExtFS();
  const { store } = state as ExtFSSearchFileSingleState;
  const { files, isComplete } = useSyncExternalStore(
    store.subscribe,
    store.getSnapshot
  );

  useEffect(() => {
    return () => {
      store.abort();
    };
  }, [store]);

  return (
    <Fragment>
      <LinearProgress sx={{ visibility: isComplete ? "hidden" : "visible" }} />
      <ListItems items={files} itemSize={68}>
        <SearchFile />
      </ListItems>
    </Fragment>
  );
};

export const SearchFile = () => {
  const { style, item }: ExtFSItemRecord<ExtFSSearchFileData> = useExtFSItem();

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
