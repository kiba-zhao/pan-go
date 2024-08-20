import {
  InfinitePagination as RAInfinitePagination,
  useListContext,
} from "react-admin";

import { TranslateProvider, useTranslate } from "../Global/Translation";

import Box from "@mui/material/Box";
import Card from "@mui/material/Card";
import Typography from "@mui/material/Typography";

const Pagination = () => {
  const { total } = useListContext();
  const t = useTranslate();
  return (
    <>
      <RAInfinitePagination />
      {total && total > 0 && (
        <Box position="sticky" bottom={0} textAlign="center">
          <Card
            elevation={2}
            sx={{ px: 2, py: 1, mb: 1, display: "inline-block" }}
          >
            <Typography variant="body2">
              {t("pagination", { total })}
            </Typography>
          </Card>
        </Box>
      )}
    </>
  );
};
export const InfinitePagination = () => (
  <TranslateProvider>
    <Pagination />
  </TranslateProvider>
);
