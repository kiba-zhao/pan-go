/**
 * Search File external store definition file
 *
 * @see https://18.react.dev/reference/react/useSyncExternalStore
 */
import type { Store, StoreContext } from "./Store";
import { newStore, initStoreContext, emitChange } from "./Store";

import type {
  ExtFSSearchFile,
  ExtFSSearchFileSearchCondition,
  ExtFSRemoteSearchFileSearchCondition,
  AppNode,
} from "./api";
import {
  searchExtFSRemoteSearchFileResults,
  searchExtFSSearchFileResults,
  selectAllAppNodes,
} from "./api";

type SearchFileItem = { peerId?: string } & ExtFSSearchFile;

type SearchFileStoreData = {
  files: SearchFileItem[];
  isComplete: boolean;
  errs: Record<string, any>;
};

type SearchFileContext = {
  workerId: Symbol;
  isSyncNodeComplete: boolean;
  isSyncRemoteComplete: boolean;
  abortCtrl: AbortController;
  worker: Promise<void>;
} & StoreContext<SearchFileStoreData>;

export interface SearchFileStore extends Store<SearchFileStoreData> {
  abort(reason?: any): void;
  refresh(): void;
}

/**
 * Creates a new search file store.
 *
 * Initializes a search file store with the given query and API.
 * Sets up the store context and data, and triggers an initial refresh.
 *
 * @param query - The search query string to be used for fetching files.
 * @param api - The API instance used for interacting with external services.
 * @returns A SearchFileStore with abort and refresh functionalities.
 */

export function newSearchFileStore(query: string): SearchFileStore {
  const ctx = {} as SearchFileContext;
  const data = {} as SearchFileStoreData;
  data.files = [];
  initStoreContext(ctx, data);
  refresh(ctx, query);

  const store = newStore(ctx);
  return {
    ...store,
    abort: (reason?: any) => abort(ctx, reason),
    refresh: () => refresh(ctx, query),
  };
}

/**
 * Aborts the current search operation.
 *
 * Aborts the current search operation if an abort controller has been set.
 * If a reason is provided, it will be passed to the abort controller.
 * @param reason - The reason to abort the search operation.
 */
function abort(ctx: SearchFileContext, reason?: any) {
  ctx.abortCtrl && ctx.abortCtrl.abort(reason);
}

/**
 * Refreshes the search file store.
 *
 * Resets the search file store to its initial state, and then triggers a new
 * search operation with the given query and API.
 * If the search file store is currently syncing, calling refresh will abort
 * the current search operation and trigger a new one.
 * @param query - The search query string to be used for fetching files.
 */
function refresh(ctx: SearchFileContext, query: string) {
  ctx.abortCtrl = new AbortController();
  ctx.isSyncNodeComplete = false;
  ctx.isSyncRemoteComplete = false;
  if (ctx.data.files.length > 0 || ctx.data.isComplete) {
    const data = {} as SearchFileStoreData;
    data.files = [];
    data.errs = {};
    ctx.data = data;
  }
  ctx.workerId = Symbol();
  ctx.worker = sync(ctx.workerId, ctx, query);
}

/**
 * Synchronizes the search file store with the given query and API.
 *
 * Triggers a search operation with the given query and API.
 * If the search file store is currently syncing, calling sync will abort
 * the current search operation and trigger a new one.
 * @param workerId - The symbol representing the current sync operation.
 * @param ctx - The search file store context.
 * @param query - The search query string to be used for fetching files.
 */
async function sync(workerId: Symbol, ctx: SearchFileContext, query: string) {
  if (workerId !== ctx.workerId) return;
  const generator = generateSearchFiles(query, ctx.abortCtrl.signal);
  const workers = [flushWithGenerator(workerId, ctx, generator)];

  try {
    const remotes = await selectAllAppNodes({ online: true });
    if (remotes.length > 0) {
      workers.push(syncRemotes(workerId, ctx, query, remotes));
    }
    await Promise.all(workers);
  } catch (e) {
    ctx.data.errs = { ...ctx.data.errs, "": e };
  } finally {
    if (workerId !== ctx.workerId) return;
    ctx.data = { ...ctx.data, isComplete: true };
    emitChange(ctx);
  }
}

