/**
 * ExtFS Page Definition File
 */

import ExpandLessIcon from "@mui/icons-material/ExpandLess";
import ExpandMoreIcon from "@mui/icons-material/ExpandMore";
import RefreshIcon from "@mui/icons-material/Refresh";
import SearchIcon from "@mui/icons-material/Search";

import Box from "@mui/material/Box";
import Breadcrumbs from "@mui/material/Breadcrumbs";
import ButtonBase from "@mui/material/ButtonBase";
import IconButton from "@mui/material/IconButton";
import Link from "@mui/material/Link";
import Menu from "@mui/material/Menu";
import MenuItem from "@mui/material/MenuItem";
import Stack from "@mui/material/Stack";
import { styled, useTheme } from "@mui/material/styles";
import Typography from "@mui/material/Typography";
import useMediaQuery from "@mui/material/useMediaQuery";

import { PageLayout } from "../Master/Page";
import { PageHeader } from "../Master/Header";
import { PageProvider } from "../Master/Context";
import { useTranslation } from "../i18n/Context";
import { RouteMore } from "../Route/More";
import {
  ExtFSSearchFileMode,
  SearchFiles,
  SearchNavigationBreadcrumbRoot,
  SearchNavigationMenuRoot,
  SearchFileRefresh,
} from "./SearchFile";
import { ExtFSHomeState, HomeItems } from "./Home";
import { ExtFSNodeFileMode, NodeFiles } from "./NodeFile";
import { ExtFSNodeMode, NodeItems } from "./NodeItem";
import { ExtFSRemoteMode, RemoteItems } from "./RemoteItem";
import { ExtFSRemoteFileMode, RemoteFiles } from "./RemoteFile";
import { SearchItems } from "./Search";
import type { ExtFSParentItem, ExtFSState } from "./State";
import { ExtFSProvider as ExtFSStateProvider, useExtFS } from "./State";

import { Dialog } from "../Common/Dialog";

import type { MouseEvent, ReactNode } from "react";
import { Fragment, useMemo, useState } from "react";
import { useIsFetching, useQueryClient } from "@tanstack/react-query";

import { ExtFSIcon, ExtFSI18nKey, ExtFSPath } from "./Route";

/**
 * Provider for ExtFS state. It initializes the state with the home state and an empty list of parent items.
 * @param children The children components to render.
 * @returns The children components wrapped in the ExtFS state provider.
 */
export const ExtFSProvider = ({ children }: { children?: ReactNode }) => {
  const initialState = {
    ...ExtFSHomeState,
    parentItems: [],
  } as ExtFSState;
  return (
    <ExtFSStateProvider value={initialState}>{children}</ExtFSStateProvider>
  );
};

const HomeTools = () => {
  const [extfs, _] = useExtFS();
  return (
    <Fragment>
      <Search />
      {extfs?.mode === ExtFSSearchFileMode ? (
        <SearchFileRefresh />
      ) : (
        <Refresh />
      )}
    </Fragment>
  );
};

const HomeAddons = () => {
  return <RouteMore path={ExtFSPath}></RouteMore>;
};

const HomeTitle = () => {
  const { t } = useTranslation();

  return <Typography variant="h6">{t(ExtFSI18nKey)}</Typography>;
};

/**
 * ExtFS Home Page component
 *
 * The root component of the ExtFS. It renders the top bar with the navigation,
 * the title, the search bar and the refresh button. It also renders the list of
 * items and the more button.
 *
 * This component is always rendered, and its children components are determined
 * by the state of the ExtFS.
 *
 * @returns The root component of the ExtFS.
 */
const Home = () => {
  const [extfs, _] = useExtFS();
  const { mode } = extfs || {};

  const ItemsElement = useMemo(() => {
    if (mode === ExtFSSearchFileMode) return <SearchFiles />;
    if (mode === ExtFSHomeState.mode) return <HomeItems />;
    if (mode === ExtFSNodeMode) return <NodeItems />;
    if (mode === ExtFSNodeFileMode) return <NodeFiles />;
    if (mode === ExtFSRemoteMode) return <RemoteItems />;
    if (mode === ExtFSRemoteFileMode) return <RemoteFiles />;
    return null;
  }, [mode]);

  return (
    <PageLayout>
      <PageHeader
        title={<HomeTitle />}
        tools={<HomeTools />}
        addons={<HomeAddons />}
      />
      <NavigationBar />
      <Box height={"calc(100% - 40px)"}> {ItemsElement}</Box>
    </PageLayout>
  );
};
/**
 * ExtFS Home Page
 *
 * This component renders the ExtFS home page.
 * It wraps the Home component with the ExtFSProvider
 * to provide the ExtFS state to the Home component.
 */
const ExtFSHome = () => (
  <Fragment>
    <PageProvider Component={ExtFSProvider} />
    <Home />
  </Fragment>
);
export default ExtFSHome;

