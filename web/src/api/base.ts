/**
 * Base API Definition File
 */
import { simple, withHeaderRespond } from "fetch-utils";

import jsonServerProvider from "ra-data-json-server";

const ROOT_PATH = `${import.meta.env.BASE_URL}${
  import.meta.env.VITE_API_PATH || ""
}`;

const { fetchOne, fetchMany } = simple(ROOT_PATH);
export { fetchMany, fetchOne };


/**
 * Data Provider for react-admin
 * 
 * @see https://github.com/marmelab/react-admin/tree/master/packages/ra-data-json-server
 */
export const dataProvider = jsonServerProvider(ROOT_PATH);

export const ETagHeaderRespond = withHeaderRespond("ETag");

export type ETagSearchResults<T extends any> = [string, number, T[]];

/**
 * get root path
 * 
 * @returns root path of API
 */
export const RootPath = () => ROOT_PATH;

export interface BaseAPI {
  RootPath: () => string;
}
