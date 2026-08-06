import {
  useAppExtra,
  useAppDispatch,
  withAppExtraAction,
} from "@/components/App/Context";
import {
  Dialog,
  DialogTitle,
  DialogDescription,
  DialogAction,
  DialogExtra,
  DialogAlert,
  AlertDescription,
} from "@/components/App/Dialog";
import { X as CloseIcon, CircleX, CircleCheckBig, CircleAlert } from "./Icon";
import {
  InputGroup,
  InputGroupAddon,
  InputGroupInput,
  InputGroupText,
  InputGroupButton,
  InputGroupTextarea,
} from "@/components/ui/input-group";
import { Blank } from "@/components/App/Typography";
import { Button } from "@/components/ui/button";
import { Badge } from "@/components/ui/badge";
import { ExtraState as AppExtraState } from "@/components/App/Extra";
import { useEffect, useState, useMemo } from "react";

import { useTranslation, DefaultNS } from "@/components/App/I18Next";
import { SettingsName, ExtraI18nPrefix } from "./meta";
import {
  useNameField,
  useNameMutation,
  useMemoField,
  useMemoMutation,
  usePeerPortField,
  usePeerPortMutation,
  useWebPortField,
  useWebPortMutation,
  useBroadcastAddrsField,
  useBroadcastAddrsMutation,
  usePublicAddrsField,
  usePublicAddrsMutation,
  type Settings,
  type HostSettings,
} from "./ReactQuery";
import { useForm } from "@/components/App/Form";
import { Spinner } from "@/components/ui/spinner";
import {
  toast,
  withToast,
  ToastType,
  type ToastID,
} from "@/components/App/Toast";
import { newIPAddr } from "@/lib/utils";

const SettingsExtra = () => (
  <>
    <NameEditExtra />
    <MemoEditExtra />
    <PeerPortEditExtra />
    <BroadcastAddrsEditExtra />
    <PublicAddrsEditExtra />
    <WebPortEditExtra />
  </>
);

export default SettingsExtra;

export enum ExtraType {
  NameEdit = "nameEdit",
  MemoEdit = "memoEdit",
  PeerPortEdit = "peerPortEdit",
  BroadcastAddrsEdit = "broadcastAddrsEdit",
  PublicAddrsEdit = "publicAddrsEdit",
  WebPortEdit = "webPortEdit",
}
export type ExtraState<Type extends ExtraType> = AppExtraState<
  Type,
  { open?: boolean }
>;

export const withExtraAction = (type?: ExtraType, open?: boolean) => {
  return withAppExtraAction(
    type === void 0 ? type : { type, open: open !== false },
  );
};

