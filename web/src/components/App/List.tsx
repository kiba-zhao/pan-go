import { Switch } from "@/components/ui/switch";
import { cn } from "@/lib/utils";
import { type ComponentProps } from "react";
import { ChevronRight } from "./Icon";
import { type AsProps } from "./Component";
import { Link } from "./Route";

export const List = ({
  children,
  className = "text-base",
  role = "list",
  ...props
}: ComponentProps<"ul">) => {
  return (
    <ul {...props} role={role} className={className}>
      {children}
    </ul>
  );
};

export const ListOutside = ({
  className,
  children,
  ...props
}: ComponentProps<"div">) => {
  return (
    <div
      {...props}
      className={cn(
        "flex items-center py-1 px-3 group/list-outside",
        className,
      )}
    >
      {children}
    </div>
  );
};

export const ListTitle = ({
  className,
  children,
  ...props
}: ComponentProps<"h6">) => {
  return (
    <h6
      {...props}
      className={cn(
        "text-sm py-1 px-4 text-muted-foreground group-[*]/list-outside:px-0 group-[*]/list-outside:py-0 group-[*]/list-outside:grow",
        className,
      )}
    >
      {children}
    </h6>
  );
};

const ListItemBaseClassName =
  "flex gap-3 px-4 py-2.5 has-[.listitem-content]:py-0 items-center overflow-hidden";
export enum ListItemVariant {
  Primary = "primary",
  Secondary = "secondary",
  Accent = "accent",
  Muted = "muted",
}
const ListItemActiveVariants = {
  [ListItemVariant.Primary]: "bg-primary text-primary-foreground",
  [ListItemVariant.Secondary]: "bg-secondary) text-secondary-foreground",
  [ListItemVariant.Accent]: "bg-accent text-accent-foreground",
  [ListItemVariant.Muted]: "bg-muted text-muted-foreground",
};
const ListItemHoverVariants = {
  [ListItemVariant.Primary]: "hover:bg-primary",
  [ListItemVariant.Secondary]: "hover:bg-secondary",
  [ListItemVariant.Accent]: "hover:bg-accent",
  [ListItemVariant.Muted]: "hover:bg-muted",
};
export const ListItem = ({
  children,
  className,
  disabled,
  active,
  hover = ListItemVariant.Muted,
  ...props
}: ComponentProps<"li"> & {
  disabled?: boolean;
  active?: ListItemVariant;
  hover?: ListItemVariant;
}) => {
  return (
    <li
      role="listitem"
      {...props}
      className={cn(
        "group/listitem odd:z-5 relative",
        "has-[.listitem-content]:-m-px has-[.listitem-content]:py-0 has-[.listitem-content]:py-0 has-[.listitem-main]:py-0 has-[.listitem-main]:px-0 has-[.listitem-main]:block",
        active
          ? `active ${ListItemActiveVariants[active]} active-${active}`
          : disabled
            ? "disabled hover:cursor-not-allowed"
            : `${ListItemHoverVariants[hover]} has-[.listitem-main]:cursor-pointer has-[label]:cursor-pointer hover-${hover}`,
        ListItemBaseClassName,
        className,
      )}
    >
      {children}
    </li>
  );
};

export const ListItemMainClassName = `listitem-main w-full text-left cursor-[inherit] ${ListItemBaseClassName}`;
export const ListItemButton = ({
  children,
  className,
  role = "button",
  type = "button",
  ...props
}: ComponentProps<"button">) => {
  return (
    <button
      {...props}
      role={role}
      type={type}
      className={cn(ListItemMainClassName, className)}
    >
      {children}
    </button>
  );
};

export const ListItemLink = ({
  children,
  className,
  ...props
}: ComponentProps<"a">) => {
  return (
    <a {...props} className={cn(ListItemMainClassName, className)}>
      {children}
    </a>
  );
};

export const ListItemNavLink = ({
  children,
  className,
  ...props
}: ComponentProps<typeof Link>) => {
  return (
    <Link {...props} className={cn(ListItemMainClassName, className)}>
      {children}
    </Link>
  );
};

export const ListItemAvatar = <Props extends { className?: string }>({
  as,
  className,
  ...props
}: AsProps<Props>) => {
  const Component = as;
  return <Component {...props} className={cn("size-7", className)} />;
};

const ListItemContentVariant = {
  [ListItemVariant.Primary]:
    "group-[.hover-primary]/listitem:hover:border-primary group-[.active-primary]/listitem:border-primary",
  [ListItemVariant.Secondary]:
    "group-[.hover-secondary]/listitem:hover:border-secondary group-[.active-secondary]/listitem:border-secondary",
  [ListItemVariant.Accent]:
    "group-[.hover-accent]/listitem:hover:border-accent group-[.active-accent]/listitem:border-accent",
  [ListItemVariant.Muted]:
    "group-[.hover-muted]/listitem:hover:border-muted group-[.active-muted]/listitem:border-muted",
};
export const ListItemContent = ({
  children,
  className,
  ...props
}: ComponentProps<"div"> & { disabled?: boolean }) => {
  return (
    <div
      {...props}
      className={cn(
        "listitem-content flex grow py-2.5 items-center gap-0.5 overflow-hidden",
        "group-even/listitem:border-border group-odd/listitem:border-transparent border-t border-b group-last/listitem:border-b-transparent!",
        "group-[.hover-*]/listitem:event:hover:border-transparent group-[.active]/listitem:border-transparent",
        ListItemContentVariant[ListItemVariant.Primary],
        ListItemContentVariant[ListItemVariant.Secondary],
        ListItemContentVariant[ListItemVariant.Accent],
        ListItemContentVariant[ListItemVariant.Muted],
        className,
      )}
    >
      {children}
    </div>
  );
};

export const ListItemLabel = ({
  className,
  children,
  role = "label",
  ...props
}: ComponentProps<"label">) => {
  return (
    <label
      {...props}
      role={role}
      className={cn(
        "grow cursor-[inherit] group-[.disabled]/listitem:text-current/50",
        "whitespace-nowrap overflow-hidden",
        cn(className),
      )}
    >
      {children}
    </label>
  );
};

export const ListItemText = ({
  className,
  children,
  ...props
}: ComponentProps<"p">) => {
  return (
    <h6
      {...props}
      className={cn(
        "grow cursor-[inherit] group-[.disabled]/listitem:text-current/50",
        "whitespace-nowrap overflow-hidden",
        className,
      )}
    >
      {children}
    </h6>
  );
};

const ListItemExtraClassName =
  "text-muted-foreground cursor-[inherit] group-[.disabled]/listitem:text-current/50";
export const ListItemIcon = <Props extends { className?: string }>({
  as,
  className,
  ...props
}: AsProps<Props>) => {
  const Component = as;
  return (
    <Component
      {...props}
      className={cn("size-4", ListItemExtraClassName, className)}
    />
  );
};

export const ListItemSmall = ({
  className,
  children,
  ...props
}: ComponentProps<"small">) => {
  return (
    <small
      {...props}
      className={cn("truncate", ListItemExtraClassName, className)}
    >
      {children}
    </small>
  );
};

export const ListItemMore = ({
  className,
  ...props
}: ComponentProps<typeof ChevronRight>) => {
  return <ListItemIcon as={ChevronRight} {...props} />;
};

export const ListItemSwitch = ({
  className,
  ...props
}: ComponentProps<typeof Switch>) => {
  return <Switch {...props} className={cn("cursor-[inherit]", className)} />;
};
