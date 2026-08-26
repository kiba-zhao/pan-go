import { SettingsName } from "./meta";
import { DeviceNameInput, DeviceMemoInput } from "./DeviceInfo";
import {
  useDeviceNameField,
  useDeviceNameMutation,
  useDeviceMemoField,
  useDeviceMemoMutation,
  type DeviceInfo,
} from "./ReactQuery";
import {
  withExtraState,
  DialogExtraTitle,
  DialogExtraFormAction,
  type ExtraProps,
} from "./ExtraBase";

import { useEffect, useState, useMemo, type ComponentProps } from "react";
import { useForm } from "@/components/App/Form";
import { useAppDispatch } from "@/components/App/Context";
import { useTranslation, I18nVariant } from "@/components/App/I18Next";

import { toast, withToast, ToastType } from "@/components/App/Toast";
import {
  CardDialogExtra,
  DialogExtraClose,
  DialogExtraSpinnerAlert,
} from "@/components/App/Dialog";

export const DeviceNameEditExtra = ({ extraState }: ExtraProps) => {
  const { type, open } = extraState;
  const { t } = useTranslation(SettingsName);

  const dispatch = useAppDispatch();
  const handleClose = () => {
    dispatch?.(withExtraState({ type, open: false }));
  };

  const {
    data,
    isLoading,
    isSuccess: isLoadingSuccess,
  } = useDeviceNameField({
    enabled: open,
  });
  const defaultValues = useMemo(
    () => ({
      name: data?.name || "",
    }),
    [data?.name],
  );

  const [error, setError] = useState<Error | undefined>();
  const clearErrors = () => {
    setError(undefined);
  };

  const { mutateAsync, isPending } = useDeviceNameMutation();
  const { reset, control, handleSubmit, setFocus } = useForm<
    Pick<DeviceInfo, "name">
  >({
    progressive: true,
    disabled: !open || isLoading || isPending || !isLoadingSuccess,
    defaultValues,
  });

  const onSuccess = (deviceInfo: DeviceInfo) => {
    toast.add(
      withToast(ToastType.Success, {
        description: t(`${I18nVariant.Extra}.${type}.success`),
      }),
    );
    handleClose();
  };

  const onError = (err: Error) => {
    setFocus("name");
    setError(err);
  };

  const handleFormSubmit = (values: Pick<DeviceInfo, "name">) => {
    error && clearErrors();
    mutateAsync(values).then(onSuccess, onError);
  };

  useEffect(() => {
    if (open && isLoadingSuccess) {
      setFocus("name");
      reset(defaultValues);
    }
  }, [open, defaultValues, isLoadingSuccess, reset, setFocus]);

  return (
    <CardDialogExtra
      title={<DialogExtraTitle extraState={extraState} />}
      description={
        <TextFieldDescription
          extraState={extraState}
          error={error}
          rotate={isPending}
        />
      }
      action={<DialogExtraClose />}
      footer={
        <DialogExtraFormAction control={control} form="deviceNameEditForm">
          {t(`${I18nVariant.Action}.save`)}
        </DialogExtraFormAction>
      }
      footerClassName="flex justify-end"
    >
      <form id="deviceNameEditForm" onSubmit={handleSubmit(handleFormSubmit)}>
        <DeviceNameInput
          name="name"
          control={control}
          defaultValue={defaultValues.name}
          onChange={clearErrors}
        />
      </form>
    </CardDialogExtra>
  );
};

export const DeviceMemoEditExtra = ({ extraState }: ExtraProps) => {
  const { type, open } = extraState;
  const { t } = useTranslation(SettingsName);

  const dispatch = useAppDispatch();
  const handleClose = () => {
    dispatch?.(withExtraState({ type, open: false }));
  };

  const {
    data,
    isLoading,
    isSuccess: isLoadingSuccess,
  } = useDeviceMemoField({
    enabled: open,
  });

  const defaultValues = useMemo(
    () => ({
      memo: data?.memo || "",
    }),
    [data?.memo],
  );

  const [error, setError] = useState<Error | undefined>();
  const clearErrors = () => {
    setError(undefined);
  };

  const { mutateAsync, isPending } = useDeviceMemoMutation();
  const { reset, control, handleSubmit, setFocus } = useForm<
    Pick<DeviceInfo, "memo">
  >({
    progressive: true,
    disabled: !open || isLoading || isPending || !isLoadingSuccess,
    defaultValues,
  });

  const onSuccess = (deviceInfo: DeviceInfo) => {
    toast.add(
      withToast(ToastType.Success, {
        description: t(`${I18nVariant.Extra}.${type}.success`),
      }),
    );
    handleClose();
  };

  const onError = (err: Error) => {
    setFocus("memo");
    setError(err);
  };

  const handleFormSubmit = (values: Pick<DeviceInfo, "memo">) => {
    error && clearErrors();
    mutateAsync(values).then(onSuccess, onError);
  };

  useEffect(() => {
    if (open && isLoadingSuccess) {
      setFocus("memo");
      reset(defaultValues);
    }
  }, [open, defaultValues, isLoadingSuccess, reset, setFocus]);

  return (
    <CardDialogExtra
      title={<DialogExtraTitle extraState={extraState} />}
      description={
        <TextFieldDescription
          extraState={extraState}
          error={error}
          rotate={isPending}
        />
      }
      action={<DialogExtraClose />}
      footer={
        <DialogExtraFormAction control={control} form="deviceMemoEditForm">
          {t(`${I18nVariant.Action}.save`)}
        </DialogExtraFormAction>
      }
      footerClassName="flex justify-end"
    >
      <form id="deviceMemoEditForm" onSubmit={handleSubmit(handleFormSubmit)}>
        <DeviceMemoInput
          name="memo"
          control={control}
          defaultValue={defaultValues.memo}
          onChange={clearErrors}
        />
      </form>
    </CardDialogExtra>
  );
};

const TextFieldDescription = ({
  extraState,
  error,
  rotate,
}: ExtraProps &
  Pick<ComponentProps<typeof DialogExtraSpinnerAlert>, "error" | "rotate">) => {
  const { type } = extraState;
  const { t } = useTranslation(SettingsName);

  return (
    <DialogExtraSpinnerAlert
      error={error}
      rotate={rotate}
      rotateContent={t(`${I18nVariant.Extra}.${type}.saving`)}
    >
      <p>{t(`${I18nVariant.Extra}.${type}.description`)}</p>
      <p>{t(`${I18nVariant.Extra}.${type}.sizeLimit`)}</p>
    </DialogExtraSpinnerAlert>
  );
};
