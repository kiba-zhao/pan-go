import type { ExtFSItemRecord } from "./Item";
import { ExtFSItem, ExtFSItems, ExtFSItemTag, useExtFSItem } from "./Item";
import { More, MoreHelpItem } from "./More";
import { ExtFSNodeState } from "./Node";
import { ExtFSRemoteState } from "./Remote";
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
  nodeId: string;
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
      nodeId: settings.nodeId,
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
  const { style, data }: ExtFSItemRecord<ExtFSNode | ExtFSRemoteNode> =
    useExtFSItem();

  const [extfs, setExtFS] = useExtFS();

  const remoteNode = useMemo(
    () =>
      (data as ExtFSRemoteNode).updatedAt !== void 0
        ? (data as ExtFSRemoteNode)
        : void 0,
    [data]
  );
  if (remoteNode !== void 0) {
    const handleRemoteClick = () => {
      const { parentItems } = extfs;
      const state = {
        ...ExtFSRemoteState,
        nodeId: remoteNode.nodeId,
      };
      setExtFS({
        ...state,
        parentItems: [...parentItems, { name: remoteNode.name, state }],
      });
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
          to={`${ExtFSNodeTagRoutePath}/${remoteNode.nodeId}`}
          disabled={!remoteNode.available}
          quantity={remoteNode.tagQuantity}
          pendingQuantity={remoteNode.pendingTagQuantity}
        />
      </ExtFSItem>
    );
  }

  const handleLocalClick = () => {
    setExtFS({
      ...ExtFSNodeState,
      parentItems: [{ name: localNode.name, state: ExtFSNodeState }],
    });
  };
  const localNode = data as ExtFSNode;
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
