/**
 * Master Component Definition File
 */

import { Fragment, type ReactNode } from "react";
import { ThemeProvider, ToggleColorModeButton } from "./Theme";
import { LocalesMenuButton, useTranslation } from "../i18n/Context";
import { HeaderView, HeaderProvider, PageHeaderTitle } from "./Header";
import { Provider, ProviderView } from "./Context";

import CssBaseline from "@mui/material/CssBaseline";
import Box from "@mui/material/Box";
import Card from "@mui/material/Card";
import CardContent from "@mui/material/CardContent";
import useMediaQuery from "@mui/material/useMediaQuery";
import { useTheme } from "@mui/material/styles";

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
  <ThemeProvider>
    <Provider>
      <HeaderProvider>{children}</HeaderProvider>
    </Provider>
  </ThemeProvider>
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

export const MasterErrorBoundary = () => {
  const { t } = useTranslation();

  return (
    <Master>
      <Card>
        <PageHeaderTitle title={t("ra.page.not_found")} />
        <CardContent>
          <h1>404: {t("ra.page.not_found")}</h1>
        </CardContent>
      </Card>
    </Master>
  );
};
