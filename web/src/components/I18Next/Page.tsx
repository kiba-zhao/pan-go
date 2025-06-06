import { useTranslation } from "./Context";
import type { InitOptions } from "./Context";
import { isArray } from "lodash";
import { APP_NAMESPACE } from "./constants";

type PageI18NextProps = { defaultNS?: string } & Pick<InitOptions, "ns">;
export const PageI18Next = ({ ns, defaultNS }: PageI18NextProps) => {
  const { i18n } = useTranslation();
  if (!i18n) return null;

  if (ns !== void 0) {
    const namespaces: string[] = isArray(ns) ? ns : [ns];
    namespaces.forEach((_) => {
      if (!i18n.hasLoadedNamespace(_)) i18n.loadNamespaces(_);
    });
  }

  i18n.setDefaultNamespace(defaultNS || APP_NAMESPACE);

  return null;
};
