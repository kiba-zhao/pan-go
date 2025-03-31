import type { ReactNode } from "react";
import { createContext, useCallback, useContext } from "react";

import type { TranslateFunction as RATranslateFunction } from "react-admin";
import { useTranslate as useRATranslate } from "react-admin";

export type TranslateFunction = RATranslateFunction;

const Context = createContext<TranslateFunction>(() => {
  throw new Error("Lack of Translate Provider");
});

export const useTranslate = () => useContext<TranslateFunction>(Context);

export const DEFAULT_BASE = "custom";

type TranslateProviderProps = {
  children: ReactNode;
  base?: string;
};
export const TranslateProvider = ({
  children,
  base = DEFAULT_BASE,
}: TranslateProviderProps) => {
  const raTranslate = useRATranslate();
  const translate = useCallback(
    (...args: Parameters<TranslateFunction>) =>
      raTranslate(base ? `${base}.${args[0]}` : args[0], ...args.slice(1)),
    [base]
  );
  return <Context.Provider value={translate}>{children}</Context.Provider>;
};
