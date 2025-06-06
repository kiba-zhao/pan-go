/**
 * Item More Component Definition File
 */
import {
  createReducerContext,
  ReducerStateProvider,
  useReducerState,
} from "./ReducerState";

import type { MouseEvent, ReactNode } from "react";
import { Fragment, useMemo } from "react";
import type { To } from "react-router-dom";
import { Link as RouterLink } from "../Route/Router";
import type { MenuItemOwnProps } from "@mui/material/MenuItem";

import MoreVertIcon from "@mui/icons-material/MoreVert";
import IconButton from "@mui/material/IconButton";
import ListItemIcon from "@mui/material/ListItemIcon";
import ListItemText from "@mui/material/ListItemText";
import Menu from "@mui/material/Menu";
import MenuItem from "@mui/material/MenuItem";

type MoreState = {
  anchorEl: null | HTMLElement;
};

const MoreContext = createReducerContext<MoreState>();
const useMore = () => useReducerState(MoreContext);

/**
 * A component that renders a menu with the given children.
 * The menu is opened/closed using the MoreVertIcon.
 * The menu is rendered as a MUI Menu and the children
 * are rendered as MUI MenuItems.
 *
 * @example
 * <MoreMenu>
 *   <MenuItem>Item 1</MenuItem>
 *   <MenuItem>Item 2</MenuItem>
 * </MoreMenu>
 *
 * @param {ReactNode} children - The children to be rendered as menu items.
 * @returns {ReactElement} A React element representing the menu.
 */
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
      <IconButton onClick={handleClick} color="inherit">
        <MoreVertIcon />
      </IconButton>
      <Menu anchorEl={state.anchorEl} open={open} onClose={handleClose}>
        {children}
      </Menu>
    </Fragment>
  );
};

/**
 * A component that renders a menu with the given children.
 * The menu is opened/closed using the MoreVertIcon.
 * The menu is rendered as a MUI Menu and the children
 * are rendered as MUI MenuItems.
 *
 * @example
 * <More>
 *   <MenuItem>Item 1</MenuItem>
 *   <MenuItem>Item 2</MenuItem>
 * </More>
 *
 * @param {ReactNode} children - The children to be rendered as menu items.
 * @returns {ReactElement} A React element representing the menu.
 */
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
/**
 * A component that renders a content element and an icon element as children.
 * If the icon is not provided, the content element is returned as is.
 * If the icon is provided, the content element is wrapped in a ListText
 * and the icon is rendered as a ListIcon.
 *
 * @example
 * <More>
 *   <MenuItem>
 *     <MoreItemContent icon={<SettingsIcon />}>
 *       Item 1
 *     </MoreItemContent>
 *   </MenuItem>
 *   <MenuItem>
 *     Item 2
 *   </MenuItem>
 * </More>
 *
 * @param {ReactNode} icon - The icon to be rendered as a ListIcon.
 * @param {ReactNode} children - The content to be rendered as a ListText.
 * @returns {ReactElement} A React element representing the menu item content.
 */
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

export type MoreLinkItemProps = { to: To } & Pick<
  MenuItemOwnProps,
  "selected"
> &
  MoreItemContentProps;
/**
 * A component that renders a menu item with a link to a given route.
 * This component is used to render a menu item in the More component.
 *
 * @param {To} to - The target route of the link.
 * @param {ReactNode} icon - The icon to be rendered as a ListIcon.
 * @param {ReactNode} children - The content to be rendered as a ListText.
 *
 * @returns {ReactElement} A React element representing the menu item.
 *
 * @example
 * <More>
 *   <MenuItem>
 *     <MoreLinkItemContent icon={<SettingsIcon />} to="/extfs/settings">
 *       Settings
 *     </MoreLinkItemContent>
 *   </MenuItem>
 * </More>
 */
export const MoreLinkItem = ({
  icon,
  children,
  to,
  selected,
}: MoreLinkItemProps) => {
  const [_, setMenuState] = useMore();
  const handleItemClick = () => {
    setMenuState({ anchorEl: null });
  };

  return (
    <MenuItem
      onClick={handleItemClick}
      component={RouterLink}
      to={to}
      selected={selected}
    >
      <MoreItemContent icon={icon}>{children}</MoreItemContent>
    </MenuItem>
  );
};

export type MoreButtonItemProps = {
  onClick?: () => void;
} & MoreItemContentProps;
/**
 * A component that renders a menu item with a button.
 * This component is used to render a menu item in the More component.
 *
 * @param {ReactNode} icon - The icon to be rendered as a ListIcon.
 * @param {ReactNode} children - The content to be rendered as a ListText.
 * @param {() => void} onClick - An optional function to be called when the button is clicked.
 *
 * @returns {ReactElement} A React element representing the menu item.
 *
 * @example
 * <More>
 *   <MenuItem>
 *     <MoreButtonItem icon={<SettingsIcon />} onClick={() => console.log("Clicked")}>
 *       Settings
 *     </MoreButtonItem>
 *   </MenuItem>
 * </More>
 */
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
