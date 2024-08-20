import { TranslateProvider, useTranslate } from "../Global/Translation";
import type { QRScanChangedEvent } from "./Base";
import { QRCode, QRProvider, QRScan } from "./Base";

import type { ReactNode } from "react";
import { Fragment, useMemo, useState } from "react";

import QrCodeScannerIcon from "@mui/icons-material/QrCodeScanner";
import Button from "@mui/material/Button";
import Dialog from "@mui/material/Dialog";
import DialogActions from "@mui/material/DialogActions";
import DialogContent from "@mui/material/DialogContent";
import Stack from "@mui/material/Stack";

type NodeQRCodeProps = {
  name?: string;
  nodeId?: string;
  width?: number;
  children?: ReactNode;
  base?: string;
};
export const NodeQRCode = ({
  name = "",
  nodeId = "",
  width = 200,
  children,
}: NodeQRCodeProps) => {
  const value = useMemo(() => {
    if (!name || name.length <= 0 || !nodeId || nodeId.length <= 0) return;
    const query = new URLSearchParams({ nodeId, name });
    return `${import.meta.env.VITE_APP_NAME}://app/node?${query.toString()}`;
  }, [name, nodeId]);

  return (
    <QRProvider>
      <Stack spacing={2} alignItems="center" justifyContent={"space-between"}>
        <QRCode name={name} value={value} width={width} />
        {children}
      </Stack>
    </QRProvider>
  );
};

const InternalQRScan = () => {
  const t = useTranslate();

  const [open, setOpen] = useState(false);
  const onOpen = () => setOpen(true);
  const onClose = () => setOpen(false);
  const onChanged = (event: QRScanChangedEvent) => {
    if (!URL.canParse(event.value)) return;
    const url = new URL(event.value);
    if (url.protocol !== import.meta.env.VITE_APP_NAME) return;
    if (url.host !== "app") return;
    if (url.pathname !== "/node") return;
    const nodeId = url.searchParams.get("nodeId");
    const name = url.searchParams.get("name");
    if (!nodeId || nodeId.length <= 0 || !name || name.length <= 0) return;
    event.invalid = false;
    console.log(11111, nodeId, name);
    setOpen(false);
  };
  return (
    <Fragment>
      <Button
        variant="contained"
        size="small"
        startIcon={<QrCodeScannerIcon />}
        onClick={onOpen}
      >
        {t("button.qrscan")}
      </Button>
      <Dialog open={open} onClose={onClose}>
        <DialogContent>
          <QRScan onChanged={onChanged} width="400" height="300" />
        </DialogContent>
        <DialogActions>
          <Button size="small" onClick={onClose}>
            {t("button.cancel")}
          </Button>
        </DialogActions>
      </Dialog>
    </Fragment>
  );
};
export const NodeQRScan = () => (
  <TranslateProvider>
    <InternalQRScan />
  </TranslateProvider>
);
