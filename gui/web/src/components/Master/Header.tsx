import type { ReactNode } from "react";
import { useEffect } from "react";

import {
  createReducerContext,
  ReducerStateProvider,
  useReducerState,
  useReducerStateValue,
} from "../Common/ReducerState";

import AppBar from "@mui/material/AppBar";
import Box from "@mui/material/Box";
import Toolbar from "@mui/material/Toolbar";

type HeaderState = {
  title?: ReactNode;
  tools?: ReactNode;
  addons?: ReactNode;
};

const HeaderContextOpts = createReducerContext<HeaderState>();
const useHeader = () => useReducerState(HeaderContextOpts);
const useHeaderState = () => useReducerStateValue(HeaderContextOpts[0]);

export const HeaderProvider = ({ children }: { children: ReactNode }) => {
  const initialState = { children: null } as HeaderState;
  return (
    <ReducerStateProvider initialState={initialState} opts={HeaderContextOpts}>
      {children}
    </ReducerStateProvider>
  );
};

export const HeaderView = ({ children }: { children?: ReactNode }) => {
  const { title, tools, addons } = useHeaderState();
  return (
    <AppBar component="nav" position="static">
      <Toolbar>
        {title}
        <Box sx={{ flexGrow: 1 }} />
        {tools}
        {children}
        {addons}
      </Toolbar>
    </AppBar>
  );
};

type PageHeaderProps = { partial?: boolean } & HeaderState;
export const PageHeader = ({ partial, ...state_ }: PageHeaderProps) => {
  const [state, dispatch] = useHeader();
  const { title, tools, addons } = state_;

  useEffect(() => {
    dispatch(partial ? { ...state, ...state_ } : state_);
  }, [partial, title, tools, addons]);
  return <></>;
};
