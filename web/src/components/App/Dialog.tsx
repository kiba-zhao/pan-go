import {
  type ComponentProps,
  type MouseEvent,
  type KeyboardEvent,
  useRef,
} from "react";
import { cn } from "@/lib/utils";
import { Overlay } from "./Layout";
import { Alert } from "@/components/ui/alert";
export {
  AlertTitle,
  AlertDescription,
  AlertAction,
} from "@/components/ui/alert";

export enum DialogVariant {
  Default = "default",
}
const DialogVariants = {
  [DialogVariant.Default]:
    "flex flex-col gap-4 p-4 md:rounded-xl max-md:w-full max-md:inset-x-0 max-md:bottom-0 max-md:absolute",
};
const DialogOverlayVariants = {
  [DialogVariant.Default]: "md:flex md:items-center md:justify-center",
};
export const Dialog = ({
  children,
  className,
  variant = DialogVariant.Default,
  onClose,
  onClick,
  open,
  ...props
}: ComponentProps<"dialog"> & { variant?: DialogVariant }) => {
  const ref = useRef<HTMLDialogElement>(null);
  const handleClick = (event: MouseEvent<HTMLDialogElement>) => {
    event.stopPropagation();
    onClick?.(event);
  };

  const handleKeyDown = (event: KeyboardEvent<HTMLDivElement>) => {
    event.stopPropagation();
    if (event.key === "Escape") {
      ref.current?.close();
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
      onClick={() => ref.current?.close()}
    >
      <dialog
        ref={ref}
        {...props}
        className={cn(
          "group/dialog relative bg-sidebar text-sidebar-foreground border-border md:border-1 max-md:border-t-1",
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
        "-mx-4 -mb-4 p-4 border-t flex flex-col gap-2 md:rounded-b-xl md:flex-row md:justify-end",
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
}: ComponentProps<"p">) => {
  return (
    <p {...props} className={cn("text-sm text-muted-foreground", className)}>
      {children}
    </p>
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
