import { ExtraType } from "./ExtraBase";
import { DeviceNameEditExtra, DeviceMemoEditExtra } from "./DeviceInfoExtra";
import {
  NetworkPortEditExtra,
  BroadcastAddrsEditExtra,
  PublicAddrsEditExtra,
} from "./DeviceNetworkExtra";
import { WebPortEditExtra } from "./DeviceWebExtra";
import { ClusterSelectExtra } from "./DeviceClusterExtra";

import { AppExtra } from "@/components/App/Extra";

const SettingsExtra = () => (
  <>
    <AppExtra as={DeviceNameEditExtra} extraType={ExtraType.DeviceNameEdit} />
    <AppExtra as={DeviceMemoEditExtra} extraType={ExtraType.DeviceMemoEdit} />
    <AppExtra as={NetworkPortEditExtra} extraType={ExtraType.PeerPortEdit} />
    <AppExtra
      as={BroadcastAddrsEditExtra}
      extraType={ExtraType.BroadcastAddrsEdit}
    />
    <AppExtra as={PublicAddrsEditExtra} extraType={ExtraType.PublicAddrsEdit} />
    <AppExtra as={WebPortEditExtra} extraType={ExtraType.WebPortEdit} />
    <AppExtra as={ClusterSelectExtra} extraType={ExtraType.ClusterSelect} />
  </>
);

export default SettingsExtra;
