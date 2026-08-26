import { SettingsName } from "./meta";
import { DeviceWebPortInput } from "./DeviceWeb";
import {
  useWebPortField,
  useWebPortMutation,
  type WebHost,
} from "./ReactQuery";
import {
  withExtraState,
  DialogExtraTitle,
  DialogExtraFormAction,
  type ExtraProps,
} from "./ExtraBase";

import { useEffect, useState, useMemo } from "react";
import { useForm } from "@/components/App/Form";
import { useAppDispatch } from "@/components/App/Context";
import { useTranslation, I18nVariant } from "@/components/App/I18Next";
import { toast, withToast, ToastType } from "@/components/App/Toast";
import {
  CardDialogExtra,
  DialogExtraClose,
  DialogExtraSpinnerAlert,
} from "@/components/App/Dialog";

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
    Pick<WebHost, "webPort">
  >({
    progressive: true,
    disabled: !open || isLoading || isPending || !isLoadingSuccess,
    defaultValues,
  });

  const onSuccess = (webHost: WebHost) => {
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

  const handleFormSubmit = (values: Pick<WebHost, "webPort">) => {
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
    <CardDialogExtra
      title={<DialogExtraTitle extraState={extraState} />}
      description={
        <DialogExtraSpinnerAlert
          error={error}
          rotate={isPending}
          rotateContent={t(`${I18nVariant.Extra}.${type}.saving`)}
        >
          {t(`${I18nVariant.Extra}.${type}.description`)}
        </DialogExtraSpinnerAlert>
      }
      action={<DialogExtraClose />}
      footer={
        <DialogExtraFormAction control={control} form="webPortEditForm">
          {t(`${I18nVariant.Action}.save`)}
        </DialogExtraFormAction>
      }
      footerClassName="flex justify-end"
    >
      <form id="webPortEditForm" onSubmit={handleSubmit(handleFormSubmit)}>
        <DeviceWebPortInput
          name="webPort"
          control={control}
          defaultValue={defaultValues.webPort}
          onChange={clearErrors}
        />
      </form>
    </CardDialogExtra>
  );
};
