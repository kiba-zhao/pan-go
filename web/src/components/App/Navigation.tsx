import { useTranslation } from "./I18Next";

import { cn } from "@/lib/utils";
import { type ComponentProps } from "react";
import { Link } from "./Route";
import { AppSubject, useAppDispatch, type AppContextState } from "./Context";

import { DashboardRoutePath, DashboardName } from "@/components/Dashboard/meta";
import DashboardIcon from "@/components/Dashboard/Icon";

import ChatOutlineIcon from "@/components/Chat/Icon";
import { ChatRoutePath, ChatName } from "@/components/Chat/meta";

import ClusterOutlinedIcon from "@/components/Clusters/Icon";
import { ClustersRoutePath, ClustersName } from "@/components/Clusters/meta";

import { AppsRoutePath, AppsName } from "@/components/Apps/meta";
import AppsIcon from "@/components/Apps/Icon";

import SettingsIcon from "@/components/Settings/Icon";
import { SettingsRoutePath, SettingsName } from "@/components/Settings/meta";

import NotificationsIcon from "@/components/Notifications/Icon";
import {
  NotificationsRoutePath,
  NotificationsName,
} from "@/components/Notifications/meta";

import {
  List,
  ListItem,
  ListItemText,
  ListItemSmall,
  ListItemVariant,
  ListItemMainClassName,
  ListItemAvatar as NavListItemAvatar,
} from "./List";

export type NavigationProps = Pick<ComponentProps<typeof NavList>, "className">;
export const AppNavigation = ({ className }: NavigationProps) => {
  const { t } = useTranslation();
  const dispatch = useAppDispatch();

  const handleMouseEnter = (state: AppContextState) => dispatch?.(state);

  return (
    <NavList className={className}>
      <NavListItem
        onMouseEnter={() => handleMouseEnter({ subjectVisible: false })}
      >
        <NavListItemLink to={DashboardRoutePath}>
          <NavListItemAvatar as={DashboardIcon} />
          <NavListItemText>{t(`navigations.${DashboardName}`)}</NavListItemText>
        </NavListItemLink>
      </NavListItem>
      <NavListItem>
        <NavListItemLink to={ChatRoutePath}>
          <NavListItemAvatar as={ChatOutlineIcon} />
          <NavListItemText>{t(`navigations.${ChatName}`)}</NavListItemText>
          <NavListItemSmall>Ctrl+I</NavListItemSmall>
        </NavListItemLink>
      </NavListItem>
      <NavListItem
        onMouseEnter={() =>
          handleMouseEnter({
            subject: AppSubject.Clusters,
            subjectVisible: true,
          })
        }
      >
        <NavListItemLink to={ClustersRoutePath}>
          <NavListItemAvatar as={ClusterOutlinedIcon} />
          <NavListItemText>{t(`navigations.${ClustersName}`)}</NavListItemText>
        </NavListItemLink>
      </NavListItem>
      <NavListItem>
        <NavListItemLink to={AppsRoutePath}>
          <NavListItemAvatar as={AppsIcon} />
          <NavListItemText>{t(`navigations.${AppsName}`)}</NavListItemText>
          <NavListItemSmall>Ctrl+A</NavListItemSmall>
        </NavListItemLink>
      </NavListItem>
    </NavList>
  );
};

export const AppSecondaryNavigation = ({ className }: NavigationProps) => {
  const { t } = useTranslation();
  return (
    <NavList className={className}>
      <NavListItem>
        <NavListItemLink to={NotificationsRoutePath}>
          <NavListItemAvatar as={NotificationsIcon} />
          <NavListItemText>
            {t(`navigations.${NotificationsName}`)}
          </NavListItemText>
        </NavListItemLink>
      </NavListItem>
      <NavListItem>
        <NavListItemLink to={SettingsRoutePath}>
          <NavListItemAvatar as={SettingsIcon} />
          <NavListItemText>{t(`navigations.${SettingsName}`)}</NavListItemText>
        </NavListItemLink>
      </NavListItem>
    </NavList>
  );
};

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
