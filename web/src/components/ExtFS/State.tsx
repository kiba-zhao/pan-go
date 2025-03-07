import {
  createReducerContext,
  ReducerStateProvider,
  useReducerState,
} from "../Context/ReducerState";

import type { ReactNode } from "react";

export type ExtFSSingleState = {
  mode: string;
  queryKeyList: Array<string[]>;
};
export type ExtFSParentItem = {
  name: string;
  state: ExtFSSingleState;
};
export type ExtFSState = {
  parentItems: ExtFSParentItem[];
} & ExtFSSingleState;

const ExtFSContext = createReducerContext<ExtFSState>();

/**
 * Hook to access the ExtFS state.
 * @returns {ExtFSState} The current state of the ExtFS context.
 */
export const useExtFS = () => useReducerState(ExtFSContext);

/**
 * Provider for ExtFS state.
 *
 * @param {{ value: ExtFSState, children: ReactNode }} props
 * @prop {ExtFSState} value - Initial value of the ExtFS state.
 * @prop {ReactNode} children - Children components to render.
 * @returns {ReactElement} Element with reducer state and dispatch context.
 */
export const ExtFSProvider = ({
  value,
  children,
}: {
  value: ExtFSState;
  children: ReactNode;
}) => {
  return (
    <ReducerStateProvider<ExtFSState> initialState={value} opts={ExtFSContext}>
      {children}
    </ReducerStateProvider>
  );
};
