/**
 * ExtFSSearchItem API Definition File
 */
import { withJSONBody, withMethod, withPath, withQuery } from "fetch-utils";
import { fetchMany, fetchOne } from "./base";

export interface ExtFSSearchItemAPI {
  searchExtFSSearchItems(
    condition: ExtFSSearchItemSearchCondition
  ): Promise<ExtFSSearchItem[]>;
  deleteExtFSSearchItem(id: ExtFSSearchItem["id"]): Promise<void>;
  saveExtFSSearchItem(
    fields: ExtFSSearchItemFields,
    id?: ExtFSSearchItem["id"]
  ): Promise<ExtFSSearchItem>;
}

export type ExtFSSearchItemSearchCondition = {
  q?: string;
  limit?: number;
};

export type ExtFSSearchItem = {
  id: number;
  query: string;
  createdAt: string;
  updatedAt: string;
};

export type ExtFSSearchItemFields = Omit<
  ExtFSSearchItem,
  "id" | "createdAt" | "updatedAt" | "deletedAt"
>;

/**
 * Search ExtFSSearchItem with condition
 *
 * @param condition
 * @returns array of ExtFSSearchItem
 */
export async function searchExtFSSearchItems(
  condition: ExtFSSearchItemSearchCondition
): Promise<ExtFSSearchItem[]> {
  const { limit, ...queries } = condition;
  const _limit = limit != void 0 ? limit.toString() : void 0;
  const [_, searchItems] = await fetchMany(
    withPath("extfs/search-items", "merge"),
    withQuery(_limit ? { ...queries, _end: _limit } : queries, "merge")
  );
  return searchItems;
}

/**
 * Delete ExtFSSearchItem with id
 *
 * @param id id of ExtFSSearchItem
 */
export async function deleteExtFSSearchItem(id: ExtFSSearchItem["id"]) {
  return await fetchOne(
    withPath(`extfs/search-items/${id}`, "merge"),
    withMethod("DELETE")
  );
}

/**
 * update ExtFSSearchItem with id
 *
 * @param fields fields of ExtFSSearchItem
 * @param id id of ExtFSSearchItem
 * @returns ExtFSSearchItem
 */
export async function saveExtFSSearchItem(
  fields: ExtFSSearchItemFields,
  id?: ExtFSSearchItem["id"]
): Promise<ExtFSSearchItem> {
  return await fetchOne(
    withPath(`extfs/search-items/${id || ""}`, "merge"),
    withMethod(id ? "PATCH" : "POST"),
    withJSONBody(fields)
  );
}
