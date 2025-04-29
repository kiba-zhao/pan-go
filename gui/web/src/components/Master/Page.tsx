import Paper from "@mui/material/Paper";
import Stack from "@mui/material/Stack";
import Box from "@mui/material/Box";
import Typography from "@mui/material/Typography";
import Button from "@mui/material/Button";
import { ExtFSPath } from "../ExtFS/Route";

import { useTranslation } from "../i18n/Context";
import { Link, To } from "../Route/Router";
import { PageHeader } from "./Header";

export const PageLayout = ({ children }: { children?: React.ReactNode }) => (
  <Paper sx={{ paddingY: 1, paddingX: 3, height: "100%", overflowY: "auto" }}>
    {children}
  </Paper>
);

export const PageTopBar = ({ children }: { children?: React.ReactNode }) => (
  <Stack direction="row" spacing={1} paddingTop={2} paddingBottom={3}>
    {children}
  </Stack>
);

export const PageNotFound = ({ to, text }: { to?: To; text?: string }) => {
  const { t } = useTranslation();
  return (
    <PageLayout>
      <PageHeader />
      <Box textAlign="center">
        <Typography variant="h1">404</Typography>
        <Typography variant="h6" paddingBottom={3}>
          {text || t("errors.PageNotFound")}
        </Typography>
        <Button variant="contained" component={Link} to={to || ExtFSPath}>
          {t("custom.button.back")}
        </Button>
      </Box>
    </PageLayout>
  );
};
