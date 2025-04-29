/**
 * Reducer Context Provider Definition File
 */
import { createContext, useContext, useReducer } from "react";

import type { Dispatch, Context as ReactContext, ReactNode } from "react";

type ReducerState<T extends any> = T;
type ReducerStateAction<T extends any> = ReducerState<T>;

const defaultReducer = <T extends any>(
  _: ReducerState<T>,
  action: ReducerStateAction<T>
) => {
  return action;
};

export const mergeReducer = <T extends object>(
  state: ReducerState<T>,
  action: ReducerStateAction<T>
): ReducerState<T> => ({
  ...state,
  ...action,
});

type ReducerContext<T extends any> = ReactContext<ReducerState<T> | null>;
type ReducerDispatchContext<T extends any> = ReactContext<Dispatch<
  ReducerStateAction<T>
> | null>;

type ReducerContextOpts<T extends any> = [
  ReducerContext<T>,
  ReducerDispatchContext<T>
];
/**
 * Create a reducer context pair.
 *
 * @template T - Type of the reducer state.
 * @returns A pair of reducer context and dispatch context.
 */
export function createReducerContext<T extends any>(): ReducerContextOpts<T> {
  const context = createContext<ReducerState<T> | null>(null);
  const dispatchContext = createContext<Dispatch<ReducerStateAction<T>> | null>(
    null
  );
  return [context, dispatchContext];
}

const [DefaultReducerContext, DefaultReducerDispatchContext] =
  createReducerContext();

/**
 * Provider for reducer state.
 *
 * @template T - Type of the reducer state.
 * @param {Object} props
 * @param {ReducerState<T>} props.initialState - Initial value of the reducer state.
 * @param {ReactNode} props.children - Children components to wrap.
 * @param {(state: ReducerState<T>, action: ReducerStateAction<T>) => T} [props.reducer] - Reducer function.
 * @param {ReducerContextOpts<T> | undefined} [props.opts] - Optional context pair.
 * @returns {ReactElement} Element with reducer state and dispatch context.
 */
export const ReducerStateProvider = <T extends any>({
  initialState,
  children,
  reducer,
  opts,
}: {
  initialState: ReducerState<T>;
  children: ReactNode;
  reducer?: (state: ReducerState<T>, action: ReducerStateAction<T>) => T;
  opts?: ReducerContextOpts<T>;
}) => {
  const [state, dispatch] = useReducer(reducer || defaultReducer, initialState);
  const ReducerContext =
    opts?.[0] || (DefaultReducerContext as ReducerContext<T>);
  const ReducerDispatchContext =
    opts?.[1] || (DefaultReducerDispatchContext as ReducerDispatchContext<T>);
  return (
    <ReducerContext.Provider value={state as ReducerState<T>}>
      <ReducerDispatchContext.Provider value={dispatch}>
        {children}
      </ReducerDispatchContext.Provider>
    </ReducerContext.Provider>
  );
};

type useReducerStateReturnType<T extends any> = [T, Dispatch<T>];
/**
 * Hook to get the current reducer state and dispatch.
 *
 * @template T - Type of the reducer state.
 * @param {ReducerContextOpts<T> | undefined} [opts] - Optional context pair.
 * @returns {[T, Dispatch<T>]} A pair of current reducer state and dispatch function.
 */
export const useReducerState = <T extends any>(
  opts?: ReducerContextOpts<T>
): useReducerStateReturnType<T> => [
  useReducerStateValue<T>(opts ? opts[0] : void 0),
  useReducerStateValueSetter<T>(opts ? opts[1] : void 0),
];

/**
 * Hook to get the current reducer state value.
 *
 * @template T - Type of the reducer state.
 * @param {ReducerContext<T> | undefined} [ctx] - Optional context.
 * @returns {T} The current reducer state value.
 */
export const useReducerStateValue = <T extends any>(
  ctx?: ReducerContext<T>
): T => {
  const state = ctx ? useContext(ctx) : useContext(DefaultReducerContext);
  return state as T;
};

/**
 * Hook to get the current reducer state setter.
 *
 * @template T - Type of the reducer state.
 * @param {ReducerDispatchContext<T> | undefined} [ctx] - Optional context.
 * @returns {Dispatch<T>} The current reducer state setter function.
 */
export const useReducerStateValueSetter = <T extends any>(
  ctx?: ReducerDispatchContext<T>
): Dispatch<T> => {
  const dispatch = ctx
    ? useContext(ctx)
    : useContext(DefaultReducerDispatchContext);
  return dispatch ? dispatch : noop;
};

function noop() {}
