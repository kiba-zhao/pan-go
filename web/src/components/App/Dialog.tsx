import { Overlay } from "./Layout";
import { CircleAlert, CloseIcon } from "./Icon";
import {
  UnknownI18nKey,
  useTranslation,
  I18nVariant,
  useAppI18n,
  withAppError,
} from "./I18Next";
import { withAppExtraAction, useAppDispatch, useAppExtra } from "./Context";

import {
  type ComponentProps,
  type MouseEvent,
  type KeyboardEvent,
  useRef,
  type RefObject,
  type PropsWithChildren,
  type ReactNode,
} from "react";
import { cn } from "@/lib/utils";
import {
  Alert,
  AlertTitle,
  AlertDescription,
  AlertAction,
} from "@/components/ui/alert";
export { AlertTitle, AlertDescription, AlertAction };
import { Spinner } from "@/components/ui/spinner";
import { Button } from "@/components/ui/button";
import {
  Card,
  CardHeader,
  CardAction,
  CardTitle,
  CardDescription,
  CardContent,
  CardFooter,
} from "@/components/ui/card";

export enum DialogVariant {
  Default = "default",
  Modal = "modal",
}
const DialogVariants = {
  [DialogVariant.Default]:
    "relative w-full md:max-w-md max-md:inset-x-0 max-md:bottom-0 max-md:absolute",
  [DialogVariant.Modal]: void 0,
};
const DialogOverlayVariants = {
  [DialogVariant.Default]: "md:flex md:justify-center md:items-center",
  [DialogVariant.Modal]: "md:flex md:justify-center md:items-start pt-16",
};
type DialogProps = Omit<ComponentProps<"dialog">, "ref"> & {
  variant?: DialogVariant;
  ref?: RefObject<HTMLDialogElement>;
};
export const Dialog = ({
  children,
  className,
  variant = DialogVariant.Default,
  onClose,
  onClick,
  ref,
  open,
  ...props
}: DialogProps) => {
  const dialogRef = ref || useRef<HTMLDialogElement>(null);

  const handleClick = (event: MouseEvent<HTMLDialogElement>) => {
    event.stopPropagation();
    onClick?.(event);
  };
  const handleOverlayClick = () => {
    dialogRef.current?.close();
  };

  const handleKeyDown = (event: KeyboardEvent<HTMLDivElement>) => {
    event.stopPropagation();
    if (event.key === "Escape") {
      handleOverlayClick();
    }
  };

  return (
    <Overlay
      className={cn(
        DialogOverlayVariants[variant] ||
          DialogOverlayVariants[DialogVariant.Default],
        !open && "hidden!",
      )}
      onKeyDown={handleKeyDown}
      onClick={handleOverlayClick}
    >
      <dialog
        ref={dialogRef}
        {...props}
        className={cn(
          // "group/dialog relative bg-sidebar text-sidebar-foreground border-border md:border-1 max-md:border-t-1",
          // "shadow-ring shadow-xl/30",
          "bg-transparent",
          DialogVariants[variant] || DialogVariants[DialogVariant.Default],
          className,
        )}
        onClose={onClose}
        onClick={handleClick}
        open={open}
      >
        {children}
      </dialog>
    </Overlay>
  );
};

export type DialogExtraState = {
  open?: boolean;
};
export function withDialogExtraState<State extends DialogExtraState>(
  state: State,
) {
  const { open, ...rest } = state;
  return withAppExtraAction({ open: open !== false, ...rest });
}

export const DialogExtra = ({
  children,
  ...props
}: Omit<DialogProps, "open" | "onClose">) => {
  const extraState = useAppExtra<DialogExtraState>();
  const dispatch = useAppDispatch();
  const handleClose = () => {
    dispatch?.(withDialogExtraState({ ...extraState, open: false }));
  };
  return (
    <Dialog {...props} open={extraState.open} onClose={handleClose}>
      {children}
    </Dialog>
  );
};

export const DialogExtraCloseAction = ({
  children,
  ...props
}: Omit<ComponentProps<typeof Button>, "onClick">) => {
  const extraState = useAppExtra<DialogExtraState>();
  const dispatch = useAppDispatch();
  const handleClose = () => {
    dispatch?.(withDialogExtraState({ ...extraState, open: false }));
  };
  return (
    <Button {...props} onClick={handleClose}>
      {children}
    </Button>
  );
};

