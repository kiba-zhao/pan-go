import {
  AppBar as AdminAppBar,
  Layout,
  LayoutProps,
  LocalesMenuButton,
  RefreshIconButton,
  TitlePortal,
  ToggleThemeButton,
  useTranslate,
} from "react-admin";

import AdminPanelSettingsIcon from "@mui/icons-material/AdminPanelSettings";
import SettingsIcon from "@mui/icons-material/Settings";
import Box from "@mui/material/Box";
import Drawer from "@mui/material/Drawer";
import IconButton from "@mui/material/IconButton";
import List from "@mui/material/List";
import ListItem from "@mui/material/ListItem";
import ListItemButton from "@mui/material/ListItemButton";
import ListItemIcon from "@mui/material/ListItemIcon";
import ListItemText from "@mui/material/ListItemText";
import Tooltip from "@mui/material/Tooltip";

import { useMemo, useState } from "react";
import { useLocation, useNavigate } from "react-router-dom";

const SettingsButton = () => {
  const t = useTranslate();
  const navigate = useNavigate();
  const location = useLocation();

  const [anchorEl, setAnchorEl] = useState<null | HTMLElement>(null);
  const open = useMemo(() => Boolean(anchorEl), [anchorEl]);
  const handleClick = (event: React.MouseEvent<HTMLButtonElement>) => {
    setAnchorEl(event.currentTarget);
  };
  const handleClose = () => {
    setAnchorEl(null);
  };

  const handleCloseAndNavigate = (to: string) => {
    handleClose();
    navigate(to);
  };

  return (
    <Box component="span">
      <Tooltip title={t("custom.settings.name")}>
        <IconButton
          aria-controls="settings-menu"
          aria-haspopup="true"
          onClick={handleClick}
        >
          <SettingsIcon />
        </IconButton>
      </Tooltip>
      <Drawer anchor="right" open={open} onClose={handleClose}>
        <List sx={{ width: 250 }}>
          <ListItem disablePadding>
            <ListItemButton
              selected={location.pathname === "/app/settings"}
              onClick={() => handleCloseAndNavigate("/app/settings")}
            >
              <ListItemIcon>
                <AdminPanelSettingsIcon />
              </ListItemIcon>
              <ListItemText primary={t("custom.app/settings.name")} />
            </ListItemButton>
          </ListItem>
        </List>
      </Drawer>
    </Box>
  );
};

const AppBar = () => (
  <AdminAppBar
    toolbar={
      <>
        <TitlePortal />

        <LocalesMenuButton />
        <ToggleThemeButton />
        <SettingsButton />
        <RefreshIconButton />
      </>
    }
  ></AdminAppBar>
);

export const AppLayout = (props: LayoutProps) => (
  <Layout {...props} appBar={AppBar} />
);
