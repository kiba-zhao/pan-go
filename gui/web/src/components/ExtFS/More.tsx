import { useTranslation } from "../I18Next/Context";
import { MoreLinkItem } from "../Common/MoreButton";
import { ExtFSPath, ExtFSI18nKey, ExtFSIcon } from "./Route";

export const ExtFSMore = ({ selected }: { selected?: boolean }) => {
  const { t } = useTranslation();
  return (
    <MoreLinkItem to={ExtFSPath} icon={<ExtFSIcon />} selected={selected}>
      {t(ExtFSI18nKey)}
    </MoreLinkItem>
  );
};
