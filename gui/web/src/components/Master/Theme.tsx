import {
  createTheme,
  ThemeProvider as MuiThemeProvider,
  useColorScheme,
} from "@mui/material/styles";
import useMediaQuery from "@mui/material/useMediaQuery";

import IconButton from "@mui/material/IconButton";
import Brightness4Icon from "@mui/icons-material/Brightness4";
import Brightness7Icon from "@mui/icons-material/Brightness7";

import { ReactNode, useMemo } from "react";

const defaultThemeOptions = {
  colorSchemes: {
    dark: true,
  },
};

export const ThemeProvider = ({ children }: { children: ReactNode }) => {
  const theme = useMemo(() => createTheme(defaultThemeOptions), []);
  return (
    <MuiThemeProvider theme={theme} defaultMode="system">
      {children}
    </MuiThemeProvider>
  );
};

export const ToggleColorModeButton = () => {
  const prefersDarkMode = useMediaQuery("(prefers-color-scheme: dark)");
  const { mode, setMode } = useColorScheme();
  const colorDarkMode = useMemo(
    () =>
      mode === "dark" || (mode === "system" && prefersDarkMode) ? true : false,
    [mode, prefersDarkMode]
  );

  if (!mode) {
    return null;
  }

  const handleClick = () => {
    setMode(colorDarkMode ? "light" : "dark");
  };
  return (
    <IconButton sx={{ ml: 1 }} onClick={handleClick} color="inherit">
      {colorDarkMode ? <Brightness7Icon /> : <Brightness4Icon />}
    </IconButton>
  );
};
