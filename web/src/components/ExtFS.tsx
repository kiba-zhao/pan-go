import { Title, useTranslate } from "react-admin";

import ExpandLessIcon from "@mui/icons-material/ExpandLess";
import ExpandMoreIcon from "@mui/icons-material/ExpandMore";
import RefreshIcon from "@mui/icons-material/Refresh";
import SearchIcon from "@mui/icons-material/Search";
import StorageIcon from "@mui/icons-material/Storage";
import Box from "@mui/material/Box";
import Breadcrumbs from "@mui/material/Breadcrumbs";
import ButtonBase from "@mui/material/ButtonBase";
import IconButton from "@mui/material/IconButton";
import Link from "@mui/material/Link";
import type { MenuProps } from "@mui/material/Menu";
import Menu from "@mui/material/Menu";
import MenuItem from "@mui/material/MenuItem";
import Paper from "@mui/material/Paper";
import Stack from "@mui/material/Stack";
import { styled, useTheme } from "@mui/material/styles";
import Typography from "@mui/material/Typography";
import useMediaQuery from "@mui/material/useMediaQuery";

import {
  ExtFSSearchFileMode,
  SearchFileMore,
  SearchFiles,
  SearchNavigationBreadcrumbRoot,
  SearchNavigationMenuRoot,
} from "./ExtFS/SearchFile";
import { ExtFSHomeState, HomeItems, HomeMore } from "./ExtFS/Home";
import { ExtFSNodeFileMode, NodeFileMore, NodeFiles } from "./ExtFS/NodeFile";
import { ExtFSNodeMode, NodeItems, NodeMore } from "./ExtFS/NodeItem";
import { ExtFSRemoteMode, RemoteItems, RemoteMore } from "./ExtFS/RemoteItem";
import {
  ExtFSRemoteFileMode,
  RemoteFileMore,
  RemoteFiles,
} from "./ExtFS/RemoteFile";
import { SearchItems } from "./ExtFS/Search";
import type { ExtFSParentItem, ExtFSState } from "./ExtFS/State";
import { ExtFSProvider as ExtFSStateProvider, useExtFS } from "./ExtFS/State";

import { Dialog } from "./Feedback/Dialog";

import type { MouseEvent, ReactNode } from "react";
import {
  Children,
  Fragment,
  useEffect,
  useMemo,
  useRef,
  useState,
} from "react";

import { useIsFetching, useQueryClient } from "@tanstack/react-query";

export const ExtFSIcon = StorageIcon;
export const ExtFSRoutePath = "/extfs";

export const ExtFSProvider = ({ children }: { children: ReactNode }) => {
  const initialState = {
    ...ExtFSHomeState,
    parentItems: [],
  } as ExtFSState;
  return (
    <ExtFSStateProvider value={initialState}>{children}</ExtFSStateProvider>
  );
};

const Home = () => {
  const t = useTranslate();

  const [extfs, _] = useExtFS();

  const [ItemsElement, MoreElement] = useMemo(() => {
    if (extfs.mode === ExtFSSearchFileMode)
      return [<SearchFiles />, <SearchFileMore />];
    if (extfs.mode === ExtFSHomeState.mode)
      return [<HomeItems />, <HomeMore />];
    if (extfs.mode === ExtFSNodeMode) return [<NodeItems />, <NodeMore />];
    if (extfs.mode === ExtFSNodeFileMode)
      return [<NodeFiles />, <NodeFileMore />];
    if (extfs.mode === ExtFSRemoteMode)
      return [<RemoteItems />, <RemoteMore />];
    if (extfs.mode === ExtFSRemoteFileMode)
      return [<RemoteFiles />, <RemoteFileMore />];
    return [];
  }, [extfs.mode]);
  return (
    <Paper
      sx={{
        height: "100%",
        paddingBottom: "76px",
        paddingTop: "56px",
        marginTop: "-56px",
      }}
    >
      <Title title={t("custom.extfs.name")} />
      <TopBar>
        <Search />
        <Refresh />
        {MoreElement}
      </TopBar>
      {ItemsElement}
    </Paper>
  );
};
export default Home;

