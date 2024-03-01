import Container from "@mui/material/Container";
import Typography from "@mui/material/Typography";
import { useTranslation } from "react-i18next";

function NotFound() {
  const { t } = useTranslation();
  return (
    <Container component="main" sx={{ mt: 8, mb: 2 }} maxWidth="sm">
      <Typography variant="h2" component="h1" gutterBottom>
        {t("errors.NotFound")}
      </Typography>
      <Typography variant="h5" component="h2" gutterBottom>
        {t("errors.NotFoundDesc")}
      </Typography>
    </Container>
  );
}

export default NotFound;
