import { ReactNode, type ComponentProps, type ElementType } from "react";
import { cn } from "../../lib/utils";

type LayoutComponentType = ElementType<{
  className?: string;
  children?: ReactNode;
}>;
type LayoutProps<T extends LayoutComponentType> = ComponentProps<T> & {
  BaseComponent?: T;
};
export const ExclusiveLayout = <T extends ElementType>({
  className,
  children,
  BaseComponent = "div",
  ...props
}: LayoutProps<T>) => {
  return (
    <BaseComponent
      {...props}
      className={cn("h-full flex items-center justify-center", className)}
    >
      {children}
    </BaseComponent>
  );
};

export const PageInfoSection = ({
  logo,
  title,
  description,
}: {
  logo?: ReactNode;
  title?: string;
  description?: string;
}) => {
  return (
    <ExclusiveLayout
      BaseComponent="section"
      className="flex flex-col items-center justify-center gap-1"
    >
      {logo}
      {title && <h1 className="text-xl font-bold mt-4">{title}</h1>}
      {description && (
        <small className="text-sm font-medium text-muted-foreground">
          {description}
        </small>
      )}
    </ExclusiveLayout>
  );
};

export const Overlay = ({
  children,
  className,
  ...props
}: ComponentProps<"div">) => {
  return (
    <div
      {...props}
      className={cn(
        "z-50 bg-black/10 backdrop-blur-xs fixed inset-0",
        className,
      )}
    >
      {children}
    </div>
  );
};