type DialogExtraActionProps = {
  progress?: boolean;
  cancelText?: string;
} & ComponentProps<typeof Button>;
export const DialogExtraAction = ({
  progress,
  cancelText,
  disabled,
  children,
  ...props
}: DialogExtraActionProps) => {
  const { namespace } = useAppI18n();
  const { t } = useTranslation(namespace);

  return (
    <>
      <DialogExtraCloseAction variant="outline">
        {cancelText || t(`${I18nVariant.Action}.cancel`)}
      </DialogExtraCloseAction>
      <Button {...props} disabled={disabled || progress}>
        <Spinner className={progress ? "" : "hidden"} />
        {children || t(`${I18nVariant.Action}.submit`)}
      </Button>
    </>
  );
};

export const DialogExtraClose = ({
  size = "icon-sm",
  variant = "ghost",
  children,
  ...props
}: Omit<ComponentProps<typeof Button>, "onClick">) => {
  return (
    <DialogExtraCloseAction {...props} size={size} variant={variant}>
      {children || <CloseIcon />}
    </DialogExtraCloseAction>
  );
};

export const DialogExtraTitle = ({
  text,
  smallText,
  className = "text-muted-foreground",
  children,
  ...props
}: ComponentProps<"small"> & { text?: string; smallText?: string }) => {
  return (
    <>
      {text}
      <small {...props} className={cn("px-1", className)}>
        {smallText || children}
      </small>
    </>
  );
};

export const DialogExtraAlert = ({
  icon,
  className,
  children,
  ...props
}: ComponentProps<typeof Alert> & { icon?: ReactNode }) => (
  <Alert
    {...props}
    className={cn("border-none py-0 px-0 gap-0 bg-transparent", className)}
  >
    {icon}
    <AlertDescription>{children}</AlertDescription>
  </Alert>
);

type DialogExtraSpinnerAlertProps = {
  error?: Error;
  errorContent?: ReactNode;
  errorVariant?: ComponentProps<typeof Alert>["variant"];
  rotate?: boolean;
  rotateContent?: ReactNode;
  filledContent?: ReactNode;
  rotateVariant?: ComponentProps<typeof Alert>["variant"];
} & Omit<ComponentProps<typeof Alert>, "variant">;
export const DialogExtraSpinnerAlert = ({
  error,
  errorContent,
  errorVariant = "destructive",
  rotate,
  rotateContent,
  rotateVariant,
  children,
  ...props
}: DialogExtraSpinnerAlertProps) => {
  const { namespace } = useAppI18n();
  const { t } = useTranslation(namespace);

  let icon: ReactNode;
  let children_: ReactNode;
  let variant: ComponentProps<typeof Alert>["variant"];
  if (rotate) {
    icon = <Spinner />;
    children_ = rotateContent;
    variant = rotateVariant;
  } else if (error) {
    icon = <CircleAlert />;
    children_ = errorContent || t(...withAppError(error));
    variant = errorVariant;
  } else {
    return children;
  }

  return (
    <DialogExtraAlert {...props} variant={variant} icon={icon}>
      {children_}
    </DialogExtraAlert>
  );
};

type CardDialogExtraProps = PropsWithChildren<{
  title?: PropsWithChildren["children"];
  description?: PropsWithChildren["children"];
  action?: PropsWithChildren["children"];
  footer?: PropsWithChildren["children"];
  footerClassName?: ComponentProps<typeof CardFooter>["className"];
  cardClassName?: ComponentProps<typeof Card>["className"];
}>;
export const CardDialogExtra = ({
  children,
  title,
  description,
  action,
  footer,
  footerClassName,
  cardClassName,
}: CardDialogExtraProps) => {
  return (
    <DialogExtra>
      <Card className={cardClassName}>
        <CardHeader
          className={cn(title || description || action ? void 0 : "hidden")}
        >
          <CardAction className={cn(action ? void 0 : "hidden")}>
            {action}
          </CardAction>
          <CardTitle className={cn(title ? void 0 : "hidden")}>
            {title}
          </CardTitle>
          <CardDescription className={cn(description ? void 0 : "hidden")}>
            {description}
          </CardDescription>
        </CardHeader>
        <CardContent className={cn(children ? void 0 : "hidden")}>
          {children}
        </CardContent>
        <CardFooter className={cn(footer ? void 0 : "hidden", footerClassName)}>
          {footer}
        </CardFooter>
      </Card>
    </DialogExtra>
  );
};
