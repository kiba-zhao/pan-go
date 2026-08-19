import {
  type ComponentProps,
  type MouseEvent,
  type KeyboardEvent,
  useRef,
  type RefObject,
} from "react";
import { cn } from "@/lib/utils";
import { Overlay } from "./Layout";
import { CircleAlert } from "./Icon";
import {
  UnknownI18nKey,
  useTranslation,
  I18nVariant,
  useAppI18n,
} from "./I18Next";
import {
  Alert,
  AlertTitle,
  AlertDescription,
  AlertAction,
} from "@/components/ui/alert";
export { AlertTitle, AlertDescription, AlertAction };
import { Spinner } from "@/components/ui/spinner";

export enum DialogVariant {
  Default = "default",
  Modal = "modal",
}
const DialogVariants = {
  [DialogVariant.Default]:
    "flex flex-col gap-4 p-4 max-w-md w-full md:rounded-lg  max-md:inset-x-0 max-md:bottom-0 max-md:absolute",
  [DialogVariant.Modal]: void 0,
};
const DialogOverlayVariants = {
  [DialogVariant.Default]: "md:flex md:items-center md:justify-center",
  [DialogVariant.Modal]: "md:flex md:items-start md:justify-center pt-16",
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
          "group/dialog relative bg-sidebar text-sidebar-foreground border-border md:border-1 max-md:border-t-1",
          "shadow-ring shadow-xl/30",
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

export const DialogAction = ({
  children,
  className,
  ...props
}: ComponentProps<"div">) => {
  return (
    <div
      {...props}
      className={cn(
        "-mx-4 -mb-4 p-4 border-t flex flex-col gap-2 md:rounded-b-lg md:flex-row md:justify-end",
        className,
      )}
    >
      {children}
    </div>
  );
};

export const DialogTitle = ({
  className,
  children,
  ...props
}: ComponentProps<"h3">) => {
  return (
    <h3
      {...props}
      className={cn("text-base leading-none font-medium", className)}
    >
      {children}
    </h3>
  );
};

export const DialogDescription = ({
  className,
  children,
  ...props
}: ComponentProps<"div">) => {
  return (
    <div {...props} className={cn("text-sm text-muted-foreground", className)}>
      {children}
    </div>
  );
};

export enum DialogExtraVariant {
  Default = "default",
}
const DialogExtraVariants = {
  [DialogVariant.Default]: "absolute top-2 right-2",
};
export const DialogExtra = ({
  className,
  children,
  variant = DialogExtraVariant.Default,
  ...props
}: ComponentProps<"div"> & { variant?: DialogExtraVariant }) => {
  return (
    <div {...props} className={cn(DialogExtraVariants[variant], className)}>
      {children}
    </div>
  );
};

export const DialogAlert = ({
  className,
  children,
  ...props
}: ComponentProps<typeof Alert>) => (
  <Alert
    {...props}
    className={cn("border-none py-0 px-0 gap-0 bg-transparent", className)}
  >
    {children}
  </Alert>
);

export const DialogErrorAlert = ({
  children,
  error,
  ...props
}: ComponentProps<typeof DialogAlert> & {
  error?: Error;
}) => {
  const { namespace } = useAppI18n();
  const { t } = useTranslation(namespace);

  return (
    <DialogAlert variant="destructive" {...props}>
      <CircleAlert />
      <AlertDescription>
        {t(`${I18nVariant.Error}.${error?.name || UnknownI18nKey}`, {
          defaultValue: error?.message || "",
        })}
        {children}
      </AlertDescription>
    </DialogAlert>
  );
};

export const DialogProgressAlert = ({
  children,
  ...props
}: ComponentProps<typeof DialogAlert> & {
  error?: Error;
  namespace?: string;
}) => (
  <DialogAlert {...props}>
    <Spinner />
    <AlertDescription>{children}</AlertDescription>
  </DialogAlert>
);
