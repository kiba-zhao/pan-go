import {
  createTheme,
  ThemeProvider as MuiThemeProvider,
  useColorScheme,
} from "@mui/material/styles";

import Brightness4Icon from "@mui/icons-material/Brightness4";
import Brightness7Icon from "@mui/icons-material/Brightness7";
import IconButton from "@mui/material/IconButton";
import InitColorSchemeScript from "@mui/material/InitColorSchemeScript";

import { ReactNode, useMemo } from "react";

const defaultThemeOptions = {
  colorSchemes: {
    dark: true,
    light: true,
  },
  cssVariables: {
    colorSchemeSelector: "class",
  },
};
export const ColorSchemeScript = () => (
  <InitColorSchemeScript attribute="class" />
);

export const ThemeProvider = ({ children }: { children: ReactNode }) => {
  const theme = useMemo(() => createTheme(defaultThemeOptions), []);
  return (
    <MuiThemeProvider theme={theme} defaultMode="system">
      {children}
    </MuiThemeProvider>
  );
};

export const ToggleColorModeButton = () => {
  const { mode, systemMode, setMode } = useColorScheme();
  const colorDarkMode = useMemo(
    () => mode === "dark" || (mode === "system" && systemMode === "dark"),
    [mode, systemMode]
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
