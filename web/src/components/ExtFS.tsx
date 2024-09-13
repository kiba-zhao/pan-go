import { Title, useTranslate } from "react-admin";

import AddCircleIcon from "@mui/icons-material/AddCircle";
import BookmarkBorderIcon from "@mui/icons-material/BookmarkBorder";
import ClearIcon from "@mui/icons-material/Clear";
import CloudIcon from "@mui/icons-material/Cloud";
import ExpandLessIcon from "@mui/icons-material/ExpandLess";
import ExpandMoreIcon from "@mui/icons-material/ExpandMore";
import FolderIcon from "@mui/icons-material/Folder";
import HelpIcon from "@mui/icons-material/Help";
import InsertDriveFileIcon from "@mui/icons-material/InsertDriveFile";
import MoreVertIcon from "@mui/icons-material/MoreVert";
import OpenInNewIcon from "@mui/icons-material/OpenInNew";
import RefreshIcon from "@mui/icons-material/Refresh";
import SearchIcon from "@mui/icons-material/Search";
import SettingsIcon from "@mui/icons-material/Settings";
import StorageIcon from "@mui/icons-material/Storage";
import Avatar from "@mui/material/Avatar";
import Badge from "@mui/material/Badge";
import Box from "@mui/material/Box";
import Breadcrumbs from "@mui/material/Breadcrumbs";
import ButtonBase from "@mui/material/ButtonBase";
import Chip from "@mui/material/Chip";
import CircularProgress from "@mui/material/CircularProgress";
import Divider from "@mui/material/Divider";
import IconButton from "@mui/material/IconButton";
import InputBase from "@mui/material/InputBase";
import Link from "@mui/material/Link";
import List from "@mui/material/List";
import ListItem from "@mui/material/ListItem";
import ListItemAvatar from "@mui/material/ListItemAvatar";
import ListItemButton from "@mui/material/ListItemButton";
import ListItemIcon from "@mui/material/ListItemIcon";
import ListItemSecondaryAction from "@mui/material/ListItemSecondaryAction";
import ListItemText from "@mui/material/ListItemText";
import type { MenuProps } from "@mui/material/Menu";
import Menu from "@mui/material/Menu";
import MenuItem from "@mui/material/MenuItem";
import Paper from "@mui/material/Paper";
import Stack from "@mui/material/Stack";
import { styled, useTheme } from "@mui/material/styles";
import Typography from "@mui/material/Typography";
import useMediaQuery from "@mui/material/useMediaQuery";

import type { ExtFSItem, ExtFSItemSearchCondition } from "../API";
import { useAPI } from "../API";
import { AppNodeIcon } from "./AppNodes";
import {
  createReducerContext,
  ReducerStateProvider,
  useReducerState,
} from "./Context/ReducerState";
import { ExtFSNodeItemRoutePath } from "./ExtFSNodeItem";
import { ExtFSTagRoutePath } from "./ExtFSTag";
import { Dialog } from "./Feedback/Dialog";

import type { CSSProperties, MouseEvent, ReactNode } from "react";
import { Fragment, useEffect, useMemo, useRef, useState } from "react";

import { useIsFetching, useQuery, useQueryClient } from "@tanstack/react-query";
import { Link as RouterLink } from "react-router-dom";
import AutoSizer from "react-virtualized-auto-sizer";
import { FixedSizeList } from "react-window";

export const ExtFSIcon = StorageIcon;
export const ExtFSRoutePath = "/extfs";

type ExtFSState = {
  queryKey: string[];
  condition: ExtFSItemSearchCondition;
  defaults: ExtFSItemSearchCondition;
  parentItems: ExtFSItem[];
};

const ExtFSContext = createReducerContext<ExtFSState>();

const useExtFS = () => useReducerState(ExtFSContext);

export const ExtFSProvider = ({ children }: { children: ReactNode }) => {
  const condition = { fileType: ["N", "RN"] } as ExtFSItemSearchCondition;
  const initialState = {
    queryKey: ["extfs-items"],
    condition: condition,
    defaults: condition,
    parentItems: [],
  } as ExtFSState;
  return (
    <ReducerStateProvider<ExtFSState>
      initialState={initialState}
      opts={ExtFSContext}
    >
      {children}
    </ReducerStateProvider>
  );
};

