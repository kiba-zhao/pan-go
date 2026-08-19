import { PageInfoSection } from "./Layout";
import { BlockOutlineIcon } from "./Icon";
import { useTranslation, I18nVariant } from "./I18Next";
import { Outlets } from "./Outlets";

export const NotFoundOutlets = () => {
  const { t } = useTranslation();
  return (
    <Outlets>
      <PageInfoSection
        logo={<BlockOutlineIcon className="w-26" />}
        title={t(`${I18nVariant.Main}.notFound.title`)}
        description={t(`${I18nVariant.Main}.notFound.description`)}
      />
    </Outlets>
  );
};
