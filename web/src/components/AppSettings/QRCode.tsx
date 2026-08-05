import type { QRScanChangedEvent } from "../Common/QRCode";
import { QRCode, QRFileScan, QRProvider, QRScan } from "../Common/QRCode";
import { useTranslation } from "../I18Next/Context";

import type { ReactNode } from "react";
import { Fragment, useEffect, useMemo, useState } from "react";

import ImageIcon from "@mui/icons-material/Image";
import QrCodeScannerIcon from "@mui/icons-material/QrCodeScanner";
import Button from "@mui/material/Button";
import Dialog from "@mui/material/Dialog";
import DialogActions from "@mui/material/DialogActions";
import DialogContent from "@mui/material/DialogContent";
import Stack from "@mui/material/Stack";
import { useBrowser } from "../App/Browser";

export type NodeQRCodeValue = {
  name: string;
  peerId: string;
};

type NodeQRCodeProps = {
  name?: string;
  peerId?: string;
  width?: number;
  children?: ReactNode;
  base?: string;
};
export const NodeQRCode = ({
  name = "",
  peerId = "",
  width = 200,
  children,
}: NodeQRCodeProps) => {
  const value = useMemo(() => {
    if (!name || name.length <= 0 || !peerId || peerId.length <= 0) return;
    return toNodeQRCodeUrl({ name, peerId });
  }, [name, peerId]);

  return (
    <QRProvider>
      <Stack spacing={2} alignItems="center" justifyContent={"space-between"}>
        <QRCode name={name} value={value} width={width} />
        {children}
      </Stack>
    </QRProvider>
  );
};

export const NodeQRScan = (props: NodeQRScanProps) => {
  const { onQRScan } = props;
  const { t } = useTranslation();

  const browser = useBrowser();
  const [available, setAvailable] = useState(false);
  useEffect(() => {
    const { window } = browser || {};
    if (!window) return;
    navigator.mediaDevices.enumerateDevices().then((devices) => {
      const hasWebcam = devices.some((_) => _.kind === "videoinput");
      if (hasWebcam !== available) setAvailable(hasWebcam);
    });
  }, [props]);

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
        disabled={!available}
      >
        {t("buttons.qrscan")}
      </Button>
      {available && open && (
        <Dialog open={open} onClose={onClose}>
          <DialogActions>
            <Button size="small" onClick={onClose}>
              {t("buttons.close")}
            </Button>
          </DialogActions>
          <DialogContent>
            <QRScan onChanged={onChanged} width="400" height="300" />
          </DialogContent>
        </Dialog>
      )}
    </Fragment>
  );
};

type NodeQRScanProps = {
  onQRScan?: (value: NodeQRCodeValue) => void;
};
export const NodeFileQRScan = ({ onQRScan }: NodeQRScanProps) => {
  const { t } = useTranslation();

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
        {t("buttons.qrscan-file")}
        <QRFileScan onChange={onFileScan} />
      </Button>
    </Fragment>
  );
};

function parseNodeQRCodeValue(value: string): NodeQRCodeValue | undefined {
  if (!URL.canParse(value)) return;
  const url = new URL(value);
  if (url.protocol.slice(0, -1) !== import.meta.env.VITE_APP_NAME) return;
  if (
    !(url.host === "app" && url.pathname === "/nodes") &&
    !(url.host === "" && url.pathname === "//app/nodes")
  )
    return;
  const peerId = url.searchParams.get("peerId");
  const name = url.searchParams.get("name");
  if (!peerId || peerId.length <= 0 || !name || name.length <= 0) return;
  return { name, peerId };
}

function toNodeQRCodeUrl({ name, peerId }: NodeQRCodeValue): string {
  const query = new URLSearchParams({ peerId, name });
  return `${import.meta.env.VITE_APP_NAME}://app/nodes?${query.toString()}`;
}
