import type { ReactNode } from "react";
import { createContext, useContext } from "react";

export type BrowserState = {
  window: Window;
};

const Context = createContext<BrowserState | null>(null);

export const useBrowser = () => useContext<BrowserState | null>(Context);

type BrowserProviderProps = {
  children: ReactNode;
} & BrowserState;
export const BrowserProvider = ({ children, ...ctx }: BrowserProviderProps) => {
  return <Context.Provider value={ctx}>{children}</Context.Provider>;
};
