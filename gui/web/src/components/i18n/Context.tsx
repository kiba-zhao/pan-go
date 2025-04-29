/**
 * i18next Provider for React Admin
 *
 * @see https://react.i18next.com/
 */
import i18next from "i18next";
import LanguageDetector from "i18next-browser-languagedetector";
import resourcesToBackend from "i18next-resources-to-backend";
import {
  I18nextProvider,
  useTranslation,
  initReactI18next,
} from "react-i18next";
import { useMemo, useState } from "react";

import Box from "@mui/material/Box";
import Menu from "@mui/material/Menu";
import MenuItem from "@mui/material/MenuItem";
import LanguageIcon from "@mui/icons-material/Translate";
import IconButton from "@mui/material/IconButton";

export { useTranslation };

const defaultLanguages = [
  { locale: "zh-CN", name: "中文" },
  { locale: "en", name: "English" },
];

const importLanguage = async (language: string, namespace: string) => {
  const { default: translation } = await import(
    `../../locales/${language}/${namespace}.json`
  );

  return translation;
};

i18next
  .use(LanguageDetector)
  .use(initReactI18next)
  .use(resourcesToBackend(importLanguage))
  .init({
    lng: defaultLanguages[0].locale,
    fallbackLng: defaultLanguages[0].locale,
    supportedLngs: defaultLanguages.map((l) => l.locale),
  });

export const I18nProvider = ({ children }: { children: React.ReactNode }) => (
  <I18nextProvider i18n={i18next} defaultNS={"translation"}>
    {children}
  </I18nextProvider>
);

type Language = {
  locale: string;
  name: string;
};

export const LocalesMenuButton = ({
  languages = defaultLanguages,
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
    () => languages.find((l) => l.locale === i18next.language),
    [i18next.language, languages]
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
