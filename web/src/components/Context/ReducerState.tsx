import { createContext, useContext, useReducer } from "react";

import type { Dispatch, Context as ReactContext, ReactNode } from "react";

type ReducerState<T extends any> = T;
type ReducerStateAction<T extends any> = ReducerState<T>;

const reducer = <T extends any>(
  _: ReducerState<T>,
  action: ReducerStateAction<T>
) => {
  return action;
};

type ReducerContext<T extends any> = ReactContext<ReducerState<T> | null>;
type ReducerDispatchContext<T extends any> = ReactContext<Dispatch<
  ReducerStateAction<T>
> | null>;

type ReducerContextOpts<T extends any> = [
  ReducerContext<T>,
  ReducerDispatchContext<T>
];
export function createReducerContext<T extends any>(): ReducerContextOpts<T> {
  const context = createContext<ReducerState<T> | null>(null);
  const dispatchContext = createContext<Dispatch<ReducerStateAction<T>> | null>(
    null
  );
  return [context, dispatchContext];
}

const [DefaultReducerContext, DefaultReducerDispatchContext] =
  createReducerContext();

export const ReducerStateProvider = <T extends any>({
  initialState,
  children,
  opts,
}: {
  initialState: ReducerState<T>;
  children: ReactNode;
  opts?: ReducerContextOpts<T>;
}) => {
  const [state, dispatch] = useReducer(reducer, initialState);
  const ReducerContext = opts?.[0] || DefaultReducerContext;
  const ReducerDispatchContext = opts?.[1] || DefaultReducerDispatchContext;
  return (
    <ReducerContext.Provider value={state as ReducerState<T>}>
      <ReducerDispatchContext.Provider value={dispatch}>
        {children}
      </ReducerDispatchContext.Provider>
    </ReducerContext.Provider>
  );
};

type useReducerStateReturnType<T extends any> = [T, Dispatch<T>];
export const useReducerState = <T extends any>(
  opts?: ReducerContextOpts<T>
): useReducerStateReturnType<T> => [
  useReducerStateValue<T>(opts ? opts[0] : void 0),
  useReducerStateValueSetter<T>(opts ? opts[1] : void 0),
];

export const useReducerStateValue = <T extends any>(
  ctx?: ReducerContext<T>
): T => {
  const state = ctx ? useContext(ctx) : useContext(DefaultReducerContext);
  return state as T;
};

export const useReducerStateValueSetter = <T extends any>(
  ctx?: ReducerDispatchContext<T>
): Dispatch<T> => {
  const dispatch = ctx
    ? useContext(ctx)
    : useContext(DefaultReducerDispatchContext);
  return dispatch ? dispatch : noop;
};

function noop() {}
