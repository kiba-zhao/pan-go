import { createContext, useContext } from "react";
import * as api from "./api";

const Context = createContext<APIContext>(null!);

export function APIProvider({ children }: { children: React.ReactNode }) {
  return <Context.Provider value={api}>{children}</Context.Provider>;
}

export function useAPI(): APIContext {
  return useContext<APIContext>(Context);
}

export interface APIContext {
  SearchModules(keyword: string): Promise<ModuleSearchResult>;
  GetModule(name: string): Promise<ModuleItem>;
  SetModuleEnabled(name: string, enabled: boolean): Promise<ModuleItem>;
}

export interface ModuleItem {
  Avatar: string;
  Name: string;
  Desc: string;
  Enabled: boolean;
  ReadOnly: boolean;
  HasWeb: boolean;
}

export interface SearchResult<T> {
  Total: number;
  Items: T[];
}

export type ModuleSearchResult = SearchResult<ModuleItem>;
