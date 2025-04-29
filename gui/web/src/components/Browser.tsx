import type { ReactNode } from "react";
import { createContext, useContext } from "react";

export type BrowserContext = {
  window: Window;
};

const Context = createContext<BrowserContext | null>(null);

export const useBrowser = () => useContext<BrowserContext | null>(Context);

type BrowserProviderProps = {
  children: ReactNode;
} & BrowserContext;
export const BrowserProvider = ({ children, ...ctx }: BrowserProviderProps) => {
  return <Context.Provider value={ctx}>{children}</Context.Provider>;
};
