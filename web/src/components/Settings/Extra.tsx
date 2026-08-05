import { useAppExtra, useAppDispatch } from "@/components/App/Context";
import {
  Dialog,
  DialogTitle,
  DialogDescription,
  DialogAction,
  DialogExtra,
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
import { Button } from "@/components/ui/button";
import { Badge } from "@/components/ui/badge";
import { ExtraState as AppExtraState } from "@/components/App/Extra";
import { useRef, useEffect, useState } from "react";
import { useIMask } from "react-imask";

import { useTranslation, DefaultNS } from "@/components/App/I18Next";
import { SettingsName, ExtraI18nPrefix } from "./meta";
import { useNameField, useNameMutation, type Settings } from "./ReactQuery";
import { useForm } from "@/components/App/Form";
import { Spinner } from "@/components/ui/spinner";
import { DialogAlert, AlertDescription } from "@/components/App/Dialog";
import {
  toast,
  withToast,
  ToastType,
  type ToastID,
} from "@/components/App/Toast";

const SettingsExtra = () => (
  <>
    <MemoEditExtra />
    <NameEditExtra />
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

export const switchExtra = (
  type?: ExtraType,
  open?: boolean,
): undefined | ExtraState<ExtraType> => {
  return type === void 0 ? type : { type, open: open !== false };
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
    dispatch?.({
      extra: switchExtra(type, false),
    });
  };

  const {
    data,
    isLoading,
    isSuccess: isLoadingSuccess,
  } = useNameField({
    enabled: type === ExtraType.NameEdit && open,
  });

  const [error, setError] = useState<Error | undefined>();
  const { mutateAsync, isPending } = useNameMutation();

  const clearErrors = () => {
    setError(undefined);
  };

  const handleReset = () => {
    reset({ name: data?.name || "" });
    clearErrors();
  };

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

  const toastRef = useRef<ToastID | undefined>(undefined);
  const clearToast = () => {
    if (!toastRef.current) {
      return;
    }
    toastRef.current = undefined;
  };
  const onSuccess = (settings: Settings) => {
    toastRef.current = toast.add(
      withToast(ToastType.Success, {
        description: t("success"),
        onRemove: clearToast,
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
    clearToast();
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
      <DialogTitle className="md:w-md max-md:w-full">
        {t("title")}
        <small className="text-muted-foreground px-1">
          {t("badge.edit", CommonI18nOpts)}
        </small>
      </DialogTitle>
      <DialogDescription
        className={error || isPending ? "hidden" : ""}
      >{`${t("description")}`}</DialogDescription>
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

const MemoEditExtra = () => {
  const { type, open } = useAppExtra<ExtraState<ExtraType.MemoEdit>>({});
  const dispatch = useAppDispatch();
  const inputRef = useRef<HTMLTextAreaElement>(null);

  const handleClose = () => {
    dispatch?.({
      extra: switchExtra(type, false),
    });
  };

  useEffect(() => {
    if (type === ExtraType.MemoEdit && open) {
      inputRef.current?.focus();
    }
  }, [type, open, inputRef]);

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
        我的备忘
        <small className="text-muted-foreground">&nbsp;设置</small>
      </DialogTitle>
      <DialogDescription>
        我的备忘用于记录设备相关的信息,支持最多255个字符。
      </DialogDescription>
      <form>
        <InputGroup>
          <InputGroupTextarea
            placeholder="请输入备忘"
            name="memo"
            className="min-h-29 max-h-29 scrollbar-thin"
            maxLength={255}
            required={true}
            ref={inputRef}
          />
          <InputGroupAddon align="block-end">
            <InputGroupText className="text-xs text-muted-foreground">
              2 / 255
            </InputGroupText>
            <InputGroupButton
              variant="default"
              size="sm"
              className="ml-auto max-sm:hidden"
            >
              重置
            </InputGroupButton>
          </InputGroupAddon>
        </InputGroup>
      </form>
      <DialogAction className="bg-muted">
        <Button
          className="max-sm:hidden"
          onClick={handleClose}
          variant="outline"
        >
          取消
        </Button>
        <Button disabled>保存</Button>
      </DialogAction>
    </Dialog>
  );
};

const PeerPortEditExtra = () => {
  const { type, open } = useAppExtra<ExtraState<ExtraType.PeerPortEdit>>({});
  const dispatch = useAppDispatch();

  const { ref: inputRef } = useIMask<HTMLInputElement>({
    mask: Number,
    min: 0,
    max: 65535,
    scale: 0,
    autofix: true,
  });

  const handleClose = () => {
    dispatch?.({
      extra: switchExtra(type, false),
    });
  };

  useEffect(() => {
    if (type === ExtraType.PeerPortEdit && open) {
      inputRef.current?.focus();
    }
  }, [type, open, inputRef]);

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
        网络端口
        <small className="text-muted-foreground">&nbsp;设置</small>
      </DialogTitle>
      <DialogDescription>
        网络端口用于设备之间的通信,请输入一个可用的端口号。
      </DialogDescription>
      <form>
        <InputGroup>
          <InputGroupInput
            type="text"
            name="peerPort"
            placeholder="请输入网络端口"
            required={true}
            ref={inputRef}
            defaultValue="9000"
          />
          <InputGroupAddon align="block-end">
            <InputGroupText className="text-xs text-muted-foreground">
              0 ~ 65535
            </InputGroupText>
            <InputGroupButton
              variant="default"
              size="sm"
              className="ml-auto max-sm:hidden"
            >
              重置
            </InputGroupButton>
          </InputGroupAddon>
        </InputGroup>
      </form>
      <DialogAction className="bg-muted">
        <Button
          className="max-sm:hidden"
          onClick={handleClose}
          variant="outline"
        >
          取消
        </Button>
        <Button disabled>保存</Button>
      </DialogAction>
    </Dialog>
  );
};

const BroadcastAddrsEditExtra = () => {
  const { type, open } = useAppExtra<ExtraState<ExtraType.BroadcastAddrsEdit>>(
    {},
  );
  const dispatch = useAppDispatch();
  const inputRef = useRef<HTMLTextAreaElement>(null);

  const handleClose = () => {
    dispatch?.({
      extra: switchExtra(type, false),
    });
  };

  useEffect(() => {
    if (type === ExtraType.BroadcastAddrsEdit && open) {
      inputRef.current?.focus();
    }
  }, [type, open, inputRef]);

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
        广播地址
        <small className="text-muted-foreground">&nbsp;设置</small>
      </DialogTitle>
      <div>
        <DialogDescription>广播地址用于访问设备的网络服务。</DialogDescription>
        <DialogDescription>
          支持域名/主机名/IP地址 + 端口，每行一个地址。
        </DialogDescription>
      </div>
      <form>
        <InputGroup>
          <InputGroupTextarea
            placeholder="224.0.0.1:9001"
            name="broadcastAddrs"
            className="min-h-29 max-h-29 scrollbar-thin"
            maxLength={255}
            required={true}
            ref={inputRef}
          />
          <InputGroupAddon align="block-end">
            <Badge variant="secondary">
              <CircleCheckBig data-icon="inline-start" /> 5
            </Badge>
            <Badge variant="destructive">
              <CircleX data-icon="inline-start" /> 3
            </Badge>
            <InputGroupButton
              variant="default"
              size="sm"
              className="ml-auto max-sm:hidden"
            >
              重置
            </InputGroupButton>
          </InputGroupAddon>
        </InputGroup>
      </form>
      <DialogAction className="bg-muted">
        <Button
          className="max-sm:hidden"
          onClick={handleClose}
          variant="outline"
        >
          取消
        </Button>
        <Button disabled>保存</Button>
      </DialogAction>
    </Dialog>
  );
};

const PublicAddrsEditExtra = () => {
  const { type, open } = useAppExtra<ExtraState<ExtraType.PublicAddrsEdit>>({});
  const dispatch = useAppDispatch();
  const inputRef = useRef<HTMLTextAreaElement>(null);

  const handleClose = () => {
    dispatch?.({
      extra: switchExtra(type, false),
    });
  };

  useEffect(() => {
    if (type === ExtraType.PublicAddrsEdit && open) {
      inputRef.current?.focus();
    }
  }, [type, open, inputRef]);

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
        公共地址
        <small className="text-muted-foreground">&nbsp;设置</small>
      </DialogTitle>
      <div>
        <DialogDescription>公共地址用于访问设备的网络服务。</DialogDescription>
        <DialogDescription>
          支持域名/主机名/IP地址 + 端口，每行一个地址。
        </DialogDescription>
      </div>
      <form>
        <InputGroup>
          <InputGroupTextarea
            placeholder="www.example.com:9000"
            name="publicAddrs"
            className="min-h-29 max-h-29 scrollbar-thin"
            maxLength={255}
            required={true}
            ref={inputRef}
          />
          <InputGroupAddon align="block-end">
            <Badge variant="secondary">
              <CircleCheckBig data-icon="inline-start" /> 5
            </Badge>
            <Badge variant="destructive">
              <CircleX data-icon="inline-start" /> 3
            </Badge>
            <InputGroupButton
              variant="default"
              size="sm"
              className="ml-auto max-sm:hidden"
            >
              重置
            </InputGroupButton>
          </InputGroupAddon>
        </InputGroup>
      </form>
      <DialogAction className="bg-muted">
        <Button
          className="max-sm:hidden"
          onClick={handleClose}
          variant="outline"
        >
          取消
        </Button>
        <Button disabled>保存</Button>
      </DialogAction>
    </Dialog>
  );
};

const WebPortEditExtra = () => {
  const { type, open } = useAppExtra<ExtraState<ExtraType.WebPortEdit>>({});
  const dispatch = useAppDispatch();

  const { ref: inputRef } = useIMask<HTMLInputElement>({
    mask: Number,
    min: 0,
    max: 65535,
    scale: 0,
    autofix: true,
  });

  const handleClose = () => {
    dispatch?.({
      extra: switchExtra(type, false),
    });
  };

  useEffect(() => {
    if (type === ExtraType.WebPortEdit && open) {
      inputRef.current?.focus();
    }
  }, [type, open, inputRef]);

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
        Web端口
        <small className="text-muted-foreground">&nbsp;设置</small>
      </DialogTitle>
      <DialogDescription>
        提供Web访问的端口号,请输入一个可用的端口。
      </DialogDescription>
      <form>
        <InputGroup>
          <InputGroupInput
            type="text"
            name="webPort"
            placeholder="请输入网页服务端口"
            required={true}
            ref={inputRef}
            defaultValue="9000"
          />
          <InputGroupAddon align="block-end">
            <InputGroupText className="text-xs text-muted-foreground">
              0 ~ 65535
            </InputGroupText>
            <InputGroupButton
              variant="default"
              size="sm"
              className="ml-auto max-sm:hidden"
            >
              重置
            </InputGroupButton>
          </InputGroupAddon>
        </InputGroup>
      </form>
      <DialogAction className="bg-muted">
        <Button
          className="max-sm:hidden"
          onClick={handleClose}
          variant="outline"
        >
          取消
        </Button>
        <Button disabled>保存</Button>
      </DialogAction>
    </Dialog>
  );
};
