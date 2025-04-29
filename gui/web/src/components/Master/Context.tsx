import type { ComponentType, ReactNode } from "react";
import { useEffect } from "react";

import {
  createReducerContext,
  ReducerStateProvider,
  useReducerStateValue,
  useReducerStateValueSetter,
} from "../Common/ReducerState";
type State = {
  Component: ComponentType<{ children?: ReactNode }> | null;
};

const ContextOpts = createReducerContext<State>();
const useStateValue = () => useReducerStateValue(ContextOpts[0]);
const useStateValueDispatch = () => useReducerStateValueSetter(ContextOpts[1]);

export const Provider = ({ children }: { children: ReactNode }) => {
  const initialState = { Component: null } as State;
  return (
    <ReducerStateProvider initialState={initialState} opts={ContextOpts}>
      {children}
    </ReducerStateProvider>
  );
};

export const ProviderView = ({ children }: { children?: ReactNode }) => {
  const { Component } = useStateValue();
  if (!Component) return children;
  return <Component>{children}</Component>;
};

type PageProviderProps = State;
export const PageProvider = ({ Component }: PageProviderProps) => {
  const dispatch = useStateValueDispatch();
  useEffect(() => {
    dispatch({ Component });
  }, [Component]);
  return null;
};
