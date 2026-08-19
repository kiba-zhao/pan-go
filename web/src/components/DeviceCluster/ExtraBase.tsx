import { DialogExtraState, withDialogExtraState } from "@/components/App/Extra";

export enum ExtraType {
  ClusterSwitch = "clusterSwitch",
  ClusterRemove = "clusterRemove",
  ClusterAdd = "clusterAdd",
  ClusterEdit = "clusterEdit",
  ClusterPassphraseEdit = "clusterPassphraseEdit",
  DeviceAdd = "deviceAdd",
  DevicesRemove = "devicesRemove",
  DeviceEdit = "deviceEdit",
}

type CustomExtraState = DialogExtraState & {
  clusterId?: number;
  deviceIds?: number[];
  deviceId?: number;
};

export type ExtraState = Parameters<
  typeof withDialogExtraState<ExtraType, CustomExtraState>
>[0];
export function withExtraState(extraState: ExtraState) {
  return withDialogExtraState(extraState);
}

export type ExtraProps = { extraState: ExtraState };
