import type { CSSProperties, ReactNode } from "react";
import { Fragment } from "react";

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
import ListItemText from "@mui/material/ListItemText";
import Stack from "@mui/material/Stack";
import Typography from "@mui/material/Typography";

import type { To } from "react-router-dom";
import { Link as RouterLink } from "react-router-dom";

import type { ListItemData, ListItemsProps } from "../List/Item";
import { ListItems, useListItems } from "../List/Item";

export type ExtFSItemRecord<T extends any> = ListItemData<T>;

export const useExtFSItem = useListItems;

export type ExtFSItemsProps<T extends any> = {
  items: ListItemsProps<T>["items"];
  isFetching: ListItemsProps<T>["isFetching"];
  children: ListItemsProps<T>["children"];
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
      <ListItems items={items} isFetching={isFetching} itemSize={68}>
        {children}
      </ListItems>
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
    <ListItem style={style} secondaryAction={children} disablePadding>
      <ListItemButton onClick={() => onClick && onClick()} disabled={disabled}>
        <ListItemAvatar>
          <Avatar variant="rounded" sx={{ bgcolor: "inherit" }}>
            {avatarIcon}
          </Avatar>
        </ListItemAvatar>
        <ListItemText
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
