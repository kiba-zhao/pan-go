import type { ReactNode } from "react";
import { createContext, useContext } from "react";

export type BrowerContext = {
  window: Window;
};

const Context = createContext<BrowerContext | null>(null);

export const useBrower = () => useContext<BrowerContext | null>(Context);

type BrowerProviderProps = {
  children: ReactNode;
} & BrowerContext;
export const BrowerProvider = ({ children, ...ctx }: BrowerProviderProps) => {
  return <Context.Provider value={ctx}>{children}</Context.Provider>;
};