const CustomButton = styled(ButtonBase)(({ theme }) => ({
  backgroundColor: theme.palette.background.paper,
  "&:hover, &.Mui-focusVisible": {
    backgroundColor: theme.palette.action.hover,
    opacticy: theme.palette.action.hoverOpacity,
    border: 1,
  },
}));

const Search = () => {
  const t = useTranslate();
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
          <SearchIcon />
          <Typography sx={{ opacity: 0.5 }}>
            {t("custom.placeholder.search")}
          </Typography>
        </Stack>
      </CustomButton>
      <IconButton
        sx={fullScreen ? void 0 : { display: "none" }}
        onClick={onOpen}
      >
        <SearchIcon />
      </IconButton>
      <Dialog
        fullScreen={fullScreen}
        open={open}
        onClose={onClose}
        PaperProps={{
          sx: fullScreen ? {} : { width: "100%", overflowX: "hidden" },
        }}
      >
        <SearchItems onEsc={onClose} enabled={open} />
      </Dialog>
    </Fragment>
  );
};

const Refresh = () => {
  const [extFS, _] = useExtFS();
  const { queryKeyList } = extFS;

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
    <IconButton onClick={onClick}>
      <RefreshIcon />
    </IconButton>
  );
};

const CustomNavButton = styled(ButtonBase)(({ theme }) => ({
  "&:hover, &.Mui-focusVisible": {
    backgroundColor: theme.palette.action.hover,
    opacticy: theme.palette.action.hoverOpacity,
    border: 1,
  },
}));

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

const NavigationBar = ({
  anchorElWidth,
  anchorEl,
}: {
  anchorElWidth?: number;
  anchorEl: MenuProps["anchorEl"];
}) => {
  // const t = useTranslate();
  const [expand, setExpand] = useState(false);
  const handleExpand = () => {
    setExpand(!expand);
  };

  const handleClose = () => {
    setExpand(false);
  };

  const [extFS, setExtFS] = useExtFS();

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
        ExtFS
      </MenuItem>
    );

  const BreadcrumbsRoot =
    parentItems.length > 0 &&
    parentItems[0].state.mode === ExtFSSearchFileMode ? (
      <SearchNavigationBreadcrumbRoot />
    ) : (
      <Link
        underline="hover"
        sx={{ display: "flex", alignItems: "center", cursor: "pointer" }}
        color="inherit"
        onClick={handleHomeClick}
      >
        <ExtFSIcon sx={{ mr: 0.5 }} fontSize="inherit" />
      </Link>
    );

  return (
    <Fragment>
      <Typography
        variant="h6"
        sx={{
          display:
            parentItems.length > 0 ? "none" : { xs: "block", sm: "none" },
        }}
      >
        ExtFS
      </Typography>
      <CustomNavButton
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
          <Typography variant="h6">{parentItems.at(-1)?.name}</Typography>
        </Stack>
      </CustomNavButton>
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
        sx={{ display: { xs: "none", sm: "flex" } }}
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
        {parentItems.length > 2 ? (
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
    </Fragment>
  );
};

const TopBar = ({ children }: { children?: ReactNode }) => {
  const ref = useRef<HTMLElement | null>(null);
  const [anchorEl, setAnchorEl] = useState<null | HTMLElement>(null);
  useEffect(() => {
    setAnchorEl(ref.current);
  });

  return (
    <Stack
      padding={1}
      spacing={1}
      direction="row"
      alignItems="center"
      justifyContent="space-between"
      useFlexGap
      flexWrap="wrap"
      component={Box}
      ref={ref}
    >
      <NavigationBar
        anchorElWidth={anchorEl?.clientWidth}
        anchorEl={anchorEl}
      />
      {Children.count(children) > 0 && (
        <Stack
          direction="row"
          spacing={0.5}
          alignItems="center"
          justifyContent="flex-end"
        >
          {children}
        </Stack>
      )}
    </Stack>
  );
};
