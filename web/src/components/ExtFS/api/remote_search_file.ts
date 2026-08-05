/**
 * ExtFSRemoteSearchFile API Definition File
 */
import {
  withPath,
  withQuery,
  withResponds,
  withRequestInit,
} from "fetch-utils";
import type { FetchContextHandle } from "fetch-utils";
import { ETagHeaderRespond, fetchMany } from "../../../lib/fetch";
import type { AppNode } from "../../AppNode/api";
import type {
  ExtFSSearchFileSearchCondition,
  ExtFSSearchFileSearchResults,
} from "./search_file";

export type ExtFSRemoteSearchFileSearchCondition =
  ExtFSSearchFileSearchCondition;

export type ExtFSRemoteSearchFileSearchOpts = {
  signal?: AbortSignal;
};

/**
 * Search remote search files with condition.
 * @param peerId The peer id.
 * @param condition The condition to search remote search files.
 * @param opts Options with signal.
 * @returns The search result.
 */
export async function searchExtFSRemoteSearchFileResults(
  peerId: AppNode["peerId"],
  condition: ExtFSRemoteSearchFileSearchCondition,
  opts?: ExtFSRemoteSearchFileSearchOpts,
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
    withPath(`extfs/remotes/${peerId}/search-files`, "merge"),
    withQuery(condition_, "merge"),
    withResponds([ETagHeaderRespond], "merge"),
  );
  return [etag, total, searchFiles];
}
