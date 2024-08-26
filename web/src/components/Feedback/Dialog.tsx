import { TranslateProvider, useTranslate } from "../Global/Translation";

import Button from "@mui/material/Button";
import type { DialogProps as MuiDialogProps } from "@mui/material/Dialog";
import MuiDialog from "@mui/material/Dialog";
import DialogContent from "@mui/material/DialogContent";
import DialogContentText from "@mui/material/DialogContentText";
import DialogTitle from "@mui/material/DialogTitle";

import { Fragment } from "react";

export type DialogProps = MuiDialogProps;
export const Dialog = ({ children, ...props }: DialogProps) => (
  <TranslateProvider>
    <MuiDialog {...props}>{children}</MuiDialog>
  </TranslateProvider>
);

export const DialogSubmitActions = ({
  onSubmit,
  onCancel,
  label,
}: {
  onSubmit: () => void;
  onCancel: () => void;
  label?: string;
}) => {
  const t = useTranslate();
  return (
    <Fragment>
      <Button size="small" onClick={onCancel}>
        {t("button.cancel")}
      </Button>
      <Button size="small" onClick={onSubmit} autoFocus>
        {label ? label : t("button.submit")}
      </Button>
    </Fragment>
  );
};

export const DialogConfirmActions = ({
  onConfirm,
  label,
}: {
  onConfirm: (confirm: boolean) => void;
  label?: string;
}) => {
  const t = useTranslate();
  return (
    <Fragment>
      <Button size="small" onClick={() => onConfirm(false)}>
        {t("button.cancel")}
      </Button>
      <Button
        size="small"
        onClick={() => onConfirm(true)}
        autoFocus
        color="error"
      >
        {label ? label : t("button.confirm")}
      </Button>
    </Fragment>
  );
};

type DialogConfirmContentProps = {
  label?: string;
  contentLabel?: string;
};
export const DialogConfirmContent = ({
  label = "",
  contentLabel,
}: DialogConfirmContentProps) => {
  const t = useTranslate();
  return (
    <Fragment>
      <DialogTitle>{t("dialog.confirm.title", { label })}</DialogTitle>
      <DialogContent>
        <DialogContentText>
          {t("dialog.confirm.content", {
            label: contentLabel || label.toLowerCase(),
          })}
        </DialogContentText>
      </DialogContent>
    </Fragment>
  );
};
