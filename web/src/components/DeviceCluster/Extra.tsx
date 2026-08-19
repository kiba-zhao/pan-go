import { ExtraType } from "./ExtraBase";
import {
  ClusterSwitchExtra,
  ClusterRemoveExtra,
  ClusterEditExtra,
  ClusterPassphraseEditExtra,
} from "./ClusterExtra";
import { DevicesRemoveExtra, DeviceEditExtra } from "./DeviceExtra";

import { Extra as AppExtra } from "@/components/App/Extra";

export const ClusterExtra = () => {
  return (
    <>
      <AppExtra as={ClusterSwitchExtra} extraType={ExtraType.ClusterSwitch} />
      <AppExtra as={ClusterRemoveExtra} extraType={ExtraType.ClusterRemove} />
      <AppExtra as={ClusterEditExtra} extraType={ExtraType.ClusterEdit} />
      <AppExtra
        as={ClusterPassphraseEditExtra}
        extraType={ExtraType.ClusterPassphraseEdit}
      />
      <AppExtra as={DevicesRemoveExtra} extraType={ExtraType.DevicesRemove} />
      <AppExtra as={DeviceEditExtra} extraType={ExtraType.DeviceEdit} />
    </>
  );
};
