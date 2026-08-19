import { SettingsName } from "./meta";
import {
  DeviceNetworkPortInput,
  parseBroadcastAddr,
  DeviceBroadcastAddrsTextInput,
  parsePublicAddr,
  DevicePublicAddrsTextInput,
} from "./DeviceNetwork";
import {
  usePeerPortField,
  usePeerPortMutation,
  useBroadcastAddrsField,
  useBroadcastAddrsMutation,
  usePublicAddrsField,
  usePublicAddrsMutation,
  type Settings,
} from "./ReactQuery";
import { withExtraState, FormDialogExtra, type ExtraProps } from "./ExtraBase";

import { useEffect, useState, useMemo } from "react";
import { useForm } from "@/components/App/Form";
import { useAppDispatch } from "@/components/App/Context";
import { useTranslation, I18nVariant } from "@/components/App/I18Next";

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
  } = usePeerPortField({
    enabled: open,
  });
  const defaultValues = useMemo(
    () => ({
      peerPort: data?.peerPort || 0,
    }),
    [data?.peerPort],
  );

  const [error, setError] = useState<Error | undefined>();
  const clearErrors = () => {
    setError(undefined);
  };

  const { mutateAsync, isPending } = usePeerPortMutation();
  const { reset, control, handleSubmit, setFocus } = useForm<
    Pick<Settings, "peerPort">
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
    setFocus("peerPort");
    setError(err);
  };

  const handleFormSubmit = (values: Pick<Settings, "peerPort">) => {
    error && clearErrors();
    mutateAsync(values).then(onSuccess, onError);
  };

  useEffect(() => {
    if (open && isLoadingSuccess) {
      setFocus("peerPort");
      reset(defaultValues);
    }
  }, [open, defaultValues, isLoadingSuccess, setFocus, reset]);

  return (
    <FormDialogExtra
      extraState={extraState}
      control={control}
      badge={t(`${I18nVariant.Badge}.edit`)}
      desc={<FieldDialogDescription extraState={extraState} />}
      error={error}
      id="networkPortEditForm"
      onSubmit={handleSubmit(handleFormSubmit)}
    >
      <DeviceNetworkPortInput
        name="peerPort"
        control={control}
        defaultValue={defaultValues.peerPort}
        onChange={clearErrors}
      />
    </FormDialogExtra>
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

  const onSuccess = (settings: Settings) => {
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
    <FormDialogExtra
      extraState={extraState}
      control={control}
      badge={t(`${I18nVariant.Badge}.edit`)}
      desc={
        <FieldDialogDescription extraState={extraState} dataFormat={true} />
      }
      error={error}
      id="broadcastAddrsEditForm"
      onSubmit={handleSubmit(handleFormSubmit)}
    >
      <DeviceBroadcastAddrsTextInput
        name="broadcastAddrsText"
        control={control}
        defaultValue={defaultValues.broadcastAddrsText}
        onChange={clearErrors}
      />
    </FormDialogExtra>
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

  const onSuccess = (settings: Settings) => {
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
    <FormDialogExtra
      extraState={extraState}
      control={control}
      badge={t(`${I18nVariant.Badge}.edit`)}
      desc={
        <FieldDialogDescription extraState={extraState} dataFormat={true} />
      }
      error={error}
      id="publicAddrsEditForm"
      onSubmit={handleSubmit(handleFormSubmit)}
    >
      <DevicePublicAddrsTextInput
        name="publicAddrsText"
        control={control}
        defaultValue={defaultValues.publicAddrsText}
        onChange={clearErrors}
      />
    </FormDialogExtra>
  );
};

const FieldDialogDescription = ({
  extraState,
  dataFormat,
}: ExtraProps & {
  dataFormat?: boolean;
}) => {
  const { type } = extraState;
  const { t } = useTranslation(SettingsName);
  return (
    <>
      <p>{t(`${I18nVariant.Extra}.${type}.description`)}</p>
      {dataFormat && <p>{t(`${I18nVariant.Extra}.${type}.dataFormat`)}</p>}
    </>
  );
};
