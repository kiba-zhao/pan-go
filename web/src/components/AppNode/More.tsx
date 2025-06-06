import {
  AppNodeCreatePath,
  AppNodeCreateI18nKey,
  AppNodeNS,
  AppNodeIcon,
} from "./Route";
import { MoreLinkItem } from "../Common/MoreButton";
import { useTranslation } from "../I18Next/Context";

export const NewAppNodeMore = ({ selected }: { selected?: boolean }) => {
  const { t } = useTranslation();
  return (
    <MoreLinkItem
      to={AppNodeCreatePath}
      icon={<AppNodeIcon />}
      selected={selected}
    >
      {t(AppNodeCreateI18nKey, { ns: AppNodeNS })}
    </MoreLinkItem>
  );
};
