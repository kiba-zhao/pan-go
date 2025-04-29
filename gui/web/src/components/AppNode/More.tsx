import { AppNodeCreatePath, AppNodeCreateI18nKey, AppNodeIcon } from "./Route";
import { MoreLinkItem } from "../Common/MoreButton";
import { useTranslation } from "../i18n/Context";

export const NewAppNodeMore = ({ selected }: { selected?: boolean }) => {
  const { t } = useTranslation();
  return (
    <MoreLinkItem
      to={AppNodeCreatePath}
      icon={<AppNodeIcon />}
      selected={selected}
    >
      {t(AppNodeCreateI18nKey)}
    </MoreLinkItem>
  );
};
