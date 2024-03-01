import { createContext, useContext } from "react";

export interface APIContext {
  SearchModules(keyword: string): Promise<ModuleSearchResult>;
}

export const Context = createContext<APIContext>(null!);

export function useAPI(): APIContext {
  return useContext<APIContext>(Context);
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
