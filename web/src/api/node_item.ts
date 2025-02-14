import { withJSONBody, withMethod, withPath } from "fetch-utils";
import { fetchMany, fetchOne } from "./base";

export interface ExtFSNodeItemAPI {
  selectAllExtFSNodeItems(): Promise<ExtFSNodeItem[]>;
  saveExtFSNodeItem(
    fields: ExtFSNodeItemFields,
    id?: ExtFSNodeItem["id"]
  ): Promise<ExtFSNodeItem>;
  selectExtFSNodeItem(id: ExtFSNodeItem["id"]): Promise<ExtFSNodeItem>;
  deleteExtFSNodeItem(id: ExtFSNodeItem["id"]): Promise<void>;
}

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
  | "tagQuantity"
  | "pendingTagQuantity"
>;

export async function selectAllExtFSNodeItems(): Promise<ExtFSNodeItem[]> {
  const [_, nodeItems] = await fetchMany(withPath("extfs/node-items", "merge"));
  return nodeItems;
}

export async function saveExtFSNodeItem(
  fields: ExtFSNodeItemFields,
  id?: ExtFSNodeItem["id"]
): Promise<ExtFSNodeItem> {
  return await fetchOne(
    withPath(`extfs/node-items${id ? `/${id}` : ""}`, "merge"),
    withMethod(id ? "PATCH" : "POST"),
    withJSONBody(fields)
  );
}

export async function selectExtFSNodeItem(
  id: ExtFSNodeItem["id"]
): Promise<ExtFSNodeItem> {
  return await fetchOne(withPath(`extfs/node-items/${id}`, "merge"));
}

export async function deleteExtFSNodeItem(id: ExtFSNodeItem["id"]) {
  return await fetchOne(
    withPath(`extfs/node-items/${id}`, "merge"),
    withMethod("DELETE")
  );
}