const CustomButton = styled(ButtonBase)(({ theme }) => ({
  backgroundColor: theme.palette.background.paper,
  color: theme.palette.text.primary,
}));

/**
 * Search component
 *
 * This component renders a search button and a search input field.
 * It will open a dialog with the search input field when the button is clicked.
 * The dialog will be full screen on small screens and will be closed when
 * the escape key is pressed.
 */
const Search = () => {
  const { t } = useTranslation();
  const theme = useTheme();
  const fullScreen = useMediaQuery(theme.breakpoints.down("sm"));

  const [open, setOpen] = useState(false);
  const onOpen = () => setOpen(true);
  const onClose = () => setOpen(false);

  return (
    <Fragment>
      <CustomButton
        focusRipple
        sx={{
          borderRadius: 5,
          paddingY: 1,
          paddingX: 2,
          width: fullScreen ? "100%" : "auto",
          display: fullScreen ? "none" : "inline",
        }}
        onClick={onOpen}
      >
        <Stack
          direction="row"
          spacing={1}
          width="100%"
          alignItems="center"
          justifyContent="flex-start"
        >
          <SearchIcon color="inherit" opacity={0.5} />
          <Typography sx={{ opacity: 0.5 }}>
            {t("custom.placeholder.search")}
          </Typography>
        </Stack>
      </CustomButton>
      <IconButton
        sx={fullScreen ? void 0 : { display: "none" }}
        onClick={onOpen}
        color="inherit"
      >
        <SearchIcon />
      </IconButton>
      {open ? (
        <Dialog
          fullScreen={fullScreen}
          open={open}
          onClose={onClose}
          slotProps={{
            paper: {
              sx: fullScreen ? {} : { width: "100%", overflowX: "hidden" },
            },
          }}
        >
          <SearchItems onEsc={onClose} enabled={open} />
        </Dialog>
      ) : (
        void 0
      )}
    </Fragment>
  );
};

/**
 * Refresh Component
 *
 * This component renders a refresh button that, when clicked, triggers
 * the refetching of queries associated with the ExtFS state. It uses
 * `useIsFetching` to determine if any queries are currently fetching
 * and prevents additional refetching if so. The queries to be refetched
 * are determined by the `queryKeyList` from the ExtFS state.
 *
 * @returns A refresh icon button that refetches queries on click.
 */

const Refresh = () => {
  const [extFS, _] = useExtFS();
  const { queryKeyList } = extFS || { queryKeyList: [] };

  const isFetching = useIsFetching({
    predicate: (query) => queryKeyList.includes(query.queryKey as string[]),
  });

  const queryClient = useQueryClient();
  const onClick = async () => {
    if (isFetching) return;
    await Promise.all(
      queryKeyList.map((queryKey) =>
        queryClient.refetchQueries({ queryKey, type: "active" })
      )
    );
  };
  return (
    <IconButton onClick={onClick} color="inherit">
      <RefreshIcon />
    </IconButton>
  );
};

/**
 * A component that renders a '...' link that, when clicked, opens a menu
 * with the given `items`. When an item is clicked, the `onParentClick`
 * callback is called with the clicked item and all ancestors as an array.
 *
 * @param items An array of items to be shown in the menu.
 * @param onParentClick A callback that is called when an item is clicked.
 *                      It receives an array of items, with the first element
 *                      being the root, and the last element being the item
 *                      that was clicked.
 *
 * @returns A component that renders a '...' link that opens a menu with the
 *          given items.
 */
const NavigationMoreItems = ({
  items,
  onParentClick,
}: {
  items: ExtFSParentItem[];
  onParentClick: (items: ExtFSParentItem[]) => void;
}) => {
  const [anchorEl, setAnchorEl] = useState<null | HTMLElement>(null);
  const open = Boolean(anchorEl);
  const handleClick = (event: MouseEvent<HTMLAnchorElement>) => {
    setAnchorEl(event.currentTarget);
  };

  const handleClose = () => {
    setAnchorEl(null);
  };

  const handleMenuItemClick = (items: ExtFSParentItem[]) => {
    handleClose();
    onParentClick(items);
  };

  return (
    <Fragment>
      <Link underline="hover" color="inherit" onClick={handleClick}>
        ...
      </Link>
      <Menu anchorEl={anchorEl} open={open} onClose={() => handleClose()}>
        {items.map((item, offset) => (
          <MenuItem
            key={`nav-more-menu-${offset}`}
            onClick={() => handleMenuItemClick(items.slice(0, offset + 1))}
          >
            {item.name}
          </MenuItem>
        ))}
      </Menu>
    </Fragment>
  );
};

