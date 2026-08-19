import { SettingsName } from "./meta";
import { DeviceNameInput, DeviceMemoInput } from "./DeviceInfo";
import {
  useDeviceNameField,
  useDeviceNameMutation,
  useDeviceMemoField,
  useDeviceMemoMutation,
  type Settings,
} from "./ReactQuery";
import { withExtraState, FormDialogExtra, type ExtraProps } from "./ExtraBase";

import { useEffect, useState, useMemo } from "react";
import { useForm } from "@/components/App/Form";
import { useAppDispatch } from "@/components/App/Context";
import { useTranslation, I18nVariant } from "@/components/App/I18Next";

import { toast, withToast, ToastType } from "@/components/App/Toast";
import { Blank } from "@/components/App/Typography";

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
    Pick<Settings, "name">
  >({
    progressive: true,
    disabled: !open || isLoading || isPending || !isLoadingSuccess,
    defaultValues,
  });

  const onSuccess = (settings: Settings) => {
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

  const handleFormSubmit = (values: Pick<Settings, "name">) => {
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
    <FormDialogExtra
      extraState={extraState}
      control={control}
      badge={t(`${I18nVariant.Badge}.edit`)}
      desc={<TextFieldDialogDescription extraState={extraState} />}
      error={error}
      blank={<Blank />}
      id="deviceNameEditForm"
      onSubmit={handleSubmit(handleFormSubmit)}
    >
      <DeviceNameInput
        name="name"
        control={control}
        defaultValue={defaultValues.name}
        onChange={clearErrors}
      />
    </FormDialogExtra>
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
    Pick<Settings, "memo">
  >({
    progressive: true,
    disabled: !open || isLoading || isPending || !isLoadingSuccess,
    defaultValues,
  });

  const onSuccess = (settings: Settings) => {
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

  const handleFormSubmit = (values: Pick<Settings, "memo">) => {
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
    <FormDialogExtra
      extraState={extraState}
      control={control}
      badge={t(`${I18nVariant.Badge}.edit`)}
      desc={<TextFieldDialogDescription extraState={extraState} />}
      error={error}
      blank={<Blank />}
      id="deviceMemoEditForm"
      onSubmit={handleSubmit(handleFormSubmit)}
    >
      <DeviceMemoInput
        name="memo"
        control={control}
        defaultValue={defaultValues.memo}
        onChange={clearErrors}
      />
    </FormDialogExtra>
  );
};

const TextFieldDialogDescription = ({ extraState }: ExtraProps) => {
  const { type } = extraState;
  const { t } = useTranslation(SettingsName);
  return (
    <>
      <p>{t(`${I18nVariant.Extra}.${type}.description`)}</p>
      <p>{t(`${I18nVariant.Extra}.${type}.sizeLimit`)}</p>
    </>
  );
};
