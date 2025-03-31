/**
 * ExtFSSearchFile API Definition File
 */
import {
  withPath,
  withQuery,
  withResponds,
  withRequestInit,
} from "fetch-utils";
import type { FetchContextHandle } from "fetch-utils";
import type { ETagSearchResults } from "./base";
import { ETagHeaderRespond, fetchMany } from "./base";
import type { ExtFSNodeFile } from "./node_file";

export interface ExtFSSearchFileAPI {
  searchExtFSSearchFileResults(
    condition: ExtFSSearchFileSearchCondition,
    opts?: ExtFSSearchFileSearchOpts
  ): Promise<ExtFSSearchFileSearchResults>;
}

export type ExtFSSearchFile = {
  id: number;
  score: number;
  tokens: string[];
} & Omit<ExtFSNodeFile, "id">;

export type ExtFSSearchFileSearchCondition = {
  _start?: number;
  _end?: number;
  hash?: string;
  query: string;
};

export type ExtFSSearchFileSearchResults = ETagSearchResults<ExtFSSearchFile>;

export type ExtFSSearchFileSearchOpts = {
  signal?: AbortSignal;
};

/**
 * Search ExtFSSearchFile with condition
 *
 * @param condition condition of ExtFSSearchFile
 * @param opts Options with signal
 * @returns Search Result
 */
export async function searchExtFSSearchFileResults(
  condition: ExtFSSearchFileSearchCondition,
  opts?: ExtFSSearchFileSearchOpts
): Promise<ExtFSSearchFileSearchResults> {
  const { _start, _end, ...params } = condition;
  const condition_: Record<string, string> = {
    ...params,
  };
  if (_start !== void 0 && _start >= 0) condition_._start = _start.toString();
  if (_end !== void 0 && _end >= 0) condition_._end = _end.toString();

  const handles = [] as FetchContextHandle[];
  if (opts && opts.signal) {
    handles.push(withRequestInit({ signal: opts.signal }));
  }

  const [etag, total, searchFiles] = await fetchMany(
    ...handles,
    withPath("extfs/search-files", "merge"),
    withQuery(condition_, "merge"),
    withResponds([ETagHeaderRespond], "merge")
  );
  return [etag, total, searchFiles];
}
