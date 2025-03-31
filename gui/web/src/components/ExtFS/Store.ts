export type StoreListener = () => void;
export interface Store<T extends any> {
  subscribe(listener: StoreListener): StoreListener;
  getSnapshot(): T;
}

export type StoreContext<T extends any> = {
  listeners: StoreListener[];
  data: T;
};

/**
 * Creates a store that wraps a given context.
 * The context should be an object with the following properties:
 * - `listeners`: an array of listeners to be notified when the store changes
 * - `data`: the current state of the store
 * The returned store has the following methods:
 * - `subscribe`: registers a listener to be notified when the store changes
 * - `getSnapshot`: returns the current state of the store
 */
export function newStore<T extends any>(ctx: StoreContext<T>): Store<T> {
  return {
    subscribe: (listener) => subscribe(ctx, listener),
    getSnapshot: () => getSnapshot(ctx),
  };
}

/**
 * Registers a listener to be notified when the store changes.
 * The listener is added to the store context's `listeners` array.
 * The returned function can be used to remove the listener from the store context.
 * @param {StoreContext<T>} ctx - The store context to subscribe to.
 * @param {StoreListener} listener - The listener to be notified when the store changes.
 * @returns {StoreListener} - A function that can be used to remove the listener from the store context.
 */
export function subscribe<T extends any>(
  ctx: StoreContext<T>,
  listener: StoreListener
): StoreListener {
  ctx.listeners = [...ctx.listeners, listener];
  return () => {
    ctx.listeners = ctx.listeners.filter((l) => l !== listener);
  };
}

/**
 * Returns the current state of the store.
 * @param {StoreContext<T>} ctx - The store context to get the snapshot from.
 * @returns {T} The current state of the store.
 */
export function getSnapshot<T extends any>(ctx: StoreContext<T>): T {
  return ctx.data;
}

/**
 * Emits a change notification to all the listeners subscribed to the store.
 * This function iterates over the `listeners` array of the context and calls each listener.
 * @param {StoreContext<T>} ctx - The store context to emit the change from.
 */
export function emitChange<T extends any>(ctx: StoreContext<T>) {
  for (const listener of ctx.listeners) {
    listener();
  }
}

/**
 * Initializes the store context with provided data.
 *
 * @template T - Type of the data in the store context.
 * @param {StoreContext<T>} ctx - The store context to initialize.
 * @param {T} data - The initial data to set in the store context.
 *
 * This function sets the `data` property of the context to the provided `data`
 * and clears the `listeners` array, effectively resetting the store context.
 */

export function initStoreContext<T extends any>(ctx: StoreContext<T>, data: T) {
  ctx.data = data;
  ctx.listeners = [];
}
