import { useAppExtra, withAppExtraAction } from "./Context";
import { useTranslation, I18nVariant, useAppI18n } from "./I18Next";
import type { ComponentType, ComponentProps } from "react";
import { DialogTitle, DialogExtra } from "./Dialog";
import { Button } from "@/components/ui/button";
import { CloseIcon } from "./Icon";
import { cn } from "@/lib/utils";
import { Spinner } from "@/components/ui/spinner";

export type ExtraState<Type, State extends any> = {
  type?: Type;
} & State;

type ExtraProps<Type, Props> = {
  asProps?: Props;
  as: ComponentType<Props>;
  extraType: Type;
};
export const Extra = <
  Type,
  State,
  Props extends { extraState: ExtraState<Type, State> },
>({
  as,
  asProps: props,
  extraType,
}: ExtraProps<Type, Props>) => {
  const extraState = useAppExtra<ExtraState<Type, State>>();
  const Component = as;
  if (extraState.type !== extraType) {
    return null;
  }
  return <Component {...(props || ({} as Props))} extraState={extraState} />;
};

export type DialogExtraState = { open?: boolean };
export function withDialogExtraState<
  Type,
  State extends DialogExtraState = DialogExtraState,
>(state: ExtraState<Type, State>) {
  const { type, open, ...rest } = state;
  return withAppExtraAction({ type, open: open !== false, ...rest });
}

export const DialogExtraTitle = ({
  text,
  className = "text-muted-foreground",
  children,
  ...props
}: ComponentProps<"small"> & { text?: string }) => {
  return (
    <DialogTitle>
      {text}
      <small {...props} className={cn("px-1", className)}>
        {children}
      </small>
    </DialogTitle>
  );
};

export const DialogExtraAction = ({
  onClose,
  progress,
  disabled,
  children,
  ...props
}: {
  onClose?: ComponentProps<"button">["onClick"];
  progress?: boolean;
} & ComponentProps<typeof Button>) => {
  const { namespace } = useAppI18n();
  const { t } = useTranslation(namespace);

  return (
    <>
      <Button variant="outline" onClick={onClose}>
        {t(`${I18nVariant.Action}.cancel`)}
      </Button>
      <Button {...props} disabled={disabled || progress}>
        <Spinner className={progress ? "" : "hidden"} />
        {children}
      </Button>
    </>
  );
};

export const DialogExtraClose = ({
  onClose,
}: {
  onClose?: ComponentProps<typeof Button>["onClick"];
}) => {
  return (
    <DialogExtra>
      <Button size="icon-sm" variant="ghost" onClick={onClose}>
        <CloseIcon />
      </Button>
    </DialogExtra>
  );
};
