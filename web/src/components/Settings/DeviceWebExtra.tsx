import { SettingsName } from "./meta";
import { DeviceWebPortInput } from "./DeviceWeb";
import {
  useWebPortField,
  useWebPortMutation,
  type HostSettings,
} from "./ReactQuery";
import { withExtraState, FormDialogExtra, type ExtraProps } from "./ExtraBase";

import { useEffect, useState, useMemo } from "react";
import { useForm } from "@/components/App/Form";
import { useAppDispatch } from "@/components/App/Context";
import { useTranslation, I18nVariant } from "@/components/App/I18Next";
import { toast, withToast, ToastType } from "@/components/App/Toast";

export const WebPortEditExtra = ({ extraState }: ExtraProps) => {
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
  } = useWebPortField({
    enabled: open,
  });
  const defaultValues = useMemo(
    () => ({ webPort: data?.webPort || 0 }),
    [data?.webPort],
  );

  const [error, setError] = useState<Error | undefined>();
  const clearErrors = () => {
    setError(undefined);
  };

  const { mutateAsync, isPending } = useWebPortMutation();
  const { reset, control, handleSubmit, setFocus } = useForm<
    Pick<HostSettings, "webPort">
  >({
    progressive: true,
    disabled: !open || isLoading || isPending || !isLoadingSuccess,
    defaultValues,
  });

  const onSuccess = (settings: HostSettings) => {
    toast.add(
      withToast(ToastType.Success, {
        description: t(`${I18nVariant.Extra}.${type}.success`),
      }),
    );
    handleClose();
  };

  const onError = (err: Error) => {
    setFocus("webPort");
    setError(err);
  };

  const handleFormSubmit = (values: Pick<HostSettings, "webPort">) => {
    error && clearErrors();
    mutateAsync(values).then(onSuccess, onError);
  };

  useEffect(() => {
    if (open && isLoadingSuccess) {
      setFocus("webPort");
      reset(defaultValues);
    }
  }, [open, defaultValues, isLoadingSuccess, setFocus, reset]);
  return (
    <FormDialogExtra
      extraState={extraState}
      control={control}
      badge={t(`${I18nVariant.Badge}.edit`)}
      desc={<p>{t(`${I18nVariant.Extra}.${type}.description`)}</p>}
      error={error}
      id="webPortEditForm"
      onSubmit={handleSubmit(handleFormSubmit)}
    >
      <DeviceWebPortInput
        name="webPort"
        control={control}
        defaultValue={defaultValues.webPort}
        onChange={clearErrors}
      />
    </FormDialogExtra>
  );
};
