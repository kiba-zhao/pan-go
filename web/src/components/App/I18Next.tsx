/**
 * i18next Provider for React Admin
 *
 * @see https://react.i18next.com/
 */
import i18next from "i18next";
export type { InitOptions, TFunction } from "i18next";
import LanguageDetector from "i18next-browser-languagedetector";
import resourcesToBackend from "i18next-resources-to-backend";
import {
  I18nextProvider,
  useTranslation,
  initReactI18next,
} from "react-i18next";
import type { I18nextProviderProps } from "react-i18next";
import {
  createContext,
  useReducer,
  useContext,
  useMemo,
  type Dispatch,
} from "react";

export { useTranslation };

export const DefaultNS = "app";
export const SupportLanguages = [
  { locale: "zh-CN", name: "中文" },
  { locale: "en", name: "English" },
];

const NamespaceSeparator = "/";
const importLanguage = async (language: string, namespace: string) => {
  const { default: translation } = await await import(
    `../../locales/${language}/${namespace}.json`
  );

  if (translation && translation.external) {
    for (const key of Object.keys(translation.external)) {
      i18next.addResourceBundle(
        language,
        `${namespace}${NamespaceSeparator}${key}`,
        translation.external[key],
      );
    }
  }
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

type AppI18nState = {
  namespace: string;
};
type AppI18nAction = AppI18nState;

function AppI18nReducer(state: AppI18nState, action: AppI18nAction) {
  return {
    ...state,
    ...action,
  };
}

const appI18nContext = createContext<AppI18nState>({ namespace: DefaultNS });
const appI18nDispatchContext = createContext<Dispatch<AppI18nAction> | null>(
  null,
);

export const useAppI18n = () => useContext(appI18nContext);
export const useAppI18nDispatch = () => useContext(appI18nDispatchContext);

export function withAppI18nResetAction() {
  return { namespace: DefaultNS };
}
export function withAppI18nAction(payload: AppI18nState): AppI18nAction {
  return payload;
}

export type ProviderProps = Pick<
  I18nextProviderProps,
  "defaultNS" | "children"
>;
export const Provider = ({
  children,
  defaultNS = DefaultNS,
}: ProviderProps) => {
  const namespace = typeof defaultNS === "string" ? defaultNS : defaultNS[0];
  const [state, dispatch] = useReducer(AppI18nReducer, { namespace });
  return (
    <I18nextProvider i18n={i18next} defaultNS={defaultNS}>
      <appI18nContext.Provider value={state}>
        <appI18nDispatchContext.Provider value={dispatch}>
          {children}
        </appI18nDispatchContext.Provider>
      </appI18nContext.Provider>
    </I18nextProvider>
  );
};

export enum I18nVariant {
  Main = "main",
  Extra = "extra",
  HeaderExtra = "header-extra",
  Badge = "badge",
  Action = "action",
  Error = "error",
  Navigation = "navigation",
  Form = "form",
}

export const UnknownI18nKey = "Unknown";

type I18nNamespace = string;
export const useExternalNamespace = (
  namespace: I18nNamespace = DefaultNS,
  fallback = false,
) => {
  const { namespace: appNamespace } = useAppI18n();
  return useMemo<I18nNamespace | [I18nNamespace, I18nNamespace]>(() => {
    if (namespace === appNamespace) return namespace;
    const inlineNamespace = `${appNamespace}${NamespaceSeparator}${namespace}`;
    if (!fallback) return inlineNamespace;
    return [inlineNamespace, namespace];
  }, [namespace, appNamespace]);
};

type TFunction = NonNullable<ReturnType<typeof useTranslation>["t"]>;
type TFunctionOptions = NonNullable<Exclude<Parameters<TFunction>[1], string>>;
export function withAppError(
  error?: Error,
  opts?: TFunctionOptions,
): [string, TFunctionOptions] {
  const opts_ = opts || ({} as TFunctionOptions);
  if (!opts_.defaultValue) opts_.defaultValue = error?.message || "";
  return [`${I18nVariant.Error}.${error?.name || UnknownI18nKey}`, opts_];
}
