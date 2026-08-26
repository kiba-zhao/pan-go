import { SettingsName } from "./meta";
import {
  DeviceNetworkPortInput,
  parseBroadcastAddr,
  DeviceBroadcastAddrsTextInput,
  parsePublicAddr,
  DevicePublicAddrsTextInput,
} from "./DeviceNetwork";
import {
  useNetworkPortField,
  useNetworkPortMutation,
  useBroadcastAddrsField,
  useBroadcastAddrsMutation,
  usePublicAddrsField,
  usePublicAddrsMutation,
  type DeviceNetwork,
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
import {
  CardDialogExtra,
  DialogExtraClose,
  DialogExtraSpinnerAlert,
} from "@/components/App/Dialog";
import { toast, withToast, ToastType } from "@/components/App/Toast";

export const NetworkPortEditExtra = ({ extraState }: ExtraProps) => {
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
  } = useNetworkPortField({
    enabled: open,
  });
  const defaultValues = useMemo(
    () => ({
      port: data?.port || 0,
    }),
    [data?.port],
  );

  const [error, setError] = useState<Error | undefined>();
  const clearErrors = () => {
    setError(undefined);
  };

  const { mutateAsync, isPending } = useNetworkPortMutation();
  const { reset, control, handleSubmit, setFocus } = useForm<
    Pick<DeviceNetwork, "port">
  >({
    progressive: true,
    disabled: !open || isLoading || isPending || !isLoadingSuccess,
    defaultValues,
  });

  const onSuccess = (deviceNetwork: DeviceNetwork) => {
    toast.add(
      withToast(ToastType.Success, {
        description: t(`${I18nVariant.Extra}.${type}.success`),
      }),
    );
    handleClose();
  };

  const onError = (err: Error) => {
    setFocus("port");
    setError(err);
  };

  const handleFormSubmit = (values: Pick<DeviceNetwork, "port">) => {
    error && clearErrors();
    mutateAsync(values).then(onSuccess, onError);
  };

  useEffect(() => {
    if (open && isLoadingSuccess) {
      setFocus("port");
      reset(defaultValues);
    }
  }, [open, defaultValues, isLoadingSuccess, setFocus, reset]);

  return (
    <CardDialogExtra
      title={<DialogExtraTitle extraState={extraState} />}
      description={
        <FieldDialogDescription
          extraState={extraState}
          error={error}
          rotate={isPending}
        />
      }
      action={<DialogExtraClose />}
      footer={
        <DialogExtraFormAction control={control} form="networkPortEditForm">
          {t(`${I18nVariant.Action}.save`)}
        </DialogExtraFormAction>
      }
      footerClassName="flex justify-end"
    >
      <form id="networkPortEditForm" onSubmit={handleSubmit(handleFormSubmit)}>
        <DeviceNetworkPortInput
          name="port"
          control={control}
          defaultValue={defaultValues.port}
          onChange={clearErrors}
        />
      </form>
    </CardDialogExtra>
  );
};

export const BroadcastAddrsEditExtra = ({ extraState }: ExtraProps) => {
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
  } = useBroadcastAddrsField({
    enabled: open,
  });
  const defaultValues = useMemo(
    () => ({
      broadcastAddrsText: (data?.broadcastAddrs || []).join("\n"),
    }),
    [data?.broadcastAddrs],
  );

  const [error, setError] = useState<Error | undefined>();
  const clearErrors = () => {
    setError(undefined);
  };

  const { mutateAsync, isPending } = useBroadcastAddrsMutation();
  const { reset, control, handleSubmit, setFocus } = useForm<{
    broadcastAddrsText: string;
  }>({
    progressive: true,
    disabled: !open || isLoading || isPending || !isLoadingSuccess,
    defaultValues,
  });

  const onSuccess = (deviceNetwork: DeviceNetwork) => {
    toast.add(
      withToast(ToastType.Success, {
        description: t(`${I18nVariant.Extra}.${type}.success`),
      }),
    );
    handleClose();
  };

  const onError = (err: Error) => {
    setFocus("broadcastAddrsText");
    setError(err);
  };

  const handleFormSubmit = (values: { broadcastAddrsText: string }) => {
    error && clearErrors();
    const broadcastAddrs = values.broadcastAddrsText
      .split("\n")
      .map((_) => parseBroadcastAddr(_))
      .filter((_) => _ !== void 0);
    const broadcastAddrsSets = new Set(broadcastAddrs);
    mutateAsync({ broadcastAddrs: [...broadcastAddrsSets] }).then(
      onSuccess,
      onError,
    );
  };

  useEffect(() => {
    if (open && isLoadingSuccess) {
      setFocus("broadcastAddrsText");
      reset(defaultValues);
    }
  }, [open, defaultValues, isLoadingSuccess, setFocus]);
  return (
    <CardDialogExtra
      title={<DialogExtraTitle extraState={extraState} />}
      description={
        <FieldDialogDescription
          extraState={extraState}
          dataFormat={true}
          error={error}
          rotate={isPending}
        />
      }
      action={<DialogExtraClose />}
      footer={
        <DialogExtraFormAction control={control} form="broadcastAddrsEditForm">
          {t(`${I18nVariant.Action}.save`)}
        </DialogExtraFormAction>
      }
      footerClassName="flex justify-end"
    >
      <form
        id="broadcastAddrsEditForm"
        onSubmit={handleSubmit(handleFormSubmit)}
      >
        <DeviceBroadcastAddrsTextInput
          name="broadcastAddrsText"
          control={control}
          defaultValue={defaultValues.broadcastAddrsText}
          onChange={clearErrors}
        />
      </form>
    </CardDialogExtra>
  );
};

