import { MoreLinkItem } from "../Common/MoreButton";
import { useTranslation } from "../I18Next/Context";
import { AppSettingsPath, AppSettingsI18nKey, AppSettingsIcon } from "./Route";

export const AppSettingsMore = ({ selected }: { selected?: boolean }) => {
  const { t } = useTranslation();
  return (
    <MoreLinkItem
      to={AppSettingsPath}
      icon={<AppSettingsIcon />}
      selected={selected}
    >
      {t(AppSettingsI18nKey)}
    </MoreLinkItem>
  );
};
