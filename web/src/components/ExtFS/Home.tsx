import type { ExtFSItemRecord } from "./Item";
import { ExtFSItem, ExtFSItems, ExtFSItemTag, useExtFSItem } from "./Item";
import { More, MoreHelpItem } from "./More";
import { newExtFSState as newExtFSStateWithNodeItem } from "./NodeItem";
import {
  ExtFSRemoteState,
  newExtFSState as newExtFSStateWithRemote,
} from "./Remote";
import { useExtFS } from "./State";

import type { ExtFSRemoteNode } from "../../API";
import { useAPI } from "../../API";
import { AppNodeIcon } from "../AppNodes";
import { APP_SETTINGS_QUERY_KEY } from "../AppSettings";

import { useQuery } from "@tanstack/react-query";
import { useMemo } from "react";

import CloudIcon from "@mui/icons-material/Cloud";

const REMOTE_NODES_QUERY_KEY = ["extfs-remote-nodes"];
const ExtFSHomeMode = "H";

export const ExtFSHomeState = {
  mode: ExtFSHomeMode,
  queryKeyList: [APP_SETTINGS_QUERY_KEY, REMOTE_NODES_QUERY_KEY],
};

type ExtFSNode = {
  name: string;
  peerId: string;
};
export const HomeItems = () => {
  const [extfs, _] = useExtFS();
  const api = useAPI();
  const { data: settings, isFetching: isSettingsFetching } = useQuery({
    queryKey: APP_SETTINGS_QUERY_KEY,
    queryFn: async () => await api?.selectAllAppSettings(),
    enabled: extfs.mode === ExtFSHomeMode,
  });
  const nodeItem = useMemo(() => {
    if (!settings) return;
    return {
      name: settings.name,
      peerId: settings.peerId,
    };
  }, [settings]);

  const { data: remotes, isFetching: isRemotesFetching } = useQuery({
    queryKey: REMOTE_NODES_QUERY_KEY,
    queryFn: async () => await api?.selectAllExtFSRemoteNodes(),
    enabled: extfs.mode === ExtFSHomeMode,
  });

  const items = useMemo(() => {
    const items_: Array<ExtFSNode | ExtFSRemoteNode> =
      remotes && remotes.length > 0 ? remotes : [];
    if (nodeItem) {
      return [nodeItem, ...items_];
    }
    return items_;
  }, [nodeItem, remotes]);

  const isFetching = useMemo(
    () => isSettingsFetching || isRemotesFetching,
    [isSettingsFetching, isRemotesFetching]
  );

  return (
    <ExtFSItems items={items} isFetching={isFetching}>
      <HomeItem />
    </ExtFSItems>
  );
};

const ExtFSNodeTagRoutePath = "/extfs/node-tags"; // TODO: not implement
const HomeItem = () => {
  const { style, item }: ExtFSItemRecord<ExtFSNode | ExtFSRemoteNode> =
    useExtFSItem();

  const [extfs, setExtFS] = useExtFS();

  const remoteNode = useMemo(
    () =>
      (item as ExtFSRemoteNode).updatedAt !== void 0
        ? (item as ExtFSRemoteNode)
        : void 0,
    [item]
  );
  if (remoteNode !== void 0) {
    const handleRemoteClick = () => {
      const state = newExtFSStateWithRemote(extfs, remoteNode);
      setExtFS(state);
    };

    return (
      <ExtFSItem
        style={style}
        primary={remoteNode.name}
        secondary={remoteNode.updatedAt}
        avatarIcon={
          <AppNodeIcon
            fontSize="large"
            color={remoteNode.available ? "primary" : "disabled"}
          />
        }
        extIcon={<CloudIcon fontSize="small" />}
        disabled={!remoteNode.available}
        onClick={handleRemoteClick}
      >
        <ExtFSItemTag
          to={`${ExtFSNodeTagRoutePath}/${remoteNode.peerId}`}
          disabled={!remoteNode.available}
          quantity={remoteNode.tagQuantity}
          pendingQuantity={remoteNode.pendingTagQuantity}
        />
      </ExtFSItem>
    );
  }

  const localNode = item as ExtFSNode;

  const handleLocalClick = () => {
    const state = newExtFSStateWithNodeItem(extfs, localNode.name);
    setExtFS(state);
  };

  return (
    <ExtFSItem
      style={style}
      primary={localNode.name}
      secondary="-"
      avatarIcon={<AppNodeIcon fontSize="large" color="primary" />}
      onClick={handleLocalClick}
    ></ExtFSItem>
  );
};

export const HomeMore = () => {
  return (
    <More>
      <MoreHelpItem />
    </More>
  );
};
