/**
 * Master Component Definition File
 */

import { Fragment } from "react";
import { LocalesMenuButton } from "../I18Next/Menu";
import { Provider, ProviderView } from "./Context";
import { HeaderProvider, HeaderView } from "./Header";
import { ToggleColorModeButton } from "./Theme";

import Box from "@mui/material/Box";
import CssBaseline from "@mui/material/CssBaseline";
import { useTheme } from "@mui/material/styles";
import useMediaQuery from "@mui/material/useMediaQuery";

import type { ReactNode } from "react";

const MasterHeader = () => {
  const theme = useTheme();
  const isWideScreen = useMediaQuery(theme.breakpoints.up("sm"));
  return (
    <HeaderView>
      {isWideScreen && (
        <Fragment>
          <LocalesMenuButton />
          <ToggleColorModeButton />
        </Fragment>
      )}
    </HeaderView>
  );
};

const MasterMain = ({ children }: { children: ReactNode }) => (
  <Box className="master-main">{children}</Box>
);

const MasterProvider = ({ children }: { children: ReactNode }) => (
  <Provider>
    <HeaderProvider>{children}</HeaderProvider>
  </Provider>
);

type MasterProps = {
  children?: ReactNode;
};
export const Master = ({ children }: MasterProps) => (
  <MasterProvider>
    <ProviderView>
      <CssBaseline />
      <MasterHeader />
      <MasterMain>{children}</MasterMain>
    </ProviderView>
  </MasterProvider>
);
