import { TranslateProvider, useTranslate } from "../Global/Translation";

import Button from "@mui/material/Button";
import MuiDialog from "@mui/material/Dialog";
import MuiDialogActions from "@mui/material/DialogActions";
import MuiDialogContent from "@mui/material/DialogContent";
import MuiDialogTitle from "@mui/material/DialogTitle";

import { Fragment, type ReactNode } from "react";

export type DialogProps = {
  open: boolean;
  onClose?: () => void;
  title?: ReactNode;
  children?: ReactNode;
  actions?: ReactNode;
};
export const Dialog = ({
  open,
  onClose,
  title,
  children,
  actions,
}: DialogProps) => (
  <TranslateProvider>
    <MuiDialog open={open} onClose={onClose}>
      {title && <MuiDialogTitle>{title}</MuiDialogTitle>}
      {children && <MuiDialogContent>{children}</MuiDialogContent>}
      {actions && <MuiDialogActions>{actions}</MuiDialogActions>}
    </MuiDialog>
  </TranslateProvider>
);

export const DialogSubmitActions = ({
  onSubmit,
  onCancel,
}: {
  onSubmit: () => void;
  onCancel: () => void;
}) => {
  const t = useTranslate();
  return (
    <Fragment>
      <Button size="small" onClick={onCancel}>
        {t("button.cancel")}
      </Button>
      <Button size="small" onClick={onSubmit} autoFocus>
        {t("button.submit")}
      </Button>
    </Fragment>
  );
};

export const DialogConfirmActions = ({
  onConfirm,
}: {
  onConfirm: (confirm: boolean) => void;
}) => {
  const t = useTranslate();
  return (
    <Fragment>
      <Button size="small" onClick={() => onConfirm(false)}>
        {t("button.cancel")}
      </Button>
      <Button size="small" onClick={() => onConfirm(true)} autoFocus>
        {t("button.confirm")}
      </Button>
    </Fragment>
  );
};
