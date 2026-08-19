import { cn } from "@/lib/utils";
import { type ComponentProps } from "react";
import { Link } from "./Router";

import {
  List,
  ListItem,
  ListItemText,
  ListItemSmall,
  ListItemVariant,
  ListItemMainClassName,
  ListItemAvatar as NavListItemAvatar,
} from "./List";

export const NavList = ({
  children,
  className,
  ...props
}: ComponentProps<typeof List>) => {
  return (
    <nav>
      <List {...props} className={cn("group/nav-list", className)}>
        {children}
      </List>
    </nav>
  );
};

export const NavListItem = ({
  children,
  className,
  hover = ListItemVariant.Accent,
  ...props
}: ComponentProps<typeof ListItem>) => {
  return (
    <ListItem
      {...props}
      hover={hover}
      className={cn(
        // "not-has-[.listitem-link]:px-1 not-has-[.listitem-content]:py-1.5",
        className,
      )}
    >
      {children}
    </ListItem>
  );
};

export const NavListItemLink = ({
  children,
  className,
  ...props
}: ComponentProps<typeof Link>) => {
  return (
    <Link
      {...props}
      className={cn(ListItemMainClassName, "py-1.5 px-1.5", className)}
    >
      {children}
    </Link>
  );
};

export const NavListItemText = ({
  children,
  className,
  ...props
}: ComponentProps<typeof ListItemText>) => {
  return (
    <ListItemText
      {...props}
      className={cn("group-[.collapsed]/nav-list:hidden", className)}
    >
      {children}
    </ListItemText>
  );
};

export const NavListItemSmall = ({
  children,
  className,
  ...props
}: ComponentProps<typeof ListItemSmall>) => {
  return (
    <ListItemSmall
      {...props}
      className={cn(
        "font-semibold group-[.collapsed]/nav-list:hidden",
        className,
      )}
    >
      {children}
    </ListItemSmall>
  );
};

export { NavListItemAvatar };
