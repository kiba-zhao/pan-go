import { ModuleItem, ModuleSearchResult } from "./api.tsx";

// Fetch uses encapsulation

type FetchAPIConfig = {
  path: string;
  query?: URLSearchParams;
  init?: RequestInit;
};

type FetchAPIConfigHandle = (config: FetchAPIConfig) => void;

type FetchAPI = (...handles: FetchAPIConfigHandle[]) => Promise<any>;

async function createError(res: Response): Promise<any> {
  const contentType = res.headers.get("Content-Type");
  if (contentType?.trim().startsWith("application/json")) {
    return await res.json();
  }
  const message = await res.text();
  return new Error(message);
}

function use(...baseHandles: FetchAPIConfigHandle[]): FetchAPI {
  return async (...handles: FetchAPIConfigHandle[]): Promise<any> => {
    const config: FetchAPIConfig = { path: "" };
    for (const handle of baseHandles) {
      handle(config);
    }
    for (const handle of handles) {
      handle(config);
    }
    let { path, query, init } = config;
    if (query) {
      path += `?${query}`;
    }
    const res = await fetch(path, init);
    if (res.ok) {
      return await res.json();
    }
    throw await createError(res);
  };
}

function withQuery(
  params?: any,
  query?: URLSearchParams
): FetchAPIConfigHandle {
  return (config: FetchAPIConfig) => {
    if (!params) {
      config.query = query;
      return;
    }
    if (!query) {
      config.query = new URLSearchParams(params);
      return;
    }
  };
}

function withBody(body: any, init?: RequestInit): FetchAPIConfigHandle {
  return (config: FetchAPIConfig) => {
    config.init = {
      ...config.init,
      ...init,
      headers: { "Content-Type": "application/json" },
      body: JSON.stringify(body),
    };
  };
}

function withPath(path: string): FetchAPIConfigHandle {
  return (config: FetchAPIConfig) => {
    config.path += path;
  };
}

type FetchAPIMethod = "GET" | "POST" | "PUT" | "PATCH" | "DELETE";

function withAPI(path: string, method?: FetchAPIMethod): FetchAPIConfigHandle {
  return (config: FetchAPIConfig) => {
    config.path += path;
    config.init = { ...config.init, method: method || "GET" };
  };
}

// API implement

const ROOT_PATH = `${import.meta.env.BASE_URL}${
  import.meta.env.VITE_API_PATH || ""
}`;

const fetchAPI = use(withPath(ROOT_PATH));

export async function SearchModules(
  keyword: string
): Promise<ModuleSearchResult> {
  return await fetchAPI(withAPI("/modules"), withQuery({ keyword }));
}

export async function GetModule(name: string): Promise<ModuleItem> {
  return await fetchAPI(withAPI(`/modules/${name}`));
}

export async function SetModuleEnabled(
  name: string,
  enabled: boolean
): Promise<ModuleItem> {
  return await fetchAPI(
    withAPI(`/modules/${name}`, "PATCH"),
    withBody({ enabled })
  );
}
