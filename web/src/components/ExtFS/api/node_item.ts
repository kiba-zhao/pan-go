/**
 * ExtFSNodeItem API Definition File
 */
import { withJSONBody, withMethod, withPath } from "fetch-utils";
import { fetchMany, fetchOne } from "../../../lib/api";

export type ExtFSNodeItem = {
  id: number;
  name: string;
  filePath: string;
  fileType: "F" | "D";
  mimeType: string;
  size: number;
  enabled: boolean;
  available: boolean;
  createdAt: string;
  updatedAt: string;
  deletedAt: string;
};

export type ExtFSNodeItemFields = Omit<
  ExtFSNodeItem,
  | "id"
  | "createdAt"
  | "updatedAt"
  | "deletedAt"
  | "available"
  | "fileType"
  | "mimeType"
  | "size"
>;

/**
 * Select all ExtFS node items.
 *
 * @returns A list of node items.
 */
export async function selectAllExtFSNodeItems(): Promise<ExtFSNodeItem[]> {
  const [_, nodeItems] = await fetchMany(withPath("extfs/node-items", "merge"));
  return nodeItems;
}

/**
 * Save a ExtFS node item with `id` or create a new one.
 *
 * @param fields The node item fields to save.
 * @param id The node item id. If not provided, a new node item will be created.
 * @returns The saved node item.
 */
export async function saveExtFSNodeItem(
  fields: ExtFSNodeItemFields,
  id?: ExtFSNodeItem["id"],
): Promise<ExtFSNodeItem> {
  return await fetchOne(
    withPath(`extfs/node-items${id ? `/${id}` : ""}`, "merge"),
    withMethod(id ? "PATCH" : "POST"),
    withJSONBody(fields),
  );
}

/**
 * Select a ExtFS node item with `id`.
 *
 * @param id The node item id.
 * @returns node item.
 */
export async function selectExtFSNodeItem(
  id: ExtFSNodeItem["id"],
): Promise<ExtFSNodeItem> {
  return await fetchOne(withPath(`extfs/node-items/${id}`, "merge"));
}

/**
 * Delete a ExtFS node item with `id`.
 *
 * @param id The node item id.
 */
export async function deleteExtFSNodeItem(id: ExtFSNodeItem["id"]) {
  return await fetchOne(
    withPath(`extfs/node-items/${id}`, "merge"),
    withMethod("DELETE"),
  );
}
