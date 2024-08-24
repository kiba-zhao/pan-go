import { TranslateProvider, useTranslate } from "../Global/Translation";
import type { QRScanChangedEvent } from "./Base";
import { QRCode, QRFileScan, QRProvider, QRScan } from "./Base";

import type { ReactNode } from "react";
import { Fragment, useMemo, useState } from "react";

import ImageIcon from "@mui/icons-material/Image";
import QrCodeScannerIcon from "@mui/icons-material/QrCodeScanner";
import Button from "@mui/material/Button";
import Dialog from "@mui/material/Dialog";
import DialogActions from "@mui/material/DialogActions";
import DialogContent from "@mui/material/DialogContent";
import Stack from "@mui/material/Stack";

export type NodeQRCodeValue = {
  name: string;
  nodeId: string;
};

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
    return toNodeQRCodeUrl({ name, nodeId });
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

const InternalQRScan = ({ onQRScan }: NodeQRScanProps) => {
  const t = useTranslate();

  const [open, setOpen] = useState(false);
  const onOpen = () => setOpen(true);
  const onClose = () => setOpen(false);
  const onChanged = (event: QRScanChangedEvent) => {
    const v = parseNodeQRCodeValue(event.value);
    if (!v) return;
    onQRScan?.(v);
    event.invalid = false;
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
        <DialogActions>
          <Button size="small" onClick={onClose}>
            {t("button.close")}
          </Button>
        </DialogActions>
        <DialogContent>
          <QRScan onChanged={onChanged} width="400" height="300" />
        </DialogContent>
      </Dialog>
    </Fragment>
  );
};
export const NodeQRScan = (props: NodeQRScanProps) => (
  <TranslateProvider>
    <InternalQRScan {...props} />
  </TranslateProvider>
);

type NodeQRScanProps = {
  onQRScan?: (value: NodeQRCodeValue) => void;
};
const InternalQRFileScan = ({ onQRScan }: NodeQRScanProps) => {
  const t = useTranslate();

  const onFileScan = (value: string) => {
    const nodeValue = parseNodeQRCodeValue(value);
    if (!nodeValue) return;
    onQRScan?.(nodeValue);
  };
  return (
    <Fragment>
      <Button
        component="label"
        variant="contained"
        size="small"
        startIcon={<ImageIcon />}
      >
        {t("button.qrscan-file")}
        <QRFileScan onChange={onFileScan} />
      </Button>
    </Fragment>
  );
};

export const NodeFileQRScan = (props: NodeQRScanProps) => (
  <TranslateProvider>
    <InternalQRFileScan {...props} />
  </TranslateProvider>
);

function parseNodeQRCodeValue(value: string): NodeQRCodeValue | undefined {
  if (!URL.canParse(value)) return;
  const url = new URL(value);
  if (url.protocol.slice(0, -1) !== import.meta.env.VITE_APP_NAME) return;
  if (url.pathname !== "//app/node") return;
  const nodeId = url.searchParams.get("nodeId");
  const name = url.searchParams.get("name");
  if (!nodeId || nodeId.length <= 0 || !name || name.length <= 0) return;
  return { name, nodeId };
}

function toNodeQRCodeUrl({ name, nodeId }: NodeQRCodeValue): string {
  const query = new URLSearchParams({ nodeId, name });
  return `${import.meta.env.VITE_APP_NAME}://app/node?${query.toString()}`;
}
