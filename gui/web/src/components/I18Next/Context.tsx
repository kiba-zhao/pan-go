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

import type { ReactNode } from "react";

import { LANGUAGES, APP_NAMESPACE } from "./constants";

export { useTranslation };

const importLanguage = async (language: string, namespace: string) => {
  const { default: translation } = await import(
    `../../locales/${language}/${namespace}.json`
  );

  return translation;
};

i18next
  .use(LanguageDetector)
  .use(initReactI18next)
  .use(resourcesToBackend(importLanguage))
  .init({
    ns: [APP_NAMESPACE],
    defaultNS: APP_NAMESPACE,
    fallbackNS: [APP_NAMESPACE],
    lng: LANGUAGES[0].locale,
    fallbackLng: LANGUAGES[0].locale,
    supportedLngs: LANGUAGES.map((l) => l.locale),
  });

export const Provider = ({ children }: { children: ReactNode }) => {
  return <I18nextProvider i18n={i18next}>{children}</I18nextProvider>;
};
