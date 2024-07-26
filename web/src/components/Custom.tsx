import { useEffect, useId, useMemo, useRef } from "react";
import {
  InfinitePagination as RAInfinitePagination,
  useListContext,
  useTranslate,
} from "react-admin";

import Box from "@mui/material/Box";
import Card from "@mui/material/Card";
import Typography from "@mui/material/Typography";

import FolderIcon from "@mui/icons-material/Folder";
import FilledInput from "@mui/material/FilledInput";
import FormControl from "@mui/material/FormControl";
import IconButton from "@mui/material/IconButton";
import InputAdornment from "@mui/material/InputAdornment";
import InputLabel from "@mui/material/InputLabel";

import type { QRCodeRenderersOptions } from "qrcode";
import { toCanvas } from "qrcode";

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
              {t("custom.pagination", { total })}
            </Typography>
          </Card>
        </Box>
      )}
    </>
  );
};

type InputProps<T = any> = {
  label: string;
} & T;

export type FilePathInputProps = InputProps<{}>;
export const FilePathInput = ({ label }: FilePathInputProps) => {
  const id = useId();
  const elementId = useMemo(() => `custom-filepath-input-${id}`, [id]);
  return (
    <FormControl variant="filled" fullWidth>
      <InputLabel htmlFor={elementId}>{label}</InputLabel>
      <FilledInput
        id={elementId}
        endAdornment={
          <InputAdornment position="end">
            <IconButton aria-label="directions">
              <FolderIcon />
            </IconButton>
          </InputAdornment>
        }
      />
    </FormControl>
  );
};

export type QRCodeProps = { value: string } & QRCodeRenderersOptions;
export const QRCode = ({ value, ...opts }: QRCodeProps) => {
  const canvasRef = useRef<HTMLCanvasElement | null>(null);

  useEffect(() => {
    toCanvas(canvasRef.current, value, opts);
  }, [value, opts]);

  return <canvas ref={canvasRef}></canvas>;
};
