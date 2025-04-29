/**
 * ExtFS Home Component Definition File
 */
import type { ExtFSItemRecord } from "./Item";
import { ExtFSItem, ExtFSItems, useExtFSItem, ExtFSItemSettings } from "./Item";
import { newExtFSState as newExtFSStateWithNodeItem } from "./NodeItem";
import { newExtFSState as newExtFSStateWithRemote } from "./RemoteItem";
import { useExtFS } from "./State";

import { useTranslation } from "../i18n/Context";
import type { AppNode } from "./api";
import { selectAllAppSettings, selectAllAppNodes } from "./api";
import { APP_SETTINGS_QUERY_KEY } from "../AppSettings/Page";
import { generateEditPath } from "../Route/utils";
import { AppNodeIcon, AppNodePath, AppNodeEditI18nKey } from "../AppNode/Route";

import { useQuery } from "@tanstack/react-query";
import { useMemo } from "react";

import CloudIcon from "@mui/icons-material/Cloud";

export const REMOTE_NODES_QUERY_KEY = ["extfs-remote-nodes"];
const ExtFSHomeMode = "H";

export const ExtFSHomeState = {
  mode: ExtFSHomeMode,
  queryKeyList: [APP_SETTINGS_QUERY_KEY, REMOTE_NODES_QUERY_KEY],
};

type ExtFSLocalNode = {
  name: string;
  peerId: string;
};
/**
 * @function HomeItems
 * @description
 * A component that renders ExtFS items for the home route.
 * @returns {JSX.Element} An element that renders a list of items.
 * @example
 * import { HomeItems } from "./Home";
 * <HomeItems />
 */
export const HomeItems = () => {
  const [extfs, _] = useExtFS();

  const { data: settings, isFetching: isSettingsFetching } = useQuery({
    queryKey: APP_SETTINGS_QUERY_KEY,
    queryFn: async () => await selectAllAppSettings(),
    enabled: extfs.mode === ExtFSHomeMode,
  });
  const nodeItem = useMemo(() => {
    if (!settings) return;
    return {
      name: settings.name,
      peerId: settings.peerId,
    };
  }, [settings]);

  const {
    data: remotes,
    isFetching: isRemotesFetching,
    error,
  } = useQuery({
    queryKey: REMOTE_NODES_QUERY_KEY,
    queryFn: async () => await selectAllAppNodes(),
    enabled: extfs.mode === ExtFSHomeMode,
  });

  const items = useMemo(() => {
    const items_: Array<ExtFSLocalNode | AppNode> =
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
    <ExtFSItems items={items} isFetching={isFetching} error={error || void 0}>
      <HomeItem />
    </ExtFSItems>
  );
};

/**
 * A component that renders a list item for the home route.
 * @returns {JSX.Element} An element that renders a list item.
 * @example
 * import { HomeItem } from "./Home";
 * <HomeItem />
 */
const HomeItem = () => {
  const { t } = useTranslation();
  const { style, item }: ExtFSItemRecord<ExtFSLocalNode | AppNode> =
    useExtFSItem();

  const [extfs, setExtFS] = useExtFS();

  const [remoteNode, remoteSettingsUrl] = useMemo(() => {
    const rnode = item as AppNode;
    if (rnode.id === void 0) return [];
    return [rnode, generateEditPath(AppNodePath, rnode.id)];
  }, [item]);

  if (remoteNode !== void 0 && remoteSettingsUrl !== void 0) {
    const handleRemoteClick = () => {
      const state = newExtFSStateWithRemote(extfs, remoteNode);
      setExtFS(state);
    };

    return (
      <ExtFSItem
        style={style}
        primary={remoteNode.name}
        secondary={remoteNode?.updatedAt.toString()}
        avatarIcon={
          <AppNodeIcon
            fontSize="large"
            color={remoteNode.online ? "primary" : "disabled"}
          />
        }
        extIcon={<CloudIcon fontSize="small" />}
        disabled={!remoteNode?.online}
        onClick={handleRemoteClick}
      >
        <ExtFSItemSettings
          to={remoteSettingsUrl}
          title={t(AppNodeEditI18nKey)}
        />
      </ExtFSItem>
    );
  }

  const localNode = item as ExtFSLocalNode;

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
