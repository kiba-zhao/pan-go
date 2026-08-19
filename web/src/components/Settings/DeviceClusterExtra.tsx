import { SettingsName, SettingsRoutePath } from "./meta";
import { withExtraState, type ExtraProps } from "./ExtraBase";

import { useAppDispatch } from "@/components/App/Context";
import { useTranslation, I18nVariant } from "@/components/App/I18Next";

import { ClusterSwitchExtraBase } from "@/components/DeviceCluster/ClusterExtra";
import { withGoBackState } from "@/components/DeviceCluster/HeaderExtra";

export const ClusterSelectExtra = ({ extraState }: ExtraProps) => {
  const { type, open } = extraState;

  const dispatch = useAppDispatch();
  const handleClose = () => {
    dispatch?.(withExtraState({ type, open: false }));
  };

  const { t } = useTranslation(SettingsName);
  const GoBackState = withGoBackState({
    to: SettingsRoutePath,
    text: t(`${I18nVariant.Action}.backToSettings`),
  });

  return (
    <ClusterSwitchExtraBase
      open={open}
      onClose={handleClose}
      locationState={GoBackState}
    />
  );
};