/**
 * A component that renders a navigation bar for ExtFS.
 *
 * @param anchorElWidth The width of the anchor element.
 * @param anchorEl The anchor element.
 *
 * @returns A component that renders a navigation bar for ExtFS.
 */
const NavigationBar = () => {
  const { t } = useTranslation();
  const [anchorEl, setAnchorEl] = useState<null | HTMLElement>(null);
  const anchorElWidth = useMemo(() => [anchorEl?.clientWidth], [anchorEl]);

  const [expand, setExpand] = useState(false);
  const handleExpand = (event: MouseEvent<HTMLButtonElement>) => {
    setAnchorEl(expand ? null : event.currentTarget);
    setExpand(!expand);
  };

  const handleClose = () => {
    setExpand(false);
  };

  const [extFS, setExtFS] = useExtFS();
  if (!extFS) return null;

  const onParentClick = (parentItems: ExtFSParentItem[]) => {
    const parentItem = parentItems.at(-1) as ExtFSParentItem;
    setExtFS({ ...parentItem.state, parentItems });
    handleClose();
  };

  const parentItems = extFS.parentItems;

  const handleHomeClick = () => {
    setExtFS({
      ...ExtFSHomeState,
      parentItems: [],
    });
    handleClose();
  };

  const MenuRoot =
    parentItems.length > 0 &&
    parentItems[0].state.mode === ExtFSSearchFileMode ? (
      <SearchNavigationMenuRoot sx={{ width: anchorElWidth }} />
    ) : (
      <MenuItem onClick={handleHomeClick} sx={{ width: anchorElWidth }}>
        {t("custom.extfs.root")}
      </MenuItem>
    );

  const BreadcrumbsRoot =
    parentItems.length > 0 &&
    parentItems[0].state.mode === ExtFSSearchFileMode ? (
      <SearchNavigationBreadcrumbRoot />
    ) : (
      <Link
        underline="hover"
        sx={{
          display: "flex",
          alignItems: "center",
          cursor: "pointer",
        }}
        color="inherit"
        onClick={handleHomeClick}
      >
        <ExtFSIcon sx={{ mr: 0.5 }} fontSize="inherit" />
      </Link>
    );

  return (
    <Box height={40}>
      <Typography
        sx={{
          display:
            parentItems.length > 0 ? "none" : { xs: "block", sm: "none" },
        }}
      >
        {t("custom.extfs.root")}
      </Typography>
      <ButtonBase
        sx={{
          display:
            parentItems.length > 0 ? { xs: "block", sm: "none" } : "none",
        }}
        onClick={handleExpand}
      >
        <Stack
          direction="row"
          spacing={1}
          width="100%"
          alignItems="center"
          justifyContent="flex-start"
          paddingRight={3}
        >
          <ExpandLessIcon
            sx={expand ? { display: "none" } : void 0}
            fontSize="small"
          />
          <ExpandMoreIcon
            sx={!expand ? { display: "none" } : void 0}
            fontSize="small"
          />
          <Typography variant="subtitle1">
            {parentItems.at(-1)?.name}
          </Typography>
        </Stack>
      </ButtonBase>
      <Menu
        anchorEl={expand ? anchorEl : null}
        open={expand}
        onClose={handleClose}
      >
        {parentItems
          .slice(0, -1)
          .map((parentItem, parentOffset) => (
            <MenuItem
              key={`nav-menu-${parentOffset}`}
              onClick={() =>
                onParentClick(parentItems.slice(0, parentOffset + 1))
              }
              sx={{ width: anchorElWidth }}
            >
              {parentItem.name}
            </MenuItem>
          ))
          .reverse()}
        {MenuRoot}
      </Menu>
      <Breadcrumbs
        aria-label="breadcrumb"
        sx={{ display: { xs: "none", sm: "flex", height: "100%" } }}
      >
        {BreadcrumbsRoot}
        {parentItems.length > 1 ? (
          <Link
            underline="hover"
            color="inherit"
            onClick={() => onParentClick(parentItems.slice(0, 1))}
          >
            {parentItems.at(0)?.name}
          </Link>
        ) : (
          void 0
        )}
        {parentItems.length > 2 && parentItems.length < 4 ? (
          <Link
            underline="hover"
            color="inherit"
            onClick={() => onParentClick(parentItems.slice(0, 2))}
          >
            {parentItems.at(1)?.name}
          </Link>
        ) : (
          void 0
        )}
        {parentItems.length > 3 ? (
          <NavigationMoreItems
            items={parentItems.slice(1, -1)}
            onParentClick={(items) => onParentClick([parentItems[0], ...items])}
          />
        ) : (
          void 0
        )}
        {parentItems.length > 0 ? (
          <Typography
            sx={{
              color: "text.primary",
            }}
          >
            {parentItems.at(-1)?.name}
          </Typography>
        ) : (
          void 0
        )}
      </Breadcrumbs>
    </Box>
  );
};
