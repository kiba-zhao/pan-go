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

export const useExtFS = () => useReducerState(ExtFSContext);

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
