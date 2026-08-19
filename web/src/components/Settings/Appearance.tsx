import { SettingsName } from "./meta";
import { Palette } from "./Icon";
import { TOCItem, TOCChapter, FieldsClassName } from "./MainBase";

import { type ComponentProps } from "react";

import { useTranslation, I18nVariant } from "@/components/App/I18Next";
import {
  List,
  ListItem,
  ListItemSmall,
  ListItemText,
} from "@/components/App/List";

export const AppearanceTOCItem = () => {
  const { t } = useTranslation(SettingsName);
  return (
    <TOCItem
      chapter={TOCChapter.Appearance}
      text={t(`${I18nVariant.Main}.${TOCChapter.Appearance}.title`)}
      Avatar={Palette}
    />
  );
};

export const AppearanceFields = ({
  className,
}: Pick<ComponentProps<typeof List>, "className">) => {
  const { t } = useTranslation(SettingsName);
  return (
    <List className={className}>
      <ListItem className={FieldsClassName}>
        <ListItemText>
          {t(`${I18nVariant.Main}.${TOCChapter.Appearance}.theme`)}
        </ListItemText>
        <ListItemSmall>
          {t(`${I18nVariant.Main}.${TOCChapter.Appearance}.system`)}
        </ListItemSmall>
      </ListItem>
    </List>
  );
};
