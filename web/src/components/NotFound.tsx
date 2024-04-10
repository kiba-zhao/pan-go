import Card from "@mui/material/Card";
import CardContent from "@mui/material/CardContent";
import { Title, useTranslate } from "react-admin";

export default () => {
  const t = useTranslate();
  return (
    <Card>
      <Title title={t("ra.page.not_found")} />
      <CardContent>
        <h1>404: {t("ra.page.not_found")}</h1>
      </CardContent>
    </Card>
  );
};
