import { withPath, withQuery, withResponds } from "fetch-utils";
import type { ETagSearchResults } from "./base";
import { ETagHeaderRespond, fetchMany } from "./base";
import type { ExtFSNodeFile } from "./node_file";

export interface ExtFSSearchFileAPI {
  searchExtFSSearchFileResults(
    condition: ExtFSSearchFileSearchCondition
  ): Promise<ExtFSSearchFileSearchResults>;
}

export type ExtFSSearchFile = {
  id: number;
  score: number;
  tokens: string[];
} & Omit<ExtFSNodeFile, "id" | "tagQuantity" | "pendingTagQuantity">;

export type ExtFSSearchFileSearchCondition = {
  _start?: number;
  _end?: number;
  hash?: string;
  query: string;
};

export type ExtFSSearchFileSearchResults = ETagSearchResults<ExtFSSearchFile>;

export async function searchExtFSSearchFileResults(
  condition: ExtFSSearchFileSearchCondition
): Promise<ExtFSSearchFileSearchResults> {
  const { _start, _end, ...params } = condition;
  const condition_: Record<string, string> = {
    ...params,
  };
  if (_start !== void 0 && _start >= 0) condition_._start = _start.toString();
  if (_end !== void 0 && _end >= 0) condition_._end = _end.toString();

  const [etag, total, searchFiles] = await fetchMany(
    withPath("extfs/search-files", "merge"),
    withQuery(condition_, "merge"),
    withResponds([ETagHeaderRespond], "merge")
  );
  return [etag, total, searchFiles];
}
