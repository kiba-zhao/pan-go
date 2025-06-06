import Paper from "@mui/material/Paper";
import Stack from "@mui/material/Stack";
import Box from "@mui/material/Box";
import Typography from "@mui/material/Typography";
import Button from "@mui/material/Button";
import { ExtFSPath } from "../ExtFS/Route";

import { useTranslation } from "../I18Next/Context";
import { PageI18Next } from "../I18Next/Page";
import { Link, To } from "../Route/Router";
import { PageHeader } from "./Header";
import { Fragment } from "react/jsx-runtime";

export const PageLayout = ({ children }: { children?: React.ReactNode }) => (
  <Paper sx={{ paddingY: 1, paddingX: 3, height: "100%", overflowY: "auto" }}>
    {children}
  </Paper>
);

export const PageHeaderTitle = ({
  children,
}: {
  children?: React.ReactNode;
}) => {
  return (
    <Typography noWrap variant="h6" component="div">
      {children}
    </Typography>
  );
};

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
          {text || t("resources.not_found.text")}
        </Typography>
        <Button variant="contained" component={Link} to={to || ExtFSPath}>
          {t("buttons.back")}
        </Button>
      </Box>
    </PageLayout>
  );
};

export const NotFoundPage = () => (
  <Fragment>
    <PageI18Next />
    <PageNotFound />
  </Fragment>
);