const CommonI18nOpts = { ns: [SettingsName, DefaultNS], keyPrefix: "" };
const NameEditKeyPrefix = `${ExtraI18nPrefix}.${ExtraType.NameEdit}`;
const NameEditExtra = () => {
  const { t } = useTranslation(SettingsName, {
    keyPrefix: NameEditKeyPrefix,
  });

  const { type, open } = useAppExtra<ExtraState<ExtraType.NameEdit>>({});
  const dispatch = useAppDispatch();

  const handleClose = () => {
    dispatch?.(withExtraAction(type, false));
  };

  const {
    data,
    isLoading,
    isSuccess: isLoadingSuccess,
  } = useNameField({
    enabled: type === ExtraType.NameEdit && open,
  });

  const [error, setError] = useState<Error | undefined>();
  const clearErrors = () => {
    setError(undefined);
  };

  const { mutateAsync, isPending } = useNameMutation();
  const {
    reset,
    watch,
    register,
    handleSubmit,
    formState: { errors, disabled: formDisabled, isValid, isDirty },
    setFocus,
  } = useForm<Pick<Settings, "name">>({
    progressive: true,
    disabled:
      type !== ExtraType.NameEdit ||
      !open ||
      isLoading ||
      isPending ||
      !isLoadingSuccess,
  });

  const handleReset = () => {
    reset({ name: data?.name || "" });
    clearErrors();
  };

  const onSuccess = (settings: Settings) => {
    toast.add(
      withToast(ToastType.Success, {
        description: t("success"),
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
    if (type === ExtraType.NameEdit && open && isLoadingSuccess) {
      setFocus("name");
      handleReset();
    }
  }, [type, open, isLoadingSuccess, setFocus]);

  const nameValue = watch("name");

  if (type !== ExtraType.NameEdit) {
    return null;
  }

  return (
    <Dialog open={open} onClose={handleClose}>
      <DialogExtra>
        <Button size="icon-sm" variant="ghost" onClick={handleClose}>
          <CloseIcon />
        </Button>
      </DialogExtra>
      <DialogTitle>
        {t("title")}
        <small className="text-muted-foreground px-1">
          {t("badge.edit", CommonI18nOpts)}
        </small>
      </DialogTitle>
      <DialogDescription className={error || isPending ? "hidden" : ""}>
        <p>{t("description")}</p>
        <p>{t("sizeLimit")}</p>
      </DialogDescription>
      <DialogAlert
        className={!isPending && error ? void 0 : "hidden"}
        variant="destructive"
      >
        <CircleAlert />
        <AlertDescription>
          {t(error?.name || "Error", {
            ...CommonI18nOpts,
            defaultValue: error?.message || "",
          })}
          <Blank />
        </AlertDescription>
      </DialogAlert>
      <DialogAlert className={isPending ? void 0 : "hidden"}>
        <Spinner />
        <AlertDescription>
          {t("saving")}
          <Blank />
        </AlertDescription>
      </DialogAlert>
      <form id="nameEditForm" onSubmit={handleSubmit(handleFormSubmit)}>
        <InputGroup>
          <InputGroupInput
            {...register("name", {
              required: true,
              maxLength: 32,
              pattern: /^\S+$/g,
              onChange: clearErrors,
            })}
            type="text"
            placeholder={t("placeholder")}
            autoComplete="off"
            aria-invalid={errors.name ? "true" : "false"}
          />
          <InputGroupAddon align="block-end">
            <InputGroupText className="text-xs text-muted-foreground">
              {nameValue?.length || 0} / 32
            </InputGroupText>
            <InputGroupButton
              variant="default"
              size="sm"
              className="ml-auto max-md:hidden"
              onClick={handleReset}
              disabled={formDisabled || !isDirty}
            >
              {t("action.reset", CommonI18nOpts)}
            </InputGroupButton>
          </InputGroupAddon>
        </InputGroup>
      </form>
      <DialogAction className="bg-muted">
        <Button
          className="max-md:hidden"
          onClick={handleClose}
          variant="outline"
        >
          {t("action.cancel", CommonI18nOpts)}
        </Button>
        <Button
          type="submit"
          form="nameEditForm"
          disabled={formDisabled || !isValid || !isDirty}
        >
          <Spinner className={isPending ? "" : "hidden"} />
          {t("action.save", CommonI18nOpts)}
        </Button>
      </DialogAction>
    </Dialog>
  );
};

const MemoEditKeyPrefix = `${ExtraI18nPrefix}.${ExtraType.MemoEdit}`;
const MemoEditExtra = () => {
  const { t } = useTranslation(SettingsName, {
    keyPrefix: MemoEditKeyPrefix,
  });

  const { type, open } = useAppExtra<ExtraState<ExtraType.MemoEdit>>({});
  const dispatch = useAppDispatch();

  const handleClose = () => {
    dispatch?.(withExtraAction(type, false));
  };

  const {
    data,
    isLoading,
    isSuccess: isLoadingSuccess,
  } = useMemoField({
    enabled: type === ExtraType.MemoEdit && open,
  });

  const [error, setError] = useState<Error | undefined>();
  const clearErrors = () => {
    setError(undefined);
  };

  const { mutateAsync, isPending } = useMemoMutation();
  const {
    reset,
    watch,
    register,
    handleSubmit,
    formState: { errors, disabled: formDisabled, isValid, isDirty },
    setFocus,
  } = useForm<Pick<Settings, "memo">>({
    progressive: true,
    disabled:
      type !== ExtraType.MemoEdit ||
      !open ||
      isLoading ||
      isPending ||
      !isLoadingSuccess,
  });

  const handleReset = () => {
    reset({ memo: data?.memo || "" });
    clearErrors();
  };

  const onSuccess = (settings: Settings) => {
    toast.add(
      withToast(ToastType.Success, {
        description: t("success"),
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
    if (type === ExtraType.MemoEdit && open && isLoadingSuccess) {
      setFocus("memo");
      handleReset();
    }
  }, [type, open, isLoadingSuccess, setFocus]);

  const memoValue = watch("memo");

  const setMemoAs = (value?: string) => {
    if (!value) {
      return "";
    }
    const lines = value
      .trim()
      .split("\n")
      .filter((_) => !/^\s+$/g.test(_))
      .map((_) => _.trim());
    return lines.join("\n");
  };

  if (type !== ExtraType.MemoEdit) {
    return null;
  }

  return (
    <Dialog open={open} onClose={handleClose}>
      <DialogExtra>
        <Button size="icon-sm" variant="ghost" onClick={handleClose}>
          <CloseIcon />
        </Button>
      </DialogExtra>
      <DialogTitle>
        {t("title")}
        <small className="text-muted-foreground px-1">
          {t("badge.edit", CommonI18nOpts)}
        </small>
      </DialogTitle>
      <DialogDescription className={error || isPending ? "hidden" : ""}>
        <p>{t("description")}</p>
        <p>{t("sizeLimit")}</p>
      </DialogDescription>
      <DialogAlert
        className={!isPending && error ? void 0 : "hidden"}
        variant="destructive"
      >
        <CircleAlert />
        <AlertDescription>
          {t(error?.name || "Error", {
            ...CommonI18nOpts,
            defaultValue: error?.message || "",
          })}
          <Blank />
        </AlertDescription>
      </DialogAlert>
      <DialogAlert className={isPending ? void 0 : "hidden"}>
        <Spinner />
        <AlertDescription>
          {t("saving")}
          <Blank />
        </AlertDescription>
      </DialogAlert>
      <form id="memoEditForm" onSubmit={handleSubmit(handleFormSubmit)}>
        <InputGroup>
          <InputGroupTextarea
            {...register("memo", {
              maxLength: 256,
              onChange: clearErrors,
              setValueAs: setMemoAs,
            })}
            placeholder={t("placeholder")}
            className="min-h-29 max-h-29 scrollbar-thin"
            autoComplete="off"
            aria-invalid={errors.memo ? "true" : "false"}
          />
          <InputGroupAddon align="block-end">
            <InputGroupText className="text-xs text-muted-foreground">
              {memoValue?.length || 0} / 256
            </InputGroupText>
            <InputGroupButton
              variant="default"
              size="sm"
              className="ml-auto max-md:hidden"
              onClick={handleReset}
              disabled={formDisabled || !isDirty}
            >
              {t("action.reset", CommonI18nOpts)}
            </InputGroupButton>
          </InputGroupAddon>
        </InputGroup>
      </form>
      <DialogAction className="bg-muted">
        <Button
          className="max-md:hidden"
          onClick={handleClose}
          variant="outline"
        >
          {t("action.cancel", CommonI18nOpts)}
        </Button>
        <Button
          type="submit"
          form="memoEditForm"
          disabled={formDisabled || !isValid || !isDirty}
        >
          {t("action.save", CommonI18nOpts)}
        </Button>
      </DialogAction>
    </Dialog>
  );
};

const PeerPortEditKeyPrefix = `${ExtraI18nPrefix}.${ExtraType.PeerPortEdit}`;
const PeerPortEditExtra = () => {
  const { t } = useTranslation(SettingsName, {
    keyPrefix: PeerPortEditKeyPrefix,
  });

  const { type, open } = useAppExtra<ExtraState<ExtraType.PeerPortEdit>>({});
  const dispatch = useAppDispatch();

  const handleClose = () => {
    dispatch?.(withExtraAction(type, false));
  };

  const {
    data,
    isLoading,
    isSuccess: isLoadingSuccess,
  } = usePeerPortField({
    enabled: type === ExtraType.PeerPortEdit && open,
  });

  const [error, setError] = useState<Error | undefined>();
  const clearErrors = () => {
    setError(undefined);
  };

  const { mutateAsync, isPending } = usePeerPortMutation();
  const {
    reset,
    register,
    handleSubmit,
    formState: { errors, disabled: formDisabled, isValid, isDirty },
    setFocus,
  } = useForm<Pick<Settings, "peerPort">>({
    progressive: true,
    disabled:
      type !== ExtraType.PeerPortEdit ||
      !open ||
      isLoading ||
      isPending ||
      !isLoadingSuccess,
  });

  const handleReset = () => {
    reset({ peerPort: data?.peerPort });
    clearErrors();
  };

  const onSuccess = (settings: Settings) => {
    toast.add(
      withToast(ToastType.Success, {
        description: t("success"),
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
    if (type === ExtraType.PeerPortEdit && open && isLoadingSuccess) {
      setFocus("peerPort");
      handleReset();
    }
  }, [type, open, isLoadingSuccess, setFocus]);

  const peerPortValidate = (value: number) => {
    if (isNaN(value) || !Number.isInteger(value)) {
      return false;
    }
    return value >= 1 && value <= 65535;
  };

  if (type !== ExtraType.PeerPortEdit) {
    return null;
  }

  return (
    <Dialog open={open} onClose={handleClose}>
      <DialogExtra>
        <Button size="icon-sm" variant="ghost" onClick={handleClose}>
          <CloseIcon />
        </Button>
      </DialogExtra>
      <DialogTitle>
        {t("title")}
        <small className="text-muted-foreground px-1">
          {t("badge.edit", CommonI18nOpts)}
        </small>
      </DialogTitle>
      <DialogDescription className={error || isPending ? "hidden" : ""}>
        <p>{t("description")}</p>
      </DialogDescription>
      <DialogAlert
        className={!isPending && error ? void 0 : "hidden"}
        variant="destructive"
      >
        <CircleAlert />
        <AlertDescription>
          {t(error?.name || "Error", {
            ...CommonI18nOpts,
            defaultValue: error?.message || "",
          })}
        </AlertDescription>
      </DialogAlert>
      <DialogAlert className={isPending ? void 0 : "hidden"}>
        <Spinner />
        <AlertDescription>{t("saving")}</AlertDescription>
      </DialogAlert>
      <form id="peerPortEditForm" onSubmit={handleSubmit(handleFormSubmit)}>
        <InputGroup>
          <InputGroupInput
            {...register("peerPort", {
              required: true,
              valueAsNumber: true,
              min: 1,
              max: 65535,
              validate: peerPortValidate,
              onChange: clearErrors,
            })}
            type="number"
            placeholder={t("placeholder", { example: 9000 })}
            autoComplete="off"
            aria-invalid={errors.peerPort ? "true" : "false"}
          />
          <InputGroupAddon align="block-end">
            <InputGroupText className="text-xs text-muted-foreground">
              1 ~ 65535
            </InputGroupText>
            <InputGroupButton
              variant="default"
              size="sm"
              className="ml-auto max-md:hidden"
              onClick={handleReset}
              disabled={formDisabled || !isDirty}
            >
              {t("action.reset", CommonI18nOpts)}
            </InputGroupButton>
          </InputGroupAddon>
        </InputGroup>
      </form>
      <DialogAction className="bg-muted">
        <Button
          className="max-md:hidden"
          onClick={handleClose}
          variant="outline"
        >
          {t("action.cancel", CommonI18nOpts)}
        </Button>
        <Button
          type="submit"
          form="peerPortEditForm"
          disabled={formDisabled || !isValid || !isDirty}
        >
          <Spinner className={isPending ? "" : "hidden"} />
          {t("action.save", CommonI18nOpts)}
        </Button>
      </DialogAction>
    </Dialog>
  );
};

function parseBroadcastAddr(addr: string): undefined | string {
  if (/^\s+$/g.test(addr) || addr.trim().length === 0) {
    return;
  }
  addr = addr.replace(/\s+/g, "");
  const splitIdx = addr.lastIndexOf(":");
  if (splitIdx === -1) {
    return;
  }

  const portPart = addr.substring(splitIdx + 1);
  const port = Number(portPart);
  if (
    isNaN(port) ||
    !Number.isInteger(port) ||
    port < 1 ||
    port > 65535 ||
    port.toString() !== portPart
  ) {
    return;
  }

  const ipPart = newIPAddr(addr.substring(0, splitIdx));
  if (!ipPart || !ipPart.isMulticast()) {
    return;
  }

  return `${ipPart.correctForm()}:${portPart}`;
}

function countAddrsText(
  addrsText: string | undefined,
  parseFn: typeof parseBroadcastAddr,
) {
  const lines = addrsText?.split("\n") || [];
  if (lines.length === 0) {
    return [0, 0];
  }

  let validCount = 0;
  let invalidCount = 0;
  for (const line of lines) {
    if (/^\s+$/g.test(line) || line.trim().length === 0) {
      continue;
    }
    if (parseFn(line) !== void 0) {
      validCount++;
    } else {
      invalidCount++;
    }
  }

  return [validCount, invalidCount];
}

const BroadcastAddrsEditKeyPrefix = `${ExtraI18nPrefix}.${ExtraType.BroadcastAddrsEdit}`;
const BroadcastAddrsEditExtra = () => {
  const { t } = useTranslation(SettingsName, {
    keyPrefix: BroadcastAddrsEditKeyPrefix,
  });

  const { type, open } = useAppExtra<ExtraState<ExtraType.BroadcastAddrsEdit>>(
    {},
  );
  const dispatch = useAppDispatch();

  const handleClose = () => {
    dispatch?.(withExtraAction(type, false));
  };

  const {
    data,
    isLoading,
    isSuccess: isLoadingSuccess,
  } = useBroadcastAddrsField({
    enabled: type === ExtraType.BroadcastAddrsEdit && open,
  });

  const [error, setError] = useState<Error | undefined>();
  const clearErrors = () => {
    setError(undefined);
  };

  const { mutateAsync, isPending } = useBroadcastAddrsMutation();
  const {
    reset,
    watch,
    register,
    handleSubmit,
    formState: { errors, disabled: formDisabled, isValid, isDirty },
    setFocus,
  } = useForm<{ broadcastAddrsText: string }>({
    progressive: true,
    disabled:
      type !== ExtraType.BroadcastAddrsEdit ||
      !open ||
      isLoading ||
      isPending ||
      !isLoadingSuccess,
  });

  const handleReset = () => {
    reset({ broadcastAddrsText: (data?.broadcastAddrs || []).join("\n") });
    clearErrors();
  };

  const onSuccess = (settings: Settings) => {
    toast.add(
      withToast(ToastType.Success, {
        description: t("success"),
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

  const broadcastAddrText = watch("broadcastAddrsText");
  const [validCount, invalidCount] = useMemo(
    () => countAddrsText(broadcastAddrText, parseBroadcastAddr),
    [broadcastAddrText, parseBroadcastAddr],
  );

  const broadcastAddrTextValidate = (value: string) => {
    const [_, invalidCount] = countAddrsText(value, parseBroadcastAddr);
    return invalidCount === 0;
  };

  useEffect(() => {
    if (type === ExtraType.BroadcastAddrsEdit && open && isLoadingSuccess) {
      setFocus("broadcastAddrsText");
      handleReset();
    }
  }, [type, open, isLoadingSuccess, setFocus]);

  if (type !== ExtraType.BroadcastAddrsEdit) {
    return null;
  }

  return (
    <Dialog open={open} onClose={handleClose}>
      <DialogExtra>
        <Button size="icon-sm" variant="ghost" onClick={handleClose}>
          <CloseIcon />
        </Button>
      </DialogExtra>
      <DialogTitle>
        {t("title")}
        <small className="text-muted-foreground px-1">
          {t("badge.edit", CommonI18nOpts)}
        </small>
      </DialogTitle>
      <DialogDescription className={error || isPending ? "hidden" : ""}>
        <p>{t("description")}</p>
        <p>{t("dataFormat")}</p>
      </DialogDescription>
      <DialogAlert
        className={!isPending && error ? void 0 : "hidden"}
        variant="destructive"
      >
        <CircleAlert />
        <AlertDescription>
          {t(error?.name || "Error", {
            ...CommonI18nOpts,
            defaultValue: error?.message || "",
          })}
          <Blank />
        </AlertDescription>
      </DialogAlert>
      <DialogAlert className={isPending ? void 0 : "hidden"}>
        <Spinner />
        <AlertDescription>
          {t("saving")}
          <Blank />
        </AlertDescription>
      </DialogAlert>
      <form
        id="broadcastAddrsEditForm"
        onSubmit={handleSubmit(handleFormSubmit)}
      >
        <InputGroup>
          <InputGroupTextarea
            {...register("broadcastAddrsText", {
              onChange: clearErrors,
              validate: broadcastAddrTextValidate,
            })}
            placeholder={t("placeholder", { example: "224.0.0.1:9000" })}
            className="min-h-29 max-h-29 scrollbar-thin"
            autoComplete="off"
            aria-invalid={errors.broadcastAddrsText ? "true" : "false"}
          />
          <InputGroupAddon align="block-end">
            <Badge variant="secondary">
              <CircleCheckBig data-icon="inline-start" /> {validCount}
            </Badge>
            <Badge
              variant="destructive"
              className={
                errors.broadcastAddrsText && invalidCount > 0
                  ? void 0
                  : "hidden"
              }
            >
              <CircleX data-icon="inline-start" /> {invalidCount}
            </Badge>
            <InputGroupButton
              variant="default"
              size="sm"
              className="ml-auto max-md:hidden"
              onClick={handleReset}
              disabled={formDisabled || !isDirty}
            >
              {t("action.reset", CommonI18nOpts)}
            </InputGroupButton>
          </InputGroupAddon>
        </InputGroup>
      </form>
      <DialogAction className="bg-muted">
        <Button
          className="max-md:hidden"
          onClick={handleClose}
          variant="outline"
        >
          {t("action.cancel", CommonI18nOpts)}
        </Button>
        <Button
          type="submit"
          form="broadcastAddrsEditForm"
          disabled={formDisabled || !isValid || !isDirty}
        >
          {t("action.save", CommonI18nOpts)}
        </Button>
      </DialogAction>
    </Dialog>
  );
};

function parsePublicAddr(addr: string): undefined | string {
  if (/^\s+$/g.test(addr) || addr.trim().length === 0) {
    return;
  }
  addr = addr.replace(/\s+/g, "");
  const splitIdx = addr.lastIndexOf(":");
  if (splitIdx === -1) {
    return;
  }

  const portPart = addr.substring(splitIdx + 1);
  const port = Number(portPart);
  if (
    isNaN(port) ||
    !Number.isInteger(port) ||
    port < 1 ||
    port > 65535 ||
    port.toString() !== portPart
  ) {
    return;
  }

  const ipPart = addr.substring(0, splitIdx);
  const ipAddr = newIPAddr(ipPart);
  if (ipAddr) {
    if (
      ipAddr.isBroadcast() ||
      ipAddr.isMulticast() ||
      ipAddr.isLinkLocal() ||
      ipAddr.isLoopback() ||
      ipAddr.isCGNAT() ||
      ipAddr.isUnspecified()
    )
      return;
    return `${ipAddr.correctForm()}:${portPart}`;
  }

  const url = URL.parse(`https://${addr}`);
  if (!url || url.host !== addr) {
    return;
  }
  return url.host;
}

const PublicAddrsEditKeyPrefix = `${ExtraI18nPrefix}.${ExtraType.PublicAddrsEdit}`;
const PublicAddrsEditExtra = () => {
  const { t } = useTranslation(SettingsName, {
    keyPrefix: PublicAddrsEditKeyPrefix,
  });

  const { type, open } = useAppExtra<ExtraState<ExtraType.PublicAddrsEdit>>({});
  const dispatch = useAppDispatch();

  const handleClose = () => {
    dispatch?.(withExtraAction(type, false));
  };

  const {
    data,
    isLoading,
    isSuccess: isLoadingSuccess,
  } = usePublicAddrsField({
    enabled: type === ExtraType.PublicAddrsEdit && open,
  });

  const [error, setError] = useState<Error | undefined>();
  const clearErrors = () => {
    setError(undefined);
  };

  const { mutateAsync, isPending } = usePublicAddrsMutation();
  const {
    reset,
    watch,
    register,
    handleSubmit,
    formState: { errors, disabled: formDisabled, isValid, isDirty },
    setFocus,
  } = useForm<{ publicAddrsText: string }>({
    progressive: true,
    disabled:
      type !== ExtraType.PublicAddrsEdit ||
      !open ||
      isLoading ||
      isPending ||
      !isLoadingSuccess,
  });

  const handleReset = () => {
    reset({ publicAddrsText: (data?.publicAddrs || []).join("\n") });
    clearErrors();
  };

  const onSuccess = (settings: Settings) => {
    toast.add(
      withToast(ToastType.Success, {
        description: t("success"),
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

  const publicAddrText = watch("publicAddrsText");
  const [validCount, invalidCount] = useMemo(
    () => countAddrsText(publicAddrText, parsePublicAddr),
    [publicAddrText, parsePublicAddr],
  );

  const publicAddrTextValidate = (value: string) => {
    const [_, invalidCount] = countAddrsText(value, parsePublicAddr);
    return invalidCount === 0;
  };

  useEffect(() => {
    if (type === ExtraType.PublicAddrsEdit && open && isLoadingSuccess) {
      setFocus("publicAddrsText");
      handleReset();
    }
  }, [type, open, isLoadingSuccess, setFocus]);

  if (type !== ExtraType.PublicAddrsEdit) {
    return null;
  }

  return (
    <Dialog open={open} onClose={handleClose}>
      <DialogExtra>
        <Button size="icon-sm" variant="ghost" onClick={handleClose}>
          <CloseIcon />
        </Button>
      </DialogExtra>
      <DialogTitle>
        {t("title")}
        <small className="text-muted-foreground px-1">
          {t("badge.edit", CommonI18nOpts)}
        </small>
      </DialogTitle>
      <DialogDescription className={error || isPending ? "hidden" : ""}>
        <p>{t("description")}</p>
        <p>{t("dataFormat")}</p>
      </DialogDescription>
      <DialogAlert
        className={!isPending && error ? void 0 : "hidden"}
        variant="destructive"
      >
        <CircleAlert />
        <AlertDescription>
          {t(error?.name || "Error", {
            ...CommonI18nOpts,
            defaultValue: error?.message || "",
          })}
          <Blank />
        </AlertDescription>
      </DialogAlert>
      <DialogAlert className={isPending ? void 0 : "hidden"}>
        <Spinner />
        <AlertDescription>
          {t("saving")}
          <Blank />
        </AlertDescription>
      </DialogAlert>
      <form id="publicAddrsEditForm" onSubmit={handleSubmit(handleFormSubmit)}>
        <InputGroup>
          <InputGroupTextarea
            {...register("publicAddrsText", {
              onChange: clearErrors,
              validate: publicAddrTextValidate,
            })}
            placeholder={t("placeholder", { example: "www.example.com:9000" })}
            className="min-h-29 max-h-29 scrollbar-thin"
            autoComplete="off"
            aria-invalid={errors.publicAddrsText ? "true" : "false"}
          />
          <InputGroupAddon align="block-end">
            <Badge variant="secondary">
              <CircleCheckBig data-icon="inline-start" /> {validCount}
            </Badge>
            <Badge
              variant="destructive"
              className={
                errors.publicAddrsText && invalidCount > 0 ? void 0 : "hidden"
              }
            >
              <CircleX data-icon="inline-start" /> {invalidCount}
            </Badge>
            <InputGroupButton
              variant="default"
              size="sm"
              className="ml-auto max-md:hidden"
              onClick={handleReset}
              disabled={formDisabled || !isDirty}
            >
              {t("action.reset", CommonI18nOpts)}
            </InputGroupButton>
          </InputGroupAddon>
        </InputGroup>
      </form>
      <DialogAction className="bg-muted">
        <Button
          className="max-md:hidden"
          onClick={handleClose}
          variant="outline"
        >
          {t("action.cancel", CommonI18nOpts)}
        </Button>
        <Button
          type="submit"
          form="publicAddrsEditForm"
          disabled={formDisabled || !isValid || !isDirty}
        >
          {t("action.save", CommonI18nOpts)}
        </Button>
      </DialogAction>
    </Dialog>
  );
};

const WebPortEditKeyPrefix = `${ExtraI18nPrefix}.${ExtraType.WebPortEdit}`;
const WebPortEditExtra = () => {
  const { t } = useTranslation(SettingsName, {
    keyPrefix: WebPortEditKeyPrefix,
  });

  const { type, open } = useAppExtra<ExtraState<ExtraType.WebPortEdit>>({});
  const dispatch = useAppDispatch();

  const handleClose = () => {
    dispatch?.(withExtraAction(type, false));
  };

  const {
    data,
    isLoading,
    isSuccess: isLoadingSuccess,
  } = useWebPortField({
    enabled: type === ExtraType.WebPortEdit && open,
  });

  const [error, setError] = useState<Error | undefined>();
  const clearErrors = () => {
    setError(undefined);
  };

  const { mutateAsync, isPending } = useWebPortMutation();
  const {
    reset,
    register,
    handleSubmit,
    formState: { errors, disabled: formDisabled, isValid, isDirty },
    setFocus,
  } = useForm<Pick<HostSettings, "webPort">>({
    progressive: true,
    disabled:
      type !== ExtraType.WebPortEdit ||
      !open ||
      isLoading ||
      isPending ||
      !isLoadingSuccess,
  });

  const handleReset = () => {
    reset({ webPort: data?.webPort });
    clearErrors();
  };

  const onSuccess = (settings: HostSettings) => {
    toast.add(
      withToast(ToastType.Success, {
        description: t("success"),
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
    if (type === ExtraType.WebPortEdit && open && isLoadingSuccess) {
      setFocus("webPort");
      handleReset();
    }
  }, [type, open, isLoadingSuccess, setFocus]);

  const webPortValidate = (value: number) => {
    if (isNaN(value) || !Number.isInteger(value)) {
      return false;
    }
    return value >= 1 && value <= 65535;
  };

  if (type !== ExtraType.WebPortEdit) {
    return null;
  }

  return (
    <Dialog open={open} onClose={handleClose}>
      <DialogExtra>
        <Button size="icon-sm" variant="ghost" onClick={handleClose}>
          <CloseIcon />
        </Button>
      </DialogExtra>
      <DialogTitle>
        {t("title")}
        <small className="text-muted-foreground px-1">
          {t("badge.edit", CommonI18nOpts)}
        </small>
      </DialogTitle>
      <DialogDescription className={error || isPending ? "hidden" : ""}>
        <p>{t("description")}</p>
      </DialogDescription>
      <DialogAlert
        className={!isPending && error ? void 0 : "hidden"}
        variant="destructive"
      >
        <CircleAlert />
        <AlertDescription>
          {t(error?.name || "Error", {
            ...CommonI18nOpts,
            defaultValue: error?.message || "",
          })}
        </AlertDescription>
      </DialogAlert>
      <DialogAlert className={isPending ? void 0 : "hidden"}>
        <Spinner />
        <AlertDescription>{t("saving")}</AlertDescription>
      </DialogAlert>
      <form id="webPortEditForm" onSubmit={handleSubmit(handleFormSubmit)}>
        <InputGroup>
          <InputGroupInput
            {...register("webPort", {
              required: true,
              valueAsNumber: true,
              min: 1,
              max: 65535,
              validate: webPortValidate,
              onChange: clearErrors,
            })}
            type="number"
            placeholder={t("placeholder", { example: 9000 })}
            autoComplete="off"
            aria-invalid={errors.webPort ? "true" : "false"}
          />
          <InputGroupAddon align="block-end">
            <InputGroupText className="text-xs text-muted-foreground">
              1 ~ 65535
            </InputGroupText>
            <InputGroupButton
              variant="default"
              size="sm"
              className="ml-auto max-md:hidden"
              onClick={handleReset}
              disabled={formDisabled || !isDirty}
            >
              {t("action.reset", CommonI18nOpts)}
            </InputGroupButton>
          </InputGroupAddon>
        </InputGroup>
      </form>
      <DialogAction className="bg-muted">
        <Button
          className="max-md:hidden"
          onClick={handleClose}
          variant="outline"
        >
          {t("action.cancel", CommonI18nOpts)}
        </Button>
        <Button
          type="submit"
          form="webPortEditForm"
          disabled={formDisabled || !isValid || !isDirty}
        >
          <Spinner className={isPending ? "" : "hidden"} />
          {t("action.save", CommonI18nOpts)}
        </Button>
      </DialogAction>
    </Dialog>
  );
};
