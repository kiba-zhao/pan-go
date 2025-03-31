import {
  withHeaderRespond,
  newFetch,
  withPath,
  withRespond,
  withResponds,
  respondJSON,
  FetchError,
} from "fetch-utils";

import jsonServerProvider from "ra-data-json-server";

/**
 * API root path
 */
const ROOT_PATH = `${import.meta.env.BASE_URL}${
  import.meta.env.VITE_API_PATH || ""
}`;

/**
 * Custom fetch respond handle with json
 *
 * @param {Response} res fetch response
 * @returns {Promise<T>} json response
 * @throws {Error} custom http error
 */
async function respond<T extends any>(res: Response): Promise<T> {
  try {
    const json = await respondJSON(res);
    return json;
  } catch (error) {
    if (error instanceof FetchError) {
      error.name = `${error.name}.${res.status}`;
    }
    throw error;
  }
}

/**
 * X-Total-Count header respond handle
 */
const XTotalCountRespond = withHeaderRespond("X-Total-Count", Number);

const baseHandles = [withPath(ROOT_PATH), withRespond(respond)];

/**
 * Custom fetch one entity
 */
export const fetchOne = newFetch(...baseHandles);

/**
 * Custom fetch many entities
 */
export const fetchMany = newFetch(
  ...baseHandles,
  withResponds([XTotalCountRespond], "merge")
);

/**
 * JSON Server data provider
 */
export const dataProvider = jsonServerProvider(ROOT_PATH);

/**
 * ETag header respond handle
 */
export const ETagHeaderRespond = withHeaderRespond("ETag");

export type ETagSearchResults<T extends any> = [string, number, T[]];

/**
 * Get API root path
 *
 * @returns {string} API root path
 */
export const RootPath = () => ROOT_PATH;

export interface BaseAPI {
  RootPath: () => string;
}
