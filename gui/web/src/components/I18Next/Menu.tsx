import { useMemo, useState } from "react";

import Box from "@mui/material/Box";
import Menu from "@mui/material/Menu";
import MenuItem from "@mui/material/MenuItem";
import LanguageIcon from "@mui/icons-material/Translate";
import IconButton from "@mui/material/IconButton";

import { LANGUAGES } from "./constants";
import { useTranslation } from "./Context";

type Language = {
  locale: string;
  name: string;
};

export const LocalesMenuButton = ({
  languages = LANGUAGES,
}: {
  languages?: Language[];
}) => {
  const { i18n } = useTranslation();

  const [anchorEl, setAnchorEl] = useState<null | HTMLElement>(null);
  const handleOpen = (element: HTMLElement) => {
    setAnchorEl(element);
  };
  const handleClose = (): void => {
    setAnchorEl(null);
  };

  const handleLocaleChange = (locale: string): void => {
    i18n.changeLanguage(locale);
    handleClose();
  };

  const language = useMemo(
    () => languages.find((l) => l.locale === i18n.language),
    [i18n.language, languages]
  );

  return (
    <Box component={"span"}>
      <IconButton
        sx={{ ml: 1 }}
        onClick={(event) => handleOpen(event.currentTarget)}
        color="inherit"
      >
        <LanguageIcon />
      </IconButton>
      <Menu
        anchorEl={anchorEl}
        keepMounted
        open={Boolean(anchorEl)}
        onClose={handleClose}
      >
        {languages.map((lang) => (
          <MenuItem
            key={lang.locale}
            onClick={() => handleLocaleChange(lang.locale)}
            selected={lang.locale === language?.locale}
          >
            {lang.name}
          </MenuItem>
        ))}
      </Menu>
    </Box>
  );
};
