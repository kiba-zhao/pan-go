import {
  simple,
  withJSONBody,
  withMethod,
  withPath,
  withQuery,
} from "fetch-utils";
import { ModuleItem, ModuleSearchResult } from "./api.tsx";

const ROOT_PATH = `${import.meta.env.BASE_URL}${
  import.meta.env.VITE_API_PATH || ""
}`;

const { fetchOne, fetchMany } = simple(ROOT_PATH);

export async function SearchModules(
  keyword: string
): Promise<ModuleSearchResult> {
  return await fetchMany(withPath("modules", "merge"), withQuery({ keyword }));
}

export async function GetModule(name: string): Promise<ModuleItem> {
  return await fetchOne(withPath(`modules/${name}`, "merge"));
}

export async function SetModuleEnabled(
  name: string,
  enabled: boolean
): Promise<ModuleItem> {
  return await fetchOne(
    withPath(`modules/${name}`, "merge"),
    withMethod("PATCH"),
    withJSONBody({ enabled })
  );
}
