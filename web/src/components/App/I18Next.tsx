/**
 * i18next Provider for React Admin
 *
 * @see https://react.i18next.com/
 */
import i18next from "i18next";
export type { InitOptions } from "i18next";
import LanguageDetector from "i18next-browser-languagedetector";
import resourcesToBackend from "i18next-resources-to-backend";
import {
  I18nextProvider,
  useTranslation,
  initReactI18next,
} from "react-i18next";
import type { I18nextProviderProps } from "react-i18next";
import type { PropsWithChildren } from "react";

export { useTranslation };

export const DefaultNS = "app";
export const SupportLanguages = [
  { locale: "zh-CN", name: "中文" },
  { locale: "en", name: "English" },
];

const importLanguage = async (language: string, namespace: string) => {
  const { default: translation } = await import(
    `../../locales/${language}/${namespace}.json`
  );

  return translation;
};

const SupporLngs = SupportLanguages.map((l) => l.locale);
const PreferredLng = SupporLngs[0];

i18next
  .use(LanguageDetector)
  .use(initReactI18next)
  .use(resourcesToBackend(importLanguage))
  .init({
    ns: [DefaultNS],
    defaultNS: DefaultNS,
    fallbackNS: [DefaultNS],
    lng: PreferredLng,
    fallbackLng: PreferredLng,
    supportedLngs: SupporLngs,
  });

export type ProviderProps = PropsWithChildren<
  Pick<I18nextProviderProps, "defaultNS">
>;
export const Provider = ({
  children,
  defaultNS = DefaultNS,
}: ProviderProps) => {
  return (
    <I18nextProvider i18n={i18next} defaultNS={defaultNS}>
      {children}
    </I18nextProvider>
  );
};
