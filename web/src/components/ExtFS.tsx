import { Title, useTranslate } from "react-admin";

import AddCircleIcon from "@mui/icons-material/AddCircle";
import BookmarkBorderIcon from "@mui/icons-material/BookmarkBorder";
import ClearIcon from "@mui/icons-material/Clear";
import CloudIcon from "@mui/icons-material/Cloud";
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
import Dialog from "@mui/material/Dialog";
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
import Menu from "@mui/material/Menu";
import MenuItem from "@mui/material/MenuItem";
import Paper from "@mui/material/Paper";
import Stack from "@mui/material/Stack";

import ExpandLessIcon from "@mui/icons-material/ExpandLess";
import { styled, useTheme } from "@mui/material/styles";
import Typography from "@mui/material/Typography";
import useMediaQuery from "@mui/material/useMediaQuery";

import type { MenuProps } from "@mui/material/Menu";

import type { MouseEvent } from "react";
import { Fragment, useEffect, useMemo, useRef, useState } from "react";
import { AppNodeIcon } from "./AppNodes";

import { Link as RouterLink } from "react-router-dom";

export const ExfFSIcon = StorageIcon;

const Home = () => {
  const t = useTranslate();
  return (
    <Paper>
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
      <IconButton
        aria-label="more"
        id="long-button"
        aria-controls={open ? "long-menu" : undefined}
        aria-expanded={open ? "true" : undefined}
        aria-haspopup="true"
        onClick={handleClick}
      >
        <MoreVertIcon />
      </IconButton>
      <Menu
        id="long-menu"
        MenuListProps={{
          "aria-labelledby": "long-button",
        }}
        anchorEl={anchorEl}
        open={open}
        onClose={handleClose}
      >
        <MenuItem onClick={handleClose}>
          <ListItemIcon>
            <OpenInNewIcon />
          </ListItemIcon>
          <ListItemText>Open</ListItemText>
        </MenuItem>
        <MenuItem onClick={handleClose}>
          <ListItemIcon>
            <AddCircleIcon fontSize="small" />
          </ListItemIcon>
          <ListItemText>New</ListItemText>
        </MenuItem>
        <MenuItem onClick={handleClose}>
          <ListItemIcon>
            <SettingsIcon fontSize="small" />
          </ListItemIcon>
          <ListItemText>Settings</ListItemText>
        </MenuItem>
        <MenuItem onClick={handleClose}>
          <ListItemIcon>
            <HelpIcon fontSize="small" />
          </ListItemIcon>
          <ListItemText>Help</ListItemText>
        </MenuItem>
      </Menu>
    </Fragment>
  );
};

const Refresh = () => {
  return (
    <IconButton>
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

const NavigationBar = ({
  anchorElWidth,
  anchorEl,
}: {
  anchorElWidth?: number;
  anchorEl: MenuProps["anchorEl"];
}) => {
  const [expand, setExpand] = useState(false);
  const handleExpand = () => {
    setExpand(!expand);
  };

  const handleClose = () => {
    setExpand(false);
  };
  return (
    <Fragment>
      <CustomNavButton
        sx={{ display: { xs: "block", sm: "none" } }}
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
          <Typography variant="h6">Belts</Typography>
        </Stack>
      </CustomNavButton>
      <Menu anchorEl={anchorEl} open={expand} onClose={handleClose}>
        <MenuItem onClick={handleClose} sx={{ width: anchorElWidth }}>
          Floder 1
        </MenuItem>
        <MenuItem onClick={handleClose} sx={{ width: anchorElWidth }}>
          Node 1
        </MenuItem>
        <MenuItem onClick={handleClose} sx={{ width: anchorElWidth }}>
          ExtFS
        </MenuItem>
      </Menu>
      <Breadcrumbs
        aria-label="breadcrumb"
        sx={{ display: { xs: "none", sm: "flex" } }}
      >
        <Link
          underline="hover"
          sx={{ display: "flex", alignItems: "center" }}
          color="inherit"
          href="#"
        >
          <ExfFSIcon sx={{ mr: 0.5 }} fontSize="inherit" />
        </Link>
        <Link underline="hover" color="inherit" href="#">
          Node 1
        </Link>
        <Link underline="hover" color="inherit" href="#">
          ...
        </Link>
        <Typography sx={{ color: "text.primary" }}>Belts</Typography>
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

type ExtFSFileItem = {
  fileType: "E" | "S" | "D" | "F" | "N" | "RD" | "RF" | "RN";
  name: string;
  updatedAt: string;
  tagQuantity: number;
  pendingTagQuantity: number;
  disabled?: boolean;
  linkId?: string;
};
type FileItemProps = {
  record: ExtFSFileItem;
};
const FileItem = ({ record }: FileItemProps) => {
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

  const settingsUrl = useMemo(() => {
    if (record.fileType === "N") return "/extfs/local-node";
    if (record.fileType === "RN") return `/extfs/remote-node/${record.linkId}`;
    if (record.fileType === "F" || record.fileType === "D")
      return `/extfs/local-file/${record.linkId}`;
    if (record.fileType === "RF" || record.fileType === "RD")
      return `/extfs/remote-file/${record.linkId}`;
    return ".";
  }, [record.fileType, record.linkId]);

  const tagUrl = useMemo(() => {
    const searchParams = new URLSearchParams({
      fileType: record.fileType,
    });

    if (
      record.fileType.at(0) === "R" ||
      record.fileType === "F" ||
      record.fileType === "D"
    ) {
      searchParams.set("linkId", record.linkId || "");
    }
    return `/extfs/tags?${searchParams.toString()}`;
  }, [record.fileType, record.linkId]);

  return (
    <ListItem>
      <ListItemButton>
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

        <ListItemSecondaryAction>
          <IconButton
            component={RouterLink}
            to={tagUrl}
            disabled={record.disabled}
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
          <IconButton component={RouterLink} to={settingsUrl}>
            <SettingsIcon />
          </IconButton>
        </ListItemSecondaryAction>
      </ListItemButton>
    </ListItem>
  );
};

const FileItems = () => {
  const t = useTranslate();

  const records: ExtFSFileItem[] = [
    {
      name: "Local Node",
      fileType: "N",
      updatedAt: "Jan 9, 2014",
      tagQuantity: 0,
      pendingTagQuantity: 0,
    },
    {
      name: "Folder",
      fileType: "D",
      updatedAt: "Jan 9, 2014",
      tagQuantity: 2,
      pendingTagQuantity: 3,
    },
    {
      name: "File",
      fileType: "RF",
      updatedAt: "Jan 9, 2014",
      tagQuantity: 0,
      pendingTagQuantity: 0,
    },
    {
      name: "Remote Node with offline",
      fileType: "RN",
      updatedAt: "Jan 1, 2014",
      tagQuantity: 0,
      pendingTagQuantity: 0,
      disabled: true,
    },
    {
      name: "Folder with Error",
      fileType: "RD",
      updatedAt: "Jan 1, 2014",
      tagQuantity: 1,
      pendingTagQuantity: 4,
      disabled: true,
    },
    {
      name: "File with Error",
      fileType: "RF",
      updatedAt: "Jan 1, 2014",
      tagQuantity: 0,
      pendingTagQuantity: 0,
      disabled: true,
    },
  ];
  return (
    <List>
      {records.map((record) => (
        <FileItem key={record.name} record={record} />
      ))}
    </List>
  );
};
