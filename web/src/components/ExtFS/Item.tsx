import type { CSSProperties, ReactNode } from "react";
import {
  Children,
  createContext,
  createElement,
  Fragment,
  isValidElement,
  useContext,
} from "react";

import BookmarkBorderIcon from "@mui/icons-material/BookmarkBorder";
import SettingsIcon from "@mui/icons-material/Settings";
import Avatar from "@mui/material/Avatar";
import Badge from "@mui/material/Badge";
import Box from "@mui/material/Box";
import CircularProgress from "@mui/material/CircularProgress";
import IconButton from "@mui/material/IconButton";
import ListItem from "@mui/material/ListItem";
import ListItemAvatar from "@mui/material/ListItemAvatar";
import ListItemButton from "@mui/material/ListItemButton";
import ListItemSecondaryAction from "@mui/material/ListItemSecondaryAction";
import ListItemText from "@mui/material/ListItemText";
import Stack from "@mui/material/Stack";
import Typography from "@mui/material/Typography";

import type { To } from "react-router-dom";
import { Link as RouterLink } from "react-router-dom";
import AutoSizer from "react-virtualized-auto-sizer";
import { FixedSizeList } from "react-window";

export type ExtFSItemRecord<T extends any> = {
  data: T;
  style?: CSSProperties;
};
const ExtFSItemContext = createContext<ExtFSItemRecord<any>>(null!);
export const useExtFSItem = <T extends any>() =>
  useContext<ExtFSItemRecord<T>>(ExtFSItemContext);

export type ExtFSItemsProps<T extends any> = {
  items: T[];
  isFetching: boolean;
  children: ReactNode;
};
export const ExtFSItems = <T extends any>({
  items,
  isFetching,
  children,
}: ExtFSItemsProps<T>) => {
  return (
    <Fragment>
      <Box
        height="100%"
        sx={{
          display: isFetching ? "flex" : "none",
          alignItems: "center",
          justifyContent: "center",
        }}
      >
        <CircularProgress />
      </Box>
      <AutoSizer disableWidth={true} hidden={isFetching}>
        {({ height }) => (
          <FixedSizeList
            height={height}
            width="100%"
            itemSize={68}
            itemCount={items.length}
          >
            {({ index, style }) => (
              <ExtFSItemContext.Provider
                value={{ data: items[index], style }}
                key={`extfs-items-${index}`}
              >
                {Children.map(children, (child) =>
                  child && isValidElement(child)
                    ? createElement(child.type, child.props)
                    : child
                )}
              </ExtFSItemContext.Provider>
            )}
          </FixedSizeList>
        )}
      </AutoSizer>
    </Fragment>
  );
};

export type ExtFSItemProps = {
  style?: CSSProperties;
  onClick?: () => void;
  primary: string;
  secondary: string;
  avatarIcon: ReactNode;
  extIcon?: ReactNode;
  disabled?: boolean;
  children?: ReactNode;
};
export const ExtFSItem = ({
  style,
  onClick,
  primary,
  secondary,
  avatarIcon,
  extIcon,
  disabled,
  children,
}: ExtFSItemProps) => {
  return (
    <ListItem style={style}>
      <ListItemButton onClick={() => onClick && onClick()} disabled={disabled}>
        <ListItemAvatar>
          <Avatar variant="rounded" sx={{ bgcolor: "inherit" }}>
            {avatarIcon}
          </Avatar>
        </ListItemAvatar>
        <ListItemText
          sx={{ paddingRight: 10 }}
          primary={primary}
          secondaryTypographyProps={{ component: "div" }}
          secondary={
            <Stack
              direction="row"
              spacing={1}
              useFlexGap
              alignItems="center"
              justifyContent="flex-start"
              flexWrap="wrap"
            >
              {extIcon}
              <Typography>{secondary}</Typography>
            </Stack>
          }
        />
        {Children.count(children) ? (
          <ListItemSecondaryAction>{children}</ListItemSecondaryAction>
        ) : (
          void 0
        )}
      </ListItemButton>
    </ListItem>
  );
};

type ExtFSItemLinkProps = {
  to: To;
  children?: ReactNode;
  disabled?: boolean;
};
const ExtFSItemLink = ({ to, children, disabled }: ExtFSItemLinkProps) => {
  return (
    <IconButton
      component={RouterLink}
      to={to}
      disabled={disabled}
      onClick={(e) => e.stopPropagation()}
    >
      {children}
    </IconButton>
  );
};

export type ExtFSItemTagProps = {
  pendingQuantity?: number;
  quantity?: number;
} & Pick<ExtFSItemLinkProps, "to" | "disabled">;
export const ExtFSItemTag = ({
  to,
  disabled,
  pendingQuantity,
  quantity,
}: ExtFSItemTagProps) => {
  return (
    <ExtFSItemLink to={to} disabled={disabled}>
      <Badge badgeContent={pendingQuantity}>
        <Badge
          badgeContent={quantity}
          anchorOrigin={{ vertical: "bottom", horizontal: "right" }}
        >
          <BookmarkBorderIcon />
        </Badge>
      </Badge>
    </ExtFSItemLink>
  );
};

export type ExtFSItemSettingsProps = Pick<
  ExtFSItemLinkProps,
  "to" | "disabled"
>;
export const ExtFSItemSettings = ({ to, disabled }: ExtFSItemSettingsProps) => {
  return (
    <ExtFSItemLink to={to} disabled={disabled}>
      <SettingsIcon />
    </ExtFSItemLink>
  );
};
