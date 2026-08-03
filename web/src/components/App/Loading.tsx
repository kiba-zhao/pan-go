import { PageInfoSection } from "./Layout";
import { BlocksShuffle3Icon } from "./Icon";
import { useTranslation } from "./I18Next";

export const AppLoading = () => {
  const { t } = useTranslation();
  return (
    <PageInfoSection
      logo={<BlocksShuffle3Icon className="w-26" />}
      title={t("pages.loading.title")}
      description={t("pages.loading.description")}
    />
  );
};
