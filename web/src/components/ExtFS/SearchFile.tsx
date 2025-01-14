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
import { REMOTE_NODES_QUERY_KEY } from "./Home";

import type { ExtFSSearchFile, ExtFSSearchItem } from "../../api";
import { useAPI } from "../API";

import FolderIcon from "@mui/icons-material/Folder";
import InsertDriveFileIcon from "@mui/icons-material/InsertDriveFile";
import SearchOffIcon from "@mui/icons-material/SearchOff";
import LinearProgress from "@mui/material/LinearProgress";
import Link from "@mui/material/Link";
import MenuItem, { MenuItemOwnProps } from "@mui/material/MenuItem";

import { Fragment, useEffect, useMemo, useState, useRef } from "react";

import { useQuery } from "@tanstack/react-query";

export const ExtFSSearchFileMode = "SF";
const ExtFSSearchFileQueryKey = ["extfs-search-files"];
const ExtFSSearchFileState = {
  mode: ExtFSSearchFileMode,
  queryKeyList: [ExtFSSearchFileQueryKey],
};

export type ExtFSSearchFileSingleState = {
  snapshot: ExtFSState;
} & ExtFSSingleState &
  Pick<ExtFSSearchItem, "query">;

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
  };
  return {
    ...state,
    query: opts.query,
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
  const { query } = state as ExtFSSearchFileSingleState;

  const [itemsPages, setItemsPages] = useState<ExtFSSearchFileData[][]>([]);

  useEffect(() => {
    if (query === void 0 || query.length <= 0) return;
  }, [query]);

  // const api = useAPI();

  // const { data: remotes, isFetching } = useQuery({
  //   queryKey: REMOTE_NODES_QUERY_KEY,
  //   queryFn: async () => await api?.selectAllExtFSRemoteNodes(),
  //   enabled: state.mode === ExtFSSearchFileMode && !!api,
  // });

  // const {
  //   data,
  //   isFetching,
  //   error,
  //   fetchNextPage,
  //   hasNextPage,
  //   isFetchingNextPage,
  // } = useInfiniteQuery({
  //   queryKey: [...ExtFSSearchFileQueryKey, searchId],
  //   queryFn: async ({ pageParam }) =>
  //     await api?.searchExtFSSearchFileResults({
  //       searchId,
  //       _start: pageParam,
  //       _end: pageParam + 1000,
  //     }),
  //   initialPageParam: 0,
  //   getNextPageParam: (lastPage, allPages, lastPageParam) => {
  //     if (lastPage === void 0 || allPages.length <= 0) return lastPageParam;
  //     const [total, _] = lastPage;
  //     const count = allPages.reduce(
  //       (counter, page) => counter + (page !== void 0 ? page[1].length : 0),
  //       0
  //     );
  //     if (total < 0 || count < total) return count;
  //     return void 0;
  //   },
  //   enabled: state.mode === ExtFSSearchFileMode,
  // });

  // const intervalRef = useRef<ReturnType<typeof setInterval> | null>(null);
  // const cleanInterval = useCallback(() => {
  //   if (intervalRef.current === null) return;
  //   clearInterval(intervalRef.current);
  //   intervalRef.current = null;
  // }, []);
  // useEffect(() => {
  //   if (hasNextPage && intervalRef.current === null) {
  //     intervalRef.current = setInterval(() => {
  //       if (hasNextPage) {
  //         fetchNextPage();
  //       } else {
  //         cleanInterval();
  //       }
  //     }, 2000);
  //   }
  //   return cleanInterval;
  // }, [hasNextPage]);

  const items = useMemo(() => {
    if (itemsPages.length <= 0) return [];
    return itemsPages.flat();
  }, [itemsPages]);

  return (
    <Fragment>
      <LinearProgress
        sx={{ visibility: hasNextPage || isFetching ? "visible" : "hidden" }}
      />
      <ListItems items={items} itemSize={68}>
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