const Home = () => {
  const t = useTranslate();

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
      <TopBar />
      <FileItems />
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
          <Typography sx={{ opacity: 0.5 }}>Search ...</Typography>
        </Stack>
      </CustomButton>
      <IconButton
        sx={fullScreen ? void 0 : { display: "none" }}
        onClick={onOpen}
      >
        <SearchIcon />
      </IconButton>
      <Dialog fullScreen={fullScreen} open={open} onClose={onClose}>
        <Stack
          padding={1}
          component="form"
          spacing={0.5}
          direction="row"
          alignItems="center"
          justifyContent="space-between"
        >
          <SearchIcon />
          <InputBase
            placeholder="What are you looking for ?"
            fullWidth
            size="medium"
            autoFocus
          />
          <ButtonBase onClick={onClose}>
            <Chip
              label="esc"
              variant="outlined"
              sx={{ borderRadius: 1 }}
              size="small"
            />
          </ButtonBase>
        </Stack>
        <Divider />
        <List
          sx={{ pt: 0, ...(fullScreen ? {} : { height: 680, width: 552 }) }}
        >
          <ListItem disableGutters>
            <ListItemButton>
              <ListItemText
                primary="Keywords 1"
                sx={{ paddingRight: 5 }}
              ></ListItemText>
              <ListItemSecondaryAction>
                <IconButton size="small">
                  <ClearIcon fontSize="small" />
                </IconButton>
              </ListItemSecondaryAction>
            </ListItemButton>
          </ListItem>
          <ListItem disableGutters>
            <ListItemButton>
              <ListItemText
                primary="Keywords 2"
                sx={{ paddingRight: 5 }}
              ></ListItemText>
              <ListItemSecondaryAction>
                <IconButton size="small">
                  <ClearIcon fontSize="small" />
                </IconButton>
              </ListItemSecondaryAction>
            </ListItemButton>
          </ListItem>
          <ListItem disableGutters>
            <ListItemButton>
              <ListItemText
                primary="Keywords 3"
                sx={{ paddingRight: 5 }}
              ></ListItemText>
              <ListItemSecondaryAction>
                <IconButton size="small">
                  <ClearIcon fontSize="small" />
                </IconButton>
              </ListItemSecondaryAction>
            </ListItemButton>
          </ListItem>
        </List>
      </Dialog>
    </Fragment>
  );
};

const MoreMenu = ({ onClose, ...props }: MenuProps) => {
  const t = useTranslate();

  const [extFS, _] = useExtFS();

  const parentItem = useMemo(() => extFS.parentItems.at(-1), [extFS]);

  const newUrl = useMemo(() => {
    if (parentItem && parentItem.fileType === "N")
      return `${ExtFSNodeItemRoutePath}/create`;
  }, [parentItem]);

  const settingsUrl = useMemo(() => {
    if (parentItem && parentItem.fileType === "D") {
      if (extFS.parentItems.at(-2)?.fileType === "N")
        return `${ExtFSNodeItemRoutePath}/${parentItem.linkId}`;
    }
  }, [parentItem, parentItem?.fileType, parentItem?.linkId, extFS.parentItems]);
  const handleItemClick = (event: React.MouseEvent<HTMLElement>) => {
    onClose && onClose(event, "backdropClick");
  };

  return (
    <Menu onClose={onClose} {...props}>
      <MenuItem onClick={handleItemClick} sx={{ display: "none" }}>
        <ListItemIcon>
          <OpenInNewIcon />
        </ListItemIcon>
        <ListItemText>Open</ListItemText>
      </MenuItem>
      {newUrl !== void 0 ? (
        <MenuItem onClick={handleItemClick} component={RouterLink} to={newUrl}>
          <ListItemIcon>
            <AddCircleIcon fontSize="small" />
          </ListItemIcon>
          <ListItemText>{t("custom.button.new")}</ListItemText>
        </MenuItem>
      ) : (
        void 0
      )}
      {settingsUrl !== void 0 ? (
        <MenuItem
          onClick={handleItemClick}
          component={RouterLink}
          to={settingsUrl}
        >
          <ListItemIcon>
            <SettingsIcon fontSize="small" />
          </ListItemIcon>
          <ListItemText>{t("custom.button.settings")}</ListItemText>
        </MenuItem>
      ) : (
        void 0
      )}
      <MenuItem onClick={handleItemClick}>
        <ListItemIcon>
          <HelpIcon fontSize="small" />
        </ListItemIcon>
        <ListItemText>Help</ListItemText>
      </MenuItem>
    </Menu>
  );
};

