import { ModuleSearchResult } from "./api.tsx";

const ROOT_PATH = `${import.meta.env.BASE_URL}${
  import.meta.env.VITE_API_PATH || ""
}`;

async function createError(res: Response): Promise<any> {
  let contentType = res.headers.get("Content-Type");
  if (contentType?.trim().startsWith("application/json")) {
    return await res.json();
  }
  let message = await res.text();
  return new Error(message);
}

export async function SearchModules(
  keyword: string
): Promise<ModuleSearchResult> {
  let query = new URLSearchParams({ keyword });
  let res = await fetch(`${ROOT_PATH}/modules?${query}`);
  if (res.ok) {
    return await res.json();
  }
  throw await createError(res);
}
