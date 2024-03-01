import { useEffect, Fragment, useState, MouseEvent } from "react";

import AppBar from "@mui/material/AppBar";
import Toolbar from "@mui/material/Toolbar";
import Typography from "@mui/material/Typography";
import Box from "@mui/material/Box";
import Link from "@mui/material/Link";
import Button from "@mui/material/Button";
import Menu from "@mui/material/Menu";
import MenuItem from "@mui/material/MenuItem";
import IconButton from "@mui/material/IconButton";
import Tooltip from "@mui/material/Tooltip";
import ViewModuleIcon from "@mui/icons-material/ViewModule";
import SettingsIcon from "@mui/icons-material/Settings";
import TranslateIcon from "@mui/icons-material/Translate";
import KeyboardArrowDownIcon from "@mui/icons-material/KeyboardArrowDown";

import { useTranslation } from "react-i18next";
import { useHref } from "react-router-dom";

function LangSelect() {
  const { t, i18n } = useTranslation();
  const [anchorEl, setAnchorEl] = useState<null | HTMLElement>(null);
  const open = Boolean(anchorEl);
  const onOpen = (event: MouseEvent<HTMLElement>) => {
    setAnchorEl(event.currentTarget);
  };
  const onClose = () => {
    setAnchorEl(null);
  };

  const onLangChange = (lang: string) => {
    i18n.changeLanguage(lang);
    onClose();
  };

  return (
    <Fragment>
      <Button
        component="label"
        variant="text"
        color="inherit"
        startIcon={<TranslateIcon />}
        endIcon={<KeyboardArrowDownIcon />}
        onClick={onOpen}
      >
        {t(`languages.${i18n.language}`)}
      </Button>
      <Menu
        id="basic-menu"
        anchorEl={anchorEl}
        open={open}
        onClose={onClose}
        MenuListProps={{
          "aria-labelledby": "basic-button",
        }}
      >
        <MenuItem onClick={() => onLangChange("en")}>
          {t(`languages.en`)}
        </MenuItem>
        <MenuItem onClick={() => onLangChange("zh-CN")}>
          {t(`languages.zh-CN`)}
        </MenuItem>
      </Menu>
    </Fragment>
  );
}

function Header() {
  const { t } = useTranslation();
  useEffect(() => {});
  return (
    <AppBar component="nav" position="static">
      <Toolbar
        sx={{
          pr: "24px", // keep right padding when drawer closed
        }}
      >
        <Typography variant="h6" component="div" noWrap sx={{ flexGrow: 1 }}>
          <Link href={useHref("/")} color={"inherit"} underline="none">
            PAN-GO
          </Link>
        </Typography>
        <Box sx={{ display: { md: "flex" } }}>
          <LangSelect></LangSelect>
          <Tooltip
            disableFocusListener
            disableTouchListener
            title={t("Header.Modules")}
          >
            <IconButton
              size="large"
              aria-label="Modules"
              color="inherit"
              href={useHref("/modules")}
            >
              <ViewModuleIcon />
            </IconButton>
          </Tooltip>
          <Tooltip
            disableFocusListener
            disableTouchListener
            title={t("Header.Settings")}
          >
            <IconButton
              size="large"
              aria-label="Settings"
              color="inherit"
              href={useHref("/settings")}
            >
              <SettingsIcon />
            </IconButton>
          </Tooltip>
        </Box>
      </Toolbar>
    </AppBar>
  );
}

export default Header;
