import i18next from "i18next";
import LanguageDetector from "i18next-browser-languagedetector";
import resourcesToBackend from "i18next-resources-to-backend";

import {
  convertRaTranslationsToI18next,
  useI18nextProvider,
} from "ra-i18n-i18next";

import zh from "@haxqer/ra-language-chinese";
import en from "ra-language-english";

const translations: Record<string, object> = {
  en: convertRaTranslationsToI18next(en),
  "zh-CN": convertRaTranslationsToI18next(zh),
};

const importLanguage = async (language: string, namespace: string) => {
  const { default: translation } = await import(
    `./locales/${language}/${namespace}.json`
  );
  const defaults = translations[language] || translations.en;
  return { ...translation, ...defaults };
};

i18next.use(LanguageDetector).use(resourcesToBackend(importLanguage));

export const useI18nProvider = () =>
  useI18nextProvider({
    i18nextInstance: i18next,
    availableLocales: [
      { locale: "en", name: "English" },
      { locale: "zh-CN", name: "中文" },
    ],
  });
