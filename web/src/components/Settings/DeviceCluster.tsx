import { SettingsName } from "./meta";
import DeviceClusterIcon from "@/components/DeviceCluster/Icon";
import { TOCItem, TOCChapter, FieldsClassName } from "./MainBase";
import { withExtraState, ExtraType } from "./ExtraBase";
import { useTranslation, I18nVariant } from "@/components/App/I18Next";
import { useAppDispatch } from "@/components/App/Context";
import {
  List,
  ListItem,
  ListItemContent,
  ListItemSmall,
  ListItemText,
  ListItemButton,
  ListItemMore,
} from "@/components/App/List";

import { type ComponentProps, useRef } from "react";

export const DeviceClusterTOCItem = () => {
  const { t } = useTranslation(SettingsName);
  return (
    <TOCItem
      chapter={TOCChapter.DeviceCluster}
      text={t(`${I18nVariant.Main}.${TOCChapter.DeviceCluster}.title`)}
      Avatar={DeviceClusterIcon}
    />
  );
};

export const DeviceClusterFields = ({
  className,
}: Pick<ComponentProps<typeof List>, "className">) => {
  const { t } = useTranslation(SettingsName);
  // const { type, open } = useAppExtra<ExtraState>({});
  const clusterButtonRef = useRef<HTMLButtonElement>(null);
  const passportButtonRef = useRef<HTMLButtonElement>(null);

  const dispatch = useAppDispatch();
  const handleClusterClick = () => {
    dispatch?.(withExtraState({ type: ExtraType.ClusterSelect }));
  };

  return (
    <List className={className}>
      <ListItem className={FieldsClassName}>
        <ListItemButton onClick={handleClusterClick} ref={clusterButtonRef}>
          <ListItemContent>
            <ListItemText>
              {t(`${I18nVariant.Main}.${TOCChapter.DeviceCluster}.cluster`)}
            </ListItemText>
            <ListItemSmall>32</ListItemSmall>
            <ListItemMore />
          </ListItemContent>
        </ListItemButton>
      </ListItem>
      <ListItem className={FieldsClassName}>
        <ListItemButton disabled={true} ref={passportButtonRef}>
          <ListItemContent>
            <ListItemText>
              {t(`${I18nVariant.Main}.${TOCChapter.DeviceCluster}.passport`)}
            </ListItemText>
            <ListItemSmall>-</ListItemSmall>
            <ListItemMore />
          </ListItemContent>
        </ListItemButton>
      </ListItem>
    </List>
  );
};
