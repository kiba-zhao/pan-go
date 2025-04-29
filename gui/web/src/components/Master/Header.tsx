import type { ReactNode } from "react";
import { useEffect } from "react";

import {
  createReducerContext,
  ReducerStateProvider,
  useReducerState,
  useReducerStateValue,
} from "../Common/ReducerState";

import Box from "@mui/material/Box";
import AppBar from "@mui/material/AppBar";
import Toolbar from "@mui/material/Toolbar";
import Typography from "@mui/material/Typography";

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

export const PageHeaderTitle = ({ title }: { title: string }) => {
  const children = (
    <Typography noWrap variant="h6" component="div">
      {title}
    </Typography>
  );
  return <PageHeader title={children} partial />;
};

export const PageHeaderTools = ({ children }: { children?: ReactNode }) => {
  return <PageHeader tools={children} partial />;
};

export const PageHeaderAddons = ({ children }: { children?: ReactNode }) => {
  return <PageHeader addons={children} partial />;
};
