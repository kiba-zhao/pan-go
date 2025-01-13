import { simple, withHeaderRespond } from "fetch-utils";

import jsonServerProvider from "ra-data-json-server";

const ROOT_PATH = `${import.meta.env.BASE_URL}${
  import.meta.env.VITE_API_PATH || ""
}`;

const { fetchOne, fetchMany } = simple(ROOT_PATH);
export { fetchMany, fetchOne };

// export const dataProvider = simpleDataProvider({ fetchOne, fetchMany });
export const dataProvider = jsonServerProvider(ROOT_PATH);

export const ETagHeaderRespond = withHeaderRespond("ETag");

export type ETagSearchResults<T extends any> = [string, number, T[]];