const More = () => {
  const [anchorEl, setAnchorEl] = useState<null | HTMLElement>(null);
  const open = Boolean(anchorEl);
  const handleClick = (event: MouseEvent<HTMLButtonElement>) => {
    setAnchorEl(event.currentTarget);
  };
  const handleClose = () => {
    setAnchorEl(null);
  };

  return (
    <Fragment>
      <IconButton onClick={handleClick}>
        <MoreVertIcon />
      </IconButton>
      <MoreMenu anchorEl={anchorEl} open={open} onClose={handleClose} />
    </Fragment>
  );
};

const Refresh = () => {
  const [extFS, _] = useExtFS();
  const { queryKey, condition } = extFS;
  const queryKey_ = useMemo(
    () => [...queryKey, condition],
    [queryKey, condition]
  );
  const isFetching = useIsFetching({
    queryKey: queryKey_,
  });
  const queryClient = useQueryClient();
  const onClick = async () => {
    if (isFetching) return;
    await queryClient.refetchQueries({
      queryKey: queryKey_,
      type: "active",
    });
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
  items: ExtFSItem[];
  onParentClick: (items: ExtFSItem[]) => void;
}) => {
  const [anchorEl, setAnchorEl] = useState<null | HTMLElement>(null);
  const open = Boolean(anchorEl);
  const handleClick = (event: MouseEvent<HTMLAnchorElement>) => {
    setAnchorEl(event.currentTarget);
  };
  const handleClose = () => {
    setAnchorEl(null);
  };

  const handleMenuItemClick = (items: ExtFSItem[]) => {
    handleClose();
    onParentClick(items);
  };

  return (
    <Fragment>
      <Link underline="hover" color="inherit" onClick={handleClick}>
        ...
      </Link>
      <Menu anchorEl={anchorEl} open={open} onClose={handleClose}>
        {items.map((item, offset) => (
          <MenuItem
            key={`nav-more-menu-${item.id}`}
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
  const t = useTranslate();
  const [expand, setExpand] = useState(false);
  const handleExpand = () => {
    setExpand(!expand);
  };

  const handleClose = () => {
    setExpand(false);
  };

  const [extFS, setExtFS] = useExtFS();
  const { parentItems, defaults } = extFS;
  const onRootClick = () => {
    setExtFS({ ...extFS, condition: defaults, parentItems: [] });
    handleClose();
  };

  const onParentClick = (parentItems: ExtFSItem[]) => {
    const { id } = parentItems.at(-1) as ExtFSItem;

    setExtFS({ ...extFS, parentItems, condition: { parentId: id } });
    handleClose();
  };

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
      <Menu anchorEl={anchorEl} open={expand} onClose={handleClose}>
        {parentItems.slice(0, -1).map((parentItem, parentOffset) => (
          <MenuItem
            key={`nav-menu-${parentItem.id}`}
            onClick={() => onParentClick(parentItems.slice(0, parentOffset))}
            sx={{ width: anchorElWidth }}
          >
            {parentItem.name}
          </MenuItem>
        ))}

        <MenuItem onClick={onRootClick} sx={{ width: anchorElWidth }}>
          ExtFS
        </MenuItem>
      </Menu>
      <Breadcrumbs
        aria-label="breadcrumb"
        sx={{ display: { xs: "none", sm: "flex" } }}
      >
        <Link
          underline="hover"
          sx={{ display: "flex", alignItems: "center", cursor: "pointer" }}
          color="inherit"
          onClick={onRootClick}
        >
          <ExtFSIcon sx={{ mr: 0.5 }} fontSize="inherit" />
        </Link>
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

const TopBar = () => {
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
      <Stack
        direction="row"
        spacing={0.5}
        alignItems="center"
        justifyContent="flex-end"
      >
        <Search />
        <Refresh />
        <More />
      </Stack>
    </Stack>
  );
};

type FileItemProps = {
  record: ExtFSItem;
  style?: CSSProperties;
};
const FileItem = ({ record, style }: FileItemProps) => {
  const [extFS, setExtFS] = useExtFS();
  const { parentItems } = extFS;
  const parentItem = useMemo(() => extFS.parentItems.at(-1), [extFS]);

  const [avatarIcon, extIcon] = useMemo(() => {
    const extIcon =
      record.fileType.at(0) === "R" ? (
        <CloudIcon fontSize="small" />
      ) : undefined;

    if (record.fileType === "D" || record.fileType === "RD") {
      return [
        <FolderIcon
          color={record.disabled ? "disabled" : "primary"}
          fontSize="large"
        />,
        extIcon,
      ];
    }
    if (record.fileType === "F" || record.fileType === "RF") {
      return [
        <InsertDriveFileIcon
          color={record.disabled ? "disabled" : "action"}
          fontSize="large"
        />,
        extIcon,
      ];
    }
    if (record.fileType === "N" || record.fileType === "RN") {
      return [
        <AppNodeIcon
          color={record.disabled ? "disabled" : "primary"}
          fontSize="large"
        />,
        extIcon,
      ];
    }

    return [];
  }, [record.fileType, record.disabled]);

  const tagUrl = useMemo(() => {
    const searchParams = new URLSearchParams({
      fileType: record.fileType,
    });

    if (
      record.fileType.at(0) === "R" ||
      record.fileType === "F" ||
      record.fileType === "D"
    ) {
      searchParams.set("itemId", record.id || "");
    }
    return `${ExtFSTagRoutePath}?${searchParams.toString()}`;
  }, [record.fileType, record.id]);

  const settingsUrl = useMemo(() => {
    if (record.fileType === "F" || record.fileType === "D") {
      if (parentItem && parentItem.fileType === "N") {
        return `${ExtFSNodeItemRoutePath}/${record.linkId}`;
      }
    }
  }, [record.fileType, record.linkId, parentItem]);

  const onItemClick = () => {
    if (record.fileType === "F" || record.fileType === "RF") {
      // TODO: open file
      return;
    }
    // Comment: show sub items
    setExtFS({
      ...extFS,
      condition: { parentId: record.id },
      parentItems: [...parentItems, record],
    });
  };
  return (
    <ListItem style={style}>
      <ListItemButton onClick={onItemClick}>
        <ListItemAvatar>
          <Avatar variant="rounded" sx={{ bgcolor: "inherit" }}>
            {avatarIcon}
          </Avatar>
        </ListItemAvatar>
        <ListItemText
          sx={{ paddingRight: 10 }}
          primary={record.name}
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
              <Typography>{record.updatedAt}</Typography>
            </Stack>
          }
        />

        <ListItemSecondaryAction
          sx={record.fileType === "N" ? { display: "none" } : void 0}
        >
          <IconButton
            component={RouterLink}
            to={tagUrl}
            disabled={record.disabled}
            onClick={(e) => e.stopPropagation()}
          >
            <Badge badgeContent={record.pendingTagQuantity}>
              <Badge
                badgeContent={record.tagQuantity}
                anchorOrigin={{ vertical: "bottom", horizontal: "right" }}
              >
                <BookmarkBorderIcon />
              </Badge>
            </Badge>
          </IconButton>
          {settingsUrl ? (
            <IconButton
              component={RouterLink}
              to={settingsUrl}
              onClick={(e) => e.stopPropagation()}
            >
              <SettingsIcon />
            </IconButton>
          ) : (
            void 0
          )}
        </ListItemSecondaryAction>
      </ListItemButton>
    </ListItem>
  );
};

const FileItems = () => {
  const t = useTranslate();

  const api = useAPI();
  const [extFS, _setExtFS] = useExtFS();
  const { queryKey, condition } = extFS;
  const { data, isFetching, isError } = useQuery({
    queryKey: [...queryKey, condition],
    queryFn: async () => await api?.searchExtFSItems(condition),
  });

  const [_, items] = useMemo(() => data || [0, []], [data]);

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
              <FileItem
                key={`extfs-items-${items[index].id}`}
                style={style}
                record={items[index]}
              />
            )}
          </FixedSizeList>
        )}
      </AutoSizer>
    </Fragment>
  );
};