async function syncRemotes(
  workerId: Symbol,
  ctx: SearchFileContext,
  query: string,
  remotes: AppNode[],
  offset: number = 0
) {
  if (workerId !== ctx.workerId) return;
  const { peerId } = remotes[offset];
  const generator = generateSearchFiles(query, ctx.abortCtrl.signal, peerId);
  await flushWithGenerator(workerId, ctx, generator, peerId);
  if (offset < remotes.length - 1) {
    await syncRemotes(workerId, ctx, query, remotes, offset + 1);
  }
}

/**
 * Flushes the generator with the given workerId, context, and generator.
 *
 * Flushes the generator with the given workerId, context, and generator.
 * If the workerId does not match the current workerId stored in the context,
 * the function returns without doing anything.
 * If the generator signals an abort, the function breaks out of the loop.
 * If the generator signals a completion, the function emits a change with the
 * error received from the generator.
 * Otherwise, the function appends the received search files to the current
 * search files in the context and emits a change.
 * @param workerId - The symbol representing the current sync operation.
 * @param ctx - The search file store context.
 * @param generator - The async generator used for fetching search files.
 * @param peerId - The peerId of the remote node if remote search files are
 * being fetched.
 */
async function flushWithGenerator(
  workerId: Symbol,
  ctx: SearchFileContext,
  generator: AsyncGenerator<ExtFSSearchFile[]>,
  peerId?: string
) {
  if (workerId !== ctx.workerId) return;
  while (true) {
    const { value, done } = await generator.next();
    if (workerId !== ctx.workerId) break;
    if (done) {
      if (value) {
        ctx.data.errs = { ...ctx.data.errs, [peerId || ""]: value };
        emitChange(ctx);
      }
      break;
    }
    let files = value as ExtFSSearchFile[];
    if (files.length <= 0) {
      continue;
    }
    if (peerId) {
      files = files.map((_) => ({ ..._, peerId }));
    }
    let files_ = [...ctx.data.files, ...files];
    ctx.data = { ...ctx.data, files: files_ };
    emitChange(ctx);
  }
}

/**
 * Generates search files with the given query, API, and signal.
 *
 * Generates search files with the given query, API, and signal.
 * If the search file store is currently syncing, calling generateSearchFiles will
 * abort the current search operation and trigger a new one.
 * @param query - The search query string to be used for fetching files.
 * @param signal - The abort signal used for aborting the search operation.
 * @param peerId - The peerId of the remote node if remote search files are being fetched.
 * @returns An async generator that yields an array of search files.
 */
async function* generateSearchFiles(
  query: string,
  signal: AbortSignal,
  peerId?: string
): AsyncGenerator<ExtFSSearchFile[]> {
  const condition = {} as ExtFSSearchFileSearchCondition;
  condition.query = query;
  condition._start = 0;
  condition._end = 100;
  let err;

  let hash: string;
  let total: number;
  let files: ExtFSSearchFile[];
  while (true) {
    try {
      if (peerId) {
        [hash, total, files] = await searchExtFSRemoteSearchFileResults(
          peerId,
          condition as ExtFSRemoteSearchFileSearchCondition,
          { signal }
        );
      } else {
        [hash, total, files] = await searchExtFSSearchFileResults(condition, {
          signal,
        });
      }

      yield files;
      if (total === 0) {
        break;
      }

      if (condition.hash === void 0) {
        condition.hash = hash;
      } else if (condition.hash !== hash) {
        throw new Error("hash not match");
      }

      condition._start = condition._start + files.length;
      if (total > 0 && condition._start >= total) {
        break;
      }
      condition._end = condition._end + files.length;
      if (total > 0 && condition._end > total) {
        condition._end = total;
      }

      if (total < 0) {
        await new Promise((resolve) => setTimeout(resolve, 1500));
      }
    } catch (e) {
      err = e;
      break;
    }
  }

  return err;
}
