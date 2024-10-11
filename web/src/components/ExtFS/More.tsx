import {
  createReducerContext,
  ReducerStateProvider,
  useReducerState,
} from "../Context/ReducerState";

import { useExtFS } from "./State";

import type { MouseEvent, ReactNode } from "react";
import { Fragment, useMemo } from "react";
import type { To } from "react-router-dom";
import { Link as RouterLink } from "react-router-dom";

import AddCircleIcon from "@mui/icons-material/AddCircle";
import HelpIcon from "@mui/icons-material/Help";
import MoreVertIcon from "@mui/icons-material/MoreVert";
import SettingsIcon from "@mui/icons-material/Settings";
import IconButton from "@mui/material/IconButton";
import ListItemIcon from "@mui/material/ListItemIcon";
import ListItemText from "@mui/material/ListItemText";
import Menu from "@mui/material/Menu";
import MenuItem from "@mui/material/MenuItem";

import { useTranslate } from "react-admin";

type MoreState = {
  anchorEl: null | HTMLElement;
};

const MoreContext = createReducerContext<MoreState>();
const useMore = () => useReducerState(MoreContext);

const MoreMenu = ({ children }: { children: ReactNode }) => {
  const [state, setState] = useMore();
  const open = Boolean(state.anchorEl);

  const handleClick = (event: MouseEvent<HTMLButtonElement>) => {
    setState({ anchorEl: event.currentTarget });
  };
  const handleClose = () => {
    setState({ anchorEl: null });
  };

  return (
    <Fragment>
      <IconButton onClick={handleClick}>
        <MoreVertIcon />
      </IconButton>
      <Menu anchorEl={state.anchorEl} open={open} onClose={handleClose}>
        {children}
      </Menu>
    </Fragment>
  );
};

export const More = ({ children }: { children: ReactNode }) => {
  return (
    <ReducerStateProvider<MoreState>
      initialState={{ anchorEl: null }}
      opts={MoreContext}
    >
      <MoreMenu>{children}</MoreMenu>
    </ReducerStateProvider>
  );
};

type MoreItemContentProps = {
  icon?: ReactNode;
  children?: ReactNode;
};
const MoreItemContent = ({ icon, children }: MoreItemContentProps) => {
  const [contentElement, iconElement] = useMemo(() => {
    if (!icon) {
      return [children];
    }
    return [
      <ListItemText>{children}</ListItemText>,
      <ListItemIcon>{icon}</ListItemIcon>,
    ];
  }, [children, icon]);

  if (!iconElement) {
    return children;
  }
  return (
    <Fragment>
      {iconElement}
      {contentElement}
    </Fragment>
  );
};

export type MoreLinkItemProps = { to: To } & MoreItemContentProps;
export const MoreLinkItem = ({ icon, children, to }: MoreLinkItemProps) => {
  const [_, setMenuState] = useMore();
  const handleItemClick = () => {
    setMenuState({ anchorEl: null });
  };

  return (
    <MenuItem onClick={handleItemClick} component={RouterLink} to={to}>
      <MoreItemContent icon={icon}>{children}</MoreItemContent>
    </MenuItem>
  );
};

export type MoreButtonItemProps = {
  onClick?: () => void;
} & MoreItemContentProps;
export const MoreButtonItem = ({
  onClick,
  icon,
  children,
}: MoreButtonItemProps) => {
  const [_, setMoreState] = useMore();
  const handleItemClick = () => {
    setMoreState({ anchorEl: null });
    onClick && onClick();
  };
  return (
    <MenuItem onClick={handleItemClick}>
      <MoreItemContent icon={icon}>{children}</MoreItemContent>
    </MenuItem>
  );
};

export const MoreNewItem = () => {
  const t = useTranslate();
  const [state, _] = useExtFS();

  if (state.parentItems.length <= 0) return void 0;
  const url = state.parentItems.at(-1)?.newUrl;
  if (!url) return void 0;
  return (
    <MoreLinkItem icon={<AddCircleIcon fontSize="small" />} to={url}>
      {t("custom.button.new")}
    </MoreLinkItem>
  );
};

export const MoreSettingsItem = () => {
  const t = useTranslate();
  const [state, _] = useExtFS();

  if (state.parentItems.length <= 0) return void 0;
  const url = state.parentItems.at(-1)?.settingsUrl;
  if (!url) return void 0;
  return (
    <MoreLinkItem icon={<SettingsIcon fontSize="small" />} to={url}>
      {t("custom.button.new")}
    </MoreLinkItem>
  );
};

export const MoreHelpItem = () => {
  // const t = useTranslate();
  const handleClose = () => {
    // TODO: show help info
  };
  return (
    <MoreButtonItem onClick={handleClose} icon={<HelpIcon fontSize="small" />}>
      Help
    </MoreButtonItem>
  );
};
