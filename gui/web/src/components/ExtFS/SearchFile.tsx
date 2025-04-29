import { ListItems } from "../Common/Item";
import {
  ExtFSItem,
  ExtFSItemOpen,
  ExtFSItemRecord,
  ExtFSItemSettings,
  useExtFSItem,
} from "./Item";

import { newExtFSState as newExtFSStateWithNodeFile } from "./NodeFile";
import { newExtFSState as newExtFSStateWithRemoteFile } from "./RemoteFile";
import type { SearchFileStore } from "./SearchFileStore";
import { newSearchFileStore } from "./SearchFileStore";
import { ExtFSSingleState, ExtFSState, useExtFS } from "./State";

import { generateEditPath } from "../Route/utils";
import { ExtFSNodeItemPath } from "../ExtFSNodeItem/Route";
import { ExtFSBrowseFilePath } from "../ExtFSBrowseFile/Route";
import type { ExtFSSearchFile, ExtFSSearchItem } from "./api";
import { useTranslation } from "../i18n/Context";

import CloseIcon from "@mui/icons-material/Close";
import CloudIcon from "@mui/icons-material/Cloud";
import FolderIcon from "@mui/icons-material/Folder";
import InsertDriveFileIcon from "@mui/icons-material/InsertDriveFile";
import RefreshIcon from "@mui/icons-material/Refresh";
import SearchOffIcon from "@mui/icons-material/SearchOff";
import IconButton from "@mui/material/IconButton";
import LinearProgress from "@mui/material/LinearProgress";
import Link from "@mui/material/Link";
import MenuItem, { MenuItemOwnProps } from "@mui/material/MenuItem";
import Alert from "@mui/material/Alert";

import { Fragment, useMemo, useSyncExternalStore } from "react";

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

export type ExtFSSearchFileStateOpts = Pick<ExtFSSearchItem, "query">;
export function newExtFSState(
  extfs: ExtFSState,
  opts: ExtFSSearchFileStateOpts
): ExtFSState {
  const { parentItems, ...state_ } = extfs;
  const { snapshot } = state_ as ExtFSSearchFileSingleState;
  const state = {
    ...ExtFSSearchFileState,
    snapshot: snapshot || extfs,
    store: newSearchFileStore(opts.query),
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

type ExtFSSearchFileData = { peerId: string } & ExtFSSearchFile;

export const SearchFiles = () => {
  const [{ parentItems, ...state }, _] = useExtFS();
  const { store } = state as ExtFSSearchFileSingleState;
  const { files, isComplete, errs } = useSyncExternalStore(
    store.subscribe,
    store.getSnapshot
  );
  const { t } = useTranslation();

  const [error] = useMemo(() => {
    if (!errs) return [void 0, []];
    const details: Array<[string, Error]> = [];
    let error: Error | undefined;
    for (const [peerId, err] of Object.entries(errs)) {
      if (peerId.length <= 0) {
        error = err;
        continue;
      }
      details.push([peerId, err]);
    }
    if (!error && details.length > 0) {
      error = new Error();
      error.name = "SearchError";
    }
    return [error, details];
  }, [errs]);
  return (
    <Fragment>
      <LinearProgress sx={{ visibility: isComplete ? "hidden" : "visible" }} />
      {isComplete && error ? (
        <Alert severity="error">
          {t(`errors.${error.name}`, { _: error.message })}
        </Alert>
      ) : (
        void 0
      )}
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
    if (
      item?.peerId === void 0 &&
      (item?.filePath === void 0 || item?.filePath === "")
    )
      return generateEditPath(ExtFSNodeItemPath, item.itemId);
    return "";
  }, [item?.peerId, item?.filePath, item?.itemId]);

  const [extfs, setExtFS] = useExtFS();

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

  const openUrl = useMemo(() => {
    const searchParams = new URLSearchParams({
      itemId: item.itemId.toString(),
    });
    if (item.peerId !== void 0 && item.peerId.length > 0) {
      searchParams.set("peerId", item.peerId);
    }
    if (item.filePath !== void 0 && item.filePath.length > 0) {
      searchParams.set("filePath", item.filePath);
    }
    return `${ExtFSBrowseFilePath}?${searchParams.toString()}`;
  }, [item.peerId, item.itemId, item.filePath]);

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
      extIcon={item.peerId !== void 0 ? <CloudIcon fontSize="small" /> : null}
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
      <SearchOffIcon fontSize="medium" />
    </Link>
  );
};

export const SearchNavigationMenuRoot = ({
  sx,
}: {
  sx: MenuItemOwnProps["sx"];
  anchorElWidth?: number;
}) => {
  const { t } = useTranslation();
  const [extfs, setExtFS] = useExtFS();
  const handleClick = () => {
    const state = restoreExtFSState(extfs);
    setExtFS(state);
  };
  return (
    <MenuItem onClick={handleClick} sx={sx}>
      {t("extfs.search.search")}
    </MenuItem>
  );
};

export const SearchFileRefresh = () => {
  const [extfs, _] = useExtFS();
  const { parentItems, ...state } = extfs || { store: {} };
  const fileState = state as ExtFSSearchFileSingleState;

  const { store } = fileState;
  const { isComplete } = useSyncExternalStore(
    store.subscribe,
    store.getSnapshot
  );
  const handleClick = () => {
    if (isComplete) {
      store.refresh();
    } else {
      store.abort("cancelled");
    }
  };
  return (
    <IconButton onClick={handleClick}>
      {isComplete ? <RefreshIcon /> : <CloseIcon />}
    </IconButton>
  );
};
