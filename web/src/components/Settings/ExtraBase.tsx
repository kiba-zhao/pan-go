import { SettingsName } from "./meta";

import { withDialogExtraState } from "@/components/App/Extra";
import {
  useFormState,
  type UseFormStateProps,
  FieldValues,
} from "@/components/App/Form";
import { useTranslation, I18nVariant } from "@/components/App/I18Next";
import { useAppDispatch } from "@/components/App/Context";

import { type ComponentProps, type ReactNode } from "react";

import {
  Dialog,
  DialogErrorAlert,
  DialogProgressAlert,
  DialogAction,
  DialogDescription,
} from "@/components/App/Dialog";
import {
  DialogExtraClose,
  DialogExtraTitle,
  DialogExtraAction,
} from "@/components/App/Extra";

export enum ExtraType {
  DeviceNameEdit = "deviceNameEdit",
  DeviceMemoEdit = "deviceMemoEdit",
  PeerPortEdit = "peerPortEdit",
  BroadcastAddrsEdit = "broadcastAddrsEdit",
  PublicAddrsEdit = "publicAddrsEdit",
  WebPortEdit = "webPortEdit",
  ClusterSelect = "clusterSelect",
}

export type ExtraState = Parameters<typeof withDialogExtraState<ExtraType>>[0];
export { withDialogExtraState as withExtraState };

export type ExtraProps = { extraState: ExtraState };

type FormDialogExtraProps<T extends FieldValues> = ExtraProps &
  Required<Pick<UseFormStateProps<T>, "control">> & {
    title?: string;
    badge?: string;
    desc?: ReactNode;
    error?: Error;
    blank?: ReactNode;
  } & ComponentProps<"form">;
export const FormDialogExtra = <T extends FieldValues = FieldValues>({
  extraState,
  control,
  title,
  badge,
  desc,
  error,
  blank,
  id,
  children,
  ...props
}: FormDialogExtraProps<T>) => {
  const { type, open } = extraState;
  const { t } = useTranslation(SettingsName);
  const dispatch = useAppDispatch();

  const handleClose = () => {
    dispatch?.(withDialogExtraState({ type, open: false }));
  };

  const {
    disabled: formDisabled,
    isValid,
    isDirty,
    isSubmitting,
  } = useFormState<T>({ control });

  return (
    <Dialog open={open} onClose={handleClose}>
      <DialogExtraClose onClose={handleClose} />
      <DialogExtraTitle text={title || t(`${I18nVariant.Extra}.${type}.title`)}>
        {badge || t(`${I18nVariant.Extra}.${type}.subtitle`)}
      </DialogExtraTitle>
      <DialogDescription className={error || isSubmitting ? "hidden" : ""}>
        {desc}
      </DialogDescription>
      <DialogErrorAlert
        error={error}
        className={!error || isSubmitting ? "hidden" : ""}
      >
        {blank}
      </DialogErrorAlert>
      <DialogProgressAlert className={!isSubmitting ? "hidden" : ""}>
        {t(`${type}.saving`)}
        {blank}
      </DialogProgressAlert>
      <form id={id} {...props}>
        {children}
      </form>
      <DialogAction className="bg-muted">
        <DialogExtraAction
          type="submit"
          form={id}
          disabled={formDisabled || !isValid || !isDirty}
          progress={isSubmitting}
          onClose={handleClose}
        >
          {t(`${I18nVariant.Action}.save`)}
        </DialogExtraAction>
      </DialogAction>
    </Dialog>
  );
};
