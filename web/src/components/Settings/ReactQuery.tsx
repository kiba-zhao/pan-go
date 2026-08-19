import {
  useQuery,
  useMutation,
  useQueryClient,
  type UseQueryOpts,
  type UseMutationOpts,
} from "@/components/App/ReactQuery";
import {
  fetchSettings,
  fetchHostSettings,
  saveSettings,
  saveHostSettings,
} from "./api";
import { SettingsName } from "./meta";

import type { Settings, HostSettings } from "@pango/data";
export type { Settings, HostSettings };

const SettingsQueryKey = `${SettingsName}`;
const HostSettingsQueryKey = `${SettingsName}/host-settings`;

const useSettings = <TData extends any>(
  opts?: UseQueryOpts<Settings, TData>,
) => {
  return useQuery({
    queryKey: [SettingsQueryKey],
    queryFn: fetchSettings,
    ...(opts || {}),
  });
};

const useHostSettings = <TData extends any>(
  opts?: UseQueryOpts<HostSettings, TData>,
) => {
  return useQuery({
    queryKey: [HostSettingsQueryKey],
    queryFn: fetchHostSettings,
    ...(opts || {}),
  });
};

const useSettingsMutation = (opts?: UseMutationOpts<Settings>) => {
  const queryClient = useQueryClient();
  const { onSuccess, ...options } = opts || {};
  return useMutation({
    mutationFn: saveSettings,
    onSuccess: (settings, ...args) => {
      queryClient.setQueryData([SettingsQueryKey], settings);
      onSuccess?.(settings, ...args);
    },
    ...options,
  });
};

const useHostSettingsMutation = (opts?: UseMutationOpts<HostSettings>) => {
  const queryClient = useQueryClient();
  const { onSuccess, ...options } = opts || {};
  return useMutation({
    mutationFn: saveHostSettings,
    onSuccess: (settings, ...args) => {
      queryClient.setQueryData([HostSettingsQueryKey], settings);
      onSuccess?.(settings, ...args);
    },
    ...options,
  });
};

type UseFieldsOpts<TQueryFnData, TData = TQueryFnData> = Omit<
  UseQueryOpts<TQueryFnData, TData>,
  "select"
>;
type InfoFields = Pick<Settings, "name" | "memo">;
export const useDeviceInfoFields = (
  opts?: UseFieldsOpts<Settings, InfoFields>,
) => {
  return useSettings({
    ...(opts || {}),
    select: (settings) => ({
      name: settings?.name,
      memo: settings?.memo,
    }),
  });
};

type NetworkFields = Pick<
  Settings,
  "enabled" | "peerPort" | "broadcastEnabled" | "broadcastAddrs" | "publicAddrs"
>;
export const useDeviceNetworkFields = (
  opts?: UseFieldsOpts<Settings, NetworkFields>,
) => {
  return useSettings({
    ...(opts || {}),
    select: (settings) => ({
      enabled: settings?.enabled,
      peerPort: settings?.peerPort,
      broadcastEnabled: settings?.broadcastEnabled,
      broadcastAddrs: settings?.broadcastAddrs,
      publicAddrs: settings?.publicAddrs,
    }),
  });
};

type WebFields = Pick<HostSettings, "webEnabled" | "localHostOnly" | "webPort">;
export const useDeviceWebFields = (
  opts?: UseFieldsOpts<HostSettings, WebFields>,
) => {
  return useHostSettings({
    ...(opts || {}),
    select: (settings) => ({
      webEnabled: settings?.webEnabled,
      localHostOnly: settings?.localHostOnly,
      webPort: settings?.webPort,
    }),
  });
};

type NameField = Pick<Settings, "name">;
export const useDeviceNameField = (
  opts?: UseFieldsOpts<Settings, NameField>,
) => {
  return useSettings({
    ...(opts || {}),
    select: (settings) => ({
      name: settings?.name,
    }),
  });
};

type MemoField = Pick<Settings, "memo">;
export const useDeviceMemoField = (
  opts?: UseFieldsOpts<Settings, MemoField>,
) => {
  return useSettings({
    ...(opts || {}),
    select: (settings) => ({
      memo: settings?.memo,
    }),
  });
};

type PeerPortField = Pick<Settings, "peerPort">;
export const usePeerPortField = (
  opts?: UseFieldsOpts<Settings, PeerPortField>,
) => {
  return useSettings({
    ...(opts || {}),
    select: (settings) => ({
      peerPort: settings?.peerPort,
    }),
  });
};

type BroadcastAddrsField = Pick<Settings, "broadcastAddrs">;
export const useBroadcastAddrsField = (
  opts?: UseFieldsOpts<Settings, BroadcastAddrsField>,
) => {
  return useSettings({
    ...(opts || {}),
    select: (settings) => ({
      broadcastAddrs: settings?.broadcastAddrs,
    }),
  });
};

type PublicAddrsField = Pick<Settings, "publicAddrs">;

export const usePublicAddrsField = (
  opts?: UseFieldsOpts<Settings, PublicAddrsField>,
) => {
  return useSettings({
    ...(opts || {}),
    select: (settings) => ({
      publicAddrs: settings?.publicAddrs,
    }),
  });
};

type WebPortField = Pick<HostSettings, "webPort">;
export const useWebPortField = (
  opts?: UseFieldsOpts<HostSettings, WebPortField>,
) => {
  return useHostSettings({
    ...(opts || {}),
    select: (settings) => ({
      webPort: settings?.webPort,
    }),
  });
};

export {
  useSettingsMutation as useNetworkEnableMutation,
  useSettingsMutation as useBroadcastEnableMutation,
  useHostSettingsMutation as useWebEnableMutation,
  useHostSettingsMutation as useLocalHostOnlyMutation,
  useSettingsMutation as useDeviceNameMutation,
  useSettingsMutation as useDeviceMemoMutation,
  useSettingsMutation as usePeerPortMutation,
  useSettingsMutation as useBroadcastAddrsMutation,
  useSettingsMutation as usePublicAddrsMutation,
  useHostSettingsMutation as useWebPortMutation,
};
