import {
  DialogExtraState,
  withDialogExtraState,
} from "@/components/App/Dialog";
import { AppExtraState } from "@/components/App/Extra";

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

type ExtraState = {
  clusterId?: number;
  deviceIds?: number[];
  deviceId?: number;
} & AppExtraState<ExtraType> &
  DialogExtraState;

export function withExtraState(extraState: ExtraState) {
  return withDialogExtraState(extraState);
}

export type ExtraProps = { extraState: ExtraState };
