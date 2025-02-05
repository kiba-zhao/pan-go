import type { CSSProperties, ReactNode } from "react";
import { Fragment } from "react";

import { useTranslate } from "react-admin";

import SettingsIcon from "@mui/icons-material/Settings";
import Avatar from "@mui/material/Avatar";
import Box from "@mui/material/Box";
import LinearProgress from "@mui/material/LinearProgress";
import IconButton from "@mui/material/IconButton";
import ListItem from "@mui/material/ListItem";
import ListItemAvatar from "@mui/material/ListItemAvatar";
import ListItemButton from "@mui/material/ListItemButton";
import ListItemText from "@mui/material/ListItemText";
import Stack from "@mui/material/Stack";
import Typography from "@mui/material/Typography";
import Tooltip from "@mui/material/Tooltip";
import OpenInNewIcon from "@mui/icons-material/OpenInNew";

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
        <LinearProgress />
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
  title: React.ReactNode;
  to: To;
  children?: ReactNode;
  disabled?: boolean;
  target?: React.HTMLAttributeAnchorTarget;
};
const ExtFSItemLink = ({
  title,
  to,
  children,
  disabled,
  target,
}: ExtFSItemLinkProps) => {
  return (
    <Tooltip title={title}>
      <IconButton
        component={RouterLink}
        to={to}
        target={target}
        disabled={disabled}
        onClick={(e) => e.stopPropagation()}
      >
        {children}
      </IconButton>
    </Tooltip>
  );
};

export type ExtFSItemSettingsProps = Pick<
  ExtFSItemLinkProps,
  "to" | "disabled"
>;
export const ExtFSItemSettings = ({ to, disabled }: ExtFSItemSettingsProps) => {
  const t = useTranslate();
  return (
    <ExtFSItemLink
      title={t("custom.button.settings")}
      to={to}
      disabled={disabled}
    >
      <SettingsIcon />
    </ExtFSItemLink>
  );
};

export type ExtFSItemOpenProps = Pick<ExtFSItemLinkProps, "to" | "disabled">;
export const ExtFSItemOpen = ({ to, disabled }: ExtFSItemOpenProps) => {
  const t = useTranslate();
  return (
    <ExtFSItemLink
      title={t("custom.button.open")}
      to={to}
      disabled={disabled}
      target="_blank"
    >
      <OpenInNewIcon />
    </ExtFSItemLink>
  );
};
