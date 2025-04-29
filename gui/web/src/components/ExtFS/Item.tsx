/**
 * ExtFS Item Component Definition File
 */
import type { CSSProperties, ReactNode } from "react";
import { Fragment, HTMLAttributeAnchorTarget } from "react";

import OpenInNewIcon from "@mui/icons-material/OpenInNew";
import SettingsIcon from "@mui/icons-material/Settings";
import Avatar from "@mui/material/Avatar";
import Box from "@mui/material/Box";
import IconButton from "@mui/material/IconButton";
import LinearProgress from "@mui/material/LinearProgress";
import ListItem from "@mui/material/ListItem";
import ListItemAvatar from "@mui/material/ListItemAvatar";
import ListItemButton from "@mui/material/ListItemButton";
import ListItemText from "@mui/material/ListItemText";
import Stack from "@mui/material/Stack";
import Tooltip from "@mui/material/Tooltip";
import Typography from "@mui/material/Typography";
import Alert from "@mui/material/Alert";

import type { To } from "../Route/Router";
import { Link as RouterLink } from "../Route/Router";

import type { ListItemData, ListItemsProps } from "../Common/Item";
import { ListItems, useListItems } from "../Common/Item";
import { useTranslation } from "../i18n/Context";

export type ExtFSItemRecord<T extends any> = ListItemData<T>;

export const useExtFSItem = useListItems;

export type ExtFSItemsProps<T extends any> = {
  items: ListItemsProps<T>["items"];
  isFetching: ListItemsProps<T>["isFetching"];
  children: ListItemsProps<T>["children"];
  error?: Error;
};
/**
 * ExtFSItems Component
 *
 * This component renders a list of items with a loading indicator.
 *
 * @template T - The type of the items in the list.
 *
 * @param {ExtFSItemsProps<T>} props - The props for the component.
 * @param {T[]} props.items - The list of items to be displayed.
 * @param {boolean} props.isFetching - Flag indicating if the data is currently being fetched.
 * @param {ReactNode} props.children - The children elements to be rendered within the ListItems.
 *
 * @returns {JSX.Element} A JSX element that contains a loading indicator when fetching
 *                        and a list of items when data is available.
 */

export const ExtFSItems = <T extends any>({
  items,
  isFetching,
  children,
  error,
}: ExtFSItemsProps<T>) => {
  const { t } = useTranslation();

  return (
    <Fragment>
      {isFetching && (
        <Box
          height="100%"
          sx={{
            alignItems: "center",
            justifyContent: "center",
          }}
        >
          <LinearProgress />
        </Box>
      )}
      {!isFetching && error !== void 0 ? (
        <Alert severity="error">
          {t(`errors.${error.name}`, { _: error.message })}
        </Alert>
      ) : (
        void 0
      )}
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
/**
 * ExtFSItem Component
 *
 * This component renders a single item of a list with an optional avatar and
 * extended icon.
 *
 * @param {ExtFSItemProps} props - The props for the component.
 * @param {CSSProperties} props.style - The style for the list item.
 * @param {() => void} props.onClick - The function to be called when the list item is clicked.
 * @param {string} props.primary - The primary text of the list item.
 * @param {string} props.secondary - The secondary text of the list item.
 * @param {ReactNode} props.avatarIcon - The icon to be displayed as the avatar of the list item.
 * @param {ReactNode} props.extIcon - The extended icon to be displayed next to the primary text.
 * @param {boolean} props.disabled - Flag indicating if the list item is disabled.
 * @param {ReactNode} props.children - The children elements to be rendered within the ListItemButton.
 *
 * @returns {JSX.Element} A JSX element representing a list item with an avatar and optional extended icon.
 */
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
          slotProps={{ secondary: { component: "div" } }}
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
  target?: HTMLAttributeAnchorTarget;
};
/**
 * ExtFSItemLink Component
 *
 * This component renders a link button with a tooltip.
 * It is used within the `ExtFSItem` component to provide a link to a settings page.
 *
 * @param {ExtFSItemLinkProps} props - The props for the component.
 * @param {ReactNode} props.title - The text to be displayed in the tooltip.
 * @param {To} props.to - The target URL of the link.
 * @param {ReactNode} props.children - The children elements to be rendered within the IconButton.
 * @param {boolean} props.disabled - Flag indicating if the link button is disabled.
 * @param {React.HTMLAttributeAnchorTarget} props.target - The target attribute for the link.
 *
 * @returns {JSX.Element} A JSX element representing a link button with a tooltip.
 */
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
> &
  Partial<Pick<ExtFSItemLinkProps, "title">>;
/**
 * ExtFSItemSettings Component
 *
 * This component renders a link button with a tooltip.
 * It is used within the `ExtFSItem` component to provide a link to a settings page.
 *
 * @param {ExtFSItemSettingsProps} props - The props for the component.
 * @param {To} props.to - The target URL of the link.
 * @param {boolean} props.disabled - Flag indicating if the link button is disabled.
 *
 * @returns {JSX.Element} A JSX element representing a link button with a tooltip.
 */
export const ExtFSItemSettings = ({
  to,
  disabled,
  title,
}: ExtFSItemSettingsProps) => {
  const { t } = useTranslation();
  return (
    <ExtFSItemLink
      title={title ?? t("custom.button.settings")}
      to={to}
      disabled={disabled}
    >
      <SettingsIcon />
    </ExtFSItemLink>
  );
};

export type ExtFSItemOpenProps = Pick<ExtFSItemLinkProps, "to" | "disabled"> & {
  hidden?: boolean;
} & Partial<Pick<ExtFSItemLinkProps, "title">>;
/**
 * ExtFSItemOpen Component
 *
 * This component renders a link button with a tooltip.
 * It is used within the `ExtFSItem` component to provide a link to open a file or directory
 * in a new tab.
 *
 * @param {ExtFSItemOpenProps} props - The props for the component.
 * @param {To} props.to - The target URL of the link.
 * @param {boolean} props.disabled - Flag indicating if the link button is disabled.
 * @param {boolean} [props.hidden] - Flag indicating if the component should not be rendered.
 *
 * @returns {JSX.Element} A JSX element representing a link button with a tooltip.
 */
export const ExtFSItemOpen = ({
  to,
  disabled,
  hidden,
  title,
}: ExtFSItemOpenProps) => {
  const { t } = useTranslation();
  if (hidden) return;
  return (
    <ExtFSItemLink
      title={title ?? t("custom.button.open")}
      to={to}
      disabled={disabled}
      target="_blank"
    >
      <OpenInNewIcon />
    </ExtFSItemLink>
  );
};
