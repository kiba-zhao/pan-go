import {
  InfinitePagination as RAInfinitePagination,
  useListContext,
  useTranslate,
} from "react-admin";

import Box from "@mui/material/Box";
import Card from "@mui/material/Card";
import Typography from "@mui/material/Typography";

export const InfinitePagination = () => {
  const { total } = useListContext();
  const t = useTranslate();
  return (
    <>
      <RAInfinitePagination />
      {total > 0 && (
        <Box position="sticky" bottom={0} textAlign="center">
          <Card
            elevation={2}
            sx={{ px: 2, py: 1, mb: 1, display: "inline-block" }}
          >
            <Typography variant="body2">
              {t("others.pagination", { total })}
            </Typography>
          </Card>
        </Box>
      )}
    </>
  );
};