export const PublicAddrsEditExtra = ({ extraState }: ExtraProps) => {
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
  } = usePublicAddrsField({
    enabled: open,
  });
  const defaultValues = useMemo(
    () => ({
      publicAddrsText: (data?.publicAddrs || []).join("\n"),
    }),
    [data?.publicAddrs],
  );

  const [error, setError] = useState<Error | undefined>();
  const clearErrors = () => {
    setError(undefined);
  };

  const { mutateAsync, isPending } = usePublicAddrsMutation();
  const { reset, control, handleSubmit, setFocus } = useForm<{
    publicAddrsText: string;
  }>({
    progressive: true,
    disabled: !open || isLoading || isPending || !isLoadingSuccess,
    defaultValues,
  });

  const onSuccess = (deviceNetwork: DeviceNetwork) => {
    toast.add(
      withToast(ToastType.Success, {
        description: t(`${I18nVariant.Extra}.${type}.success`),
      }),
    );
    handleClose();
  };

  const onError = (err: Error) => {
    setFocus("publicAddrsText");
    setError(err);
  };

  const handleFormSubmit = (values: { publicAddrsText: string }) => {
    error && clearErrors();
    const publicAddrs = values.publicAddrsText
      .split("\n")
      .map((_) => parsePublicAddr(_))
      .filter((_) => _ !== void 0);
    const publicAddrsSets = new Set(publicAddrs);
    mutateAsync({ publicAddrs: [...publicAddrsSets] }).then(onSuccess, onError);
  };

  useEffect(() => {
    if (isLoadingSuccess) {
      setFocus("publicAddrsText");
      reset(defaultValues);
    }
  }, [open, defaultValues, isLoadingSuccess, setFocus, reset]);

  return (
    <CardDialogExtra
      title={<DialogExtraTitle extraState={extraState} />}
      description={
        <FieldDialogDescription
          extraState={extraState}
          dataFormat={true}
          error={error}
          rotate={isPending}
        />
      }
      action={<DialogExtraClose />}
      footer={
        <DialogExtraFormAction control={control} form="publicAddrsEditForm">
          {t(`${I18nVariant.Action}.save`)}
        </DialogExtraFormAction>
      }
      footerClassName="flex justify-end"
    >
      <form id="publicAddrsEditForm" onSubmit={handleSubmit(handleFormSubmit)}>
        <DevicePublicAddrsTextInput
          name="publicAddrsText"
          control={control}
          defaultValue={defaultValues.publicAddrsText}
          onChange={clearErrors}
        />
      </form>
    </CardDialogExtra>
  );
};

const FieldDialogDescription = ({
  extraState,
  dataFormat,
  error,
  rotate,
}: ExtraProps & {
  dataFormat?: boolean;
} & Pick<
    ComponentProps<typeof DialogExtraSpinnerAlert>,
    "error" | "rotate"
  >) => {
  const { type } = extraState;
  const { t } = useTranslation(SettingsName);

  return (
    <DialogExtraSpinnerAlert
      error={error}
      rotate={rotate}
      rotateContent={t(`${I18nVariant.Extra}.${type}.saving`)}
    >
      <p>{t(`${I18nVariant.Extra}.${type}.description`)}</p>
      {dataFormat && <p>{t(`${I18nVariant.Extra}.${type}.dataFormat`)}</p>}
    </DialogExtraSpinnerAlert>
  );
};
