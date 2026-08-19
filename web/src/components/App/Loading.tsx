import { PageInfoSection } from "./Layout";
import { BlocksShuffle3Icon } from "./Icon";
import { useTranslation, I18nVariant } from "./I18Next";

export const AppLoading = () => {
  const { t } = useTranslation();
  return (
    <PageInfoSection
      logo={<BlocksShuffle3Icon className="w-26" />}
      title={t(`${I18nVariant.Main}.loading.title`)}
      description={t(`${I18nVariant.Main}.loading.description`)}
    />
  );
};
