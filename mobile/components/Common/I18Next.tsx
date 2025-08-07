/**
 * i18next Provider for React Admin
 *
 * @see https://react.i18next.com/
 */
import i18next from 'i18next';
import resourcesToBackend from 'i18next-resources-to-backend';
import {I18nextProvider, initReactI18next, useTranslation} from 'react-i18next';
import * as RNLocalize from 'react-native-localize';

export {useTranslation};

import {type ReactNode} from 'react';

export const DefaultNS = 'app';
export const SupportLanguages = [
  {locale: 'zh-CN', name: '中文'},
  {locale: 'en', name: 'English'},
];

const languageSources = new Map<string, () => Promise<Record<string, any>>>();
addLanguageSource(SupportLanguages[0].locale, DefaultNS, () =>
  import(`../../locales/zh-CN/app.json`).then(m => m.default),
);
addLanguageSource(SupportLanguages[1].locale, DefaultNS, () =>
  import(`../../locales/en/app.json`).then(m => m.default),
);

const importLanguage = async (language: string, namespace: string) => {
  const key = `${namespace}.${language}`;
  return languageSources.has(key) ? languageSources.get(key)!() : null;
};

const SupporLngs = SupportLanguages.map(l => l.locale);
const PreferredLng =
  RNLocalize.getLocales().find(_ =>
    SupporLngs.includes(_.languageTag) ? _ : void 0,
  )?.languageTag || SupporLngs[0];

i18next
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

export const Provider = ({children}: {children: ReactNode}) => {
  return <I18nextProvider i18n={i18next}>{children}</I18nextProvider>;
};

type LanguageSourceLoader = () => Promise<Record<string, any>>;

export function addLanguageSource(
  lang: string,
  ns: string,
  loader: LanguageSourceLoader,
  force = false,
) {
  if (!force && languageSources.has(`${ns}.${lang}`)) return;
  languageSources.set(`${ns}.${lang}`, loader);
}
