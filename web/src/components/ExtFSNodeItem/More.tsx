import { MoreLinkItem } from "../Common/MoreButton";
import { useTranslation } from "../I18Next/Context";
import {
  ExtFSNodeItemCreatePath,
  ExtFSNodeItemCreateI18nKey,
  ExtFSNodeItemIcon,
} from "./Route";

export const NewExtFSNodeItemMore = ({ selected }: { selected?: boolean }) => {
  const { t } = useTranslation();
  return (
    <MoreLinkItem
      to={ExtFSNodeItemCreatePath}
      selected={selected}
      icon={<ExtFSNodeItemIcon />}
    >
      {t(ExtFSNodeItemCreateI18nKey)}
    </MoreLinkItem>
  );
};
