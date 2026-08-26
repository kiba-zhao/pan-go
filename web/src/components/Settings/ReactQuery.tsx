import {
  useQuery,
  useMutation,
  useQueryClient,
  type UseQueryOpts,
  type UseMutationOpts,
} from "@/components/App/ReactQuery";
import {
  fetchDeviceInfo,
  saveDeviceInfo,
  fetchDeviceNetwork,
  saveDeviceNetwork,
  fetchWebHost,
  saveWebHost,
} from "./api";
import { SettingsName } from "./meta";

import type { DeviceInfo, DeviceNetwork, WebHost } from "@pango/data";
export type { DeviceInfo, DeviceNetwork, WebHost };

const DeviceInfoQueryKey = `${SettingsName}/device-info`;
const DeviceNetworkQueryKey = `${SettingsName}/device-network`;
const WebHostQueryKey = `${SettingsName}/web-host`;

export const useDeviceInfo = <TData extends any = DeviceInfo>(
  opts?: UseQueryOpts<DeviceInfo, TData>,
) => {
  return useQuery({
    queryKey: [DeviceInfoQueryKey],
    queryFn: fetchDeviceInfo,
    ...(opts || {}),
  });
};

export const useDeviceNetwork = <TData extends any = DeviceNetwork>(
  opts?: UseQueryOpts<DeviceNetwork, TData>,
) => {
  return useQuery({
    queryKey: [DeviceNetworkQueryKey],
    queryFn: fetchDeviceNetwork,
    ...(opts || {}),
  });
};

export const useWebHost = <TData extends any = WebHost>(
  opts?: UseQueryOpts<WebHost, TData>,
) => {
  return useQuery({
    queryKey: [WebHostQueryKey],
    queryFn: fetchWebHost,
    ...(opts || {}),
  });
};

const useDeviceInfoMutation = (opts?: UseMutationOpts<DeviceInfo>) => {
  const queryClient = useQueryClient();
  const { onSuccess, ...options } = opts || {};
  return useMutation({
    mutationFn: saveDeviceInfo,
    onSuccess: (deviceInfo, ...args) => {
      queryClient.setQueryData([DeviceInfoQueryKey], deviceInfo);
      onSuccess?.(deviceInfo, ...args);
    },
    ...options,
  });
};

const useDeviceNetworkMutation = (opts?: UseMutationOpts<DeviceNetwork>) => {
  const queryClient = useQueryClient();
  const { onSuccess, ...options } = opts || {};
  return useMutation({
    mutationFn: saveDeviceNetwork,
    onSuccess: (deviceNetwork, ...args) => {
      queryClient.setQueryData([DeviceNetworkQueryKey], deviceNetwork);
      onSuccess?.(deviceNetwork, ...args);
    },
    ...options,
  });
};

const useWebHostMutation = (opts?: UseMutationOpts<WebHost>) => {
  const queryClient = useQueryClient();
  const { onSuccess, ...options } = opts || {};
  return useMutation({
    mutationFn: saveWebHost,
    onSuccess: (webHost, ...args) => {
      queryClient.setQueryData([WebHostQueryKey], webHost);
      onSuccess?.(webHost, ...args);
    },
    ...options,
  });
};

type UseFieldsOpts<TQueryFnData, TData = TQueryFnData> = Omit<
  UseQueryOpts<TQueryFnData, TData>,
  "select"
>;
type NameField = Pick<DeviceInfo, "name">;
export const useDeviceNameField = (
  opts?: UseFieldsOpts<DeviceInfo, NameField>,
) => {
  return useDeviceInfo({
    ...(opts || {}),
    select: (deviceInfo) => ({
      name: deviceInfo?.name,
    }),
  });
};

type MemoField = Pick<DeviceInfo, "memo">;
export const useDeviceMemoField = (
  opts?: UseFieldsOpts<DeviceInfo, MemoField>,
) => {
  return useDeviceInfo({
    ...(opts || {}),
    select: (deviceInfo) => ({
      memo: deviceInfo?.memo,
    }),
  });
};

type NetworkPortField = Pick<DeviceNetwork, "port">;
export const useNetworkPortField = (
  opts?: UseFieldsOpts<DeviceNetwork, NetworkPortField>,
) => {
  return useDeviceNetwork({
    ...(opts || {}),
    select: (deviceNetwork) => ({
      port: deviceNetwork?.port,
    }),
  });
};

type PublicAddrsField = Pick<DeviceNetwork, "publicAddrs">;
export const usePublicAddrsField = (
  opts?: UseFieldsOpts<DeviceNetwork, PublicAddrsField>,
) => {
  return useDeviceNetwork({
    ...(opts || {}),
    select: (deviceNetwork) => ({
      publicAddrs: deviceNetwork?.publicAddrs,
    }),
  });
};

type BroadcastAddrsField = Pick<DeviceNetwork, "broadcastAddrs">;
export const useBroadcastAddrsField = (
  opts?: UseFieldsOpts<DeviceNetwork, BroadcastAddrsField>,
) => {
  return useDeviceNetwork({
    ...(opts || {}),
    select: (deviceNetwork) => ({
      broadcastAddrs: deviceNetwork?.broadcastAddrs,
    }),
  });
};

type WebPortField = Pick<WebHost, "webPort">;
export const useWebPortField = (
  opts?: UseFieldsOpts<WebHost, WebPortField>,
) => {
  return useWebHost({
    ...(opts || {}),
    select: (webHost) => ({
      webPort: webHost?.webPort,
    }),
  });
};

export {
  useDeviceNetworkMutation as useNetworkEnableMutation,
  useDeviceNetworkMutation as useBroadcastEnableMutation,
  useWebHostMutation as useWebEnableMutation,
  useWebHostMutation as useLocalHostOnlyMutation,
  useDeviceInfoMutation as useDeviceNameMutation,
  useDeviceInfoMutation as useDeviceMemoMutation,
  useDeviceNetworkMutation as useNetworkPortMutation,
  useDeviceNetworkMutation as usePublicAddrsMutation,
  useDeviceNetworkMutation as useBroadcastAddrsMutation,
  useWebHostMutation as useWebPortMutation,
};
