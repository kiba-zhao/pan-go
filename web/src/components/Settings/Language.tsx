import { SettingsName } from "./meta";
import { Languages } from "./Icon";
import { TOCItem, TOCChapter, FieldsClassName } from "./MainBase";

import { type ComponentProps } from "react";

import { useTranslation, I18nVariant } from "@/components/App/I18Next";
import {
  List,
  ListItem,
  ListItemSmall,
  ListItemText,
} from "@/components/App/List";

export const LanguageTOCItem = () => {
  const { t } = useTranslation(SettingsName);
  return (
    <TOCItem
      className="rounded-b-sm"
      chapter={TOCChapter.Language}
      text={t(`${I18nVariant.Main}.${TOCChapter.Language}.title`)}
      Avatar={Languages}
    />
  );
};

export const LanguageFields = ({
  className,
}: Pick<ComponentProps<typeof List>, "className">) => {
  const { t } = useTranslation(SettingsName);
  return (
    <List className={className}>
      <ListItem className={FieldsClassName}>
        <ListItemText>
          {t(`${I18nVariant.Main}.${TOCChapter.Language}.lang`)}
        </ListItemText>
        <ListItemSmall>
          {t(`${I18nVariant.Main}.${TOCChapter.Language}.system`)}
        </ListItemSmall>
      </ListItem>
    </List>
  );
};
