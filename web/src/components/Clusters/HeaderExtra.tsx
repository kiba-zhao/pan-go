import { Link } from "@/components/App/Route";
import { useTranslation } from "@/components/App/I18Next";
import { ClustersName, ClustersRoutePath } from "./meta";

export const GoBack = () => {
  const { t } = useTranslation(ClustersName);

  return (
    <Link className="text-sm" to={ClustersRoutePath}>
      {t(`header-extra.back`)}
    </Link>
  );
};
