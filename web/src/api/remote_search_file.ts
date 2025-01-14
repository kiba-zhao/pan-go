import { withPath, withQuery, withResponds } from "fetch-utils";
import { ETagHeaderRespond, fetchMany } from "./base";
import type { ExtFSRemoteNode } from "./remote_node";
import type {
  ExtFSSearchFileSearchCondition,
  ExtFSSearchFileSearchResults,
} from "./search_file";

export interface ExtFSRemoteSearchFileAPI {
  searchExtFSRemoteSearchFileResults(
    condition: ExtFSRemoteSearchFileSearchCondition
  ): Promise<ExtFSSearchFileSearchResults>;
}

export type ExtFSRemoteSearchFileSearchCondition = ExtFSSearchFileSearchCondition &
  Pick<ExtFSRemoteNode, "peerId">;

export async function searchExtFSRemoteSearchFileResults(
  condition: ExtFSRemoteSearchFileSearchCondition
): Promise<ExtFSSearchFileSearchResults> {
  const { peerId, _start, _end, ...params } = condition;
  const condition_: Record<string, string> = {
    ...params,
  };
  if (_start !== void 0 && _start >= 0) condition_._start = _start.toString();
  if (_end !== void 0 && _end >= 0) condition_._end = _end.toString();

  const [etag, total, searchFiles] = await fetchMany(
    withPath(`extfs/remote/${peerId}/search-files`, "merge"),
    withQuery(condition_, "merge"),
    withResponds([ETagHeaderRespond], "merge")
  );
  return [etag, total, searchFiles];
}
