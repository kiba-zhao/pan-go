import {
  NavList,
  NavListItem,
  NavListItemLink,
  NavListItemText,
  NavListItemSmall,
  NavListItemAvatar,
} from "./Navigation";

import { useTranslation, I18nVariant } from "./I18Next";

import { type ComponentProps } from "react";
import {
  AppSubject,
  useAppDispatch,
  withAppAction,
  type AppState,
} from "./Context";

import { DashboardRoutePath, DashboardName } from "@/components/Dashboard/meta";
import DashboardIcon from "@/components/Dashboard/Icon";

import ChatOutlineIcon from "@/components/Chat/Icon";
import { ChatRoutePath, ChatName } from "@/components/Chat/meta";

import { AppsRoutePath, AppsName } from "@/components/Apps/meta";
import AppsIcon from "@/components/Apps/Icon";

import SettingsIcon from "@/components/Settings/Icon";
import { SettingsRoutePath, SettingsName } from "@/components/Settings/meta";

import NotificationsIcon from "@/components/Notifications/Icon";
import {
  NotificationsRoutePath,
  NotificationsName,
} from "@/components/Notifications/meta";

export type NavigationProps = Pick<ComponentProps<typeof NavList>, "className">;
export const AppNavigation = ({ className }: NavigationProps) => {
  const dispatch = useAppDispatch();
  const handleMouseEnter = (state: AppState) =>
    dispatch?.(withAppAction(state));

  return (
    <NavList className={className}>
      <DashboardNavListItem
        onMouseEnter={() => handleMouseEnter({ subjectVisible: false })}
      />
      <ChatNavListItem
        onMouseEnter={() =>
          handleMouseEnter({
            subject: AppSubject.Chat,
            subjectVisible: true,
          })
        }
      />
      <AppsNavListItem
        onMouseEnter={() =>
          handleMouseEnter({
            subject: AppSubject.Apps,
            subjectVisible: true,
          })
        }
      />
    </NavList>
  );
};

export const AppSecondaryNavigation = ({ className }: NavigationProps) => (
  <NavList className={className}>
    <NotificationsNavListItem />
    <SettingsNavListItem />
  </NavList>
);

type CustomNavListItemProps = Omit<
  ComponentProps<typeof NavListItem>,
  "children"
>;
const DashboardNavListItem = (props: CustomNavListItemProps) => {
  const { t } = useTranslation();
  return (
    <NavListItem {...props}>
      <NavListItemLink to={DashboardRoutePath}>
        <NavListItemAvatar as={DashboardIcon} />
        <NavListItemText>
          {t(`${I18nVariant.Navigation}.${DashboardName}`)}
        </NavListItemText>
      </NavListItemLink>
    </NavListItem>
  );
};

const ChatNavListItem = (props: CustomNavListItemProps) => {
  const { t } = useTranslation();
  return (
    <NavListItem {...props}>
      <NavListItemLink to={ChatRoutePath}>
        <NavListItemAvatar as={ChatOutlineIcon} />
        <NavListItemText>
          {t(`${I18nVariant.Navigation}.${ChatName}`)}
        </NavListItemText>
        <NavListItemSmall>Ctrl+I</NavListItemSmall>
      </NavListItemLink>
    </NavListItem>
  );
};

const AppsNavListItem = (props: CustomNavListItemProps) => {
  const { t } = useTranslation();
  return (
    <NavListItem {...props}>
      <NavListItemLink to={AppsRoutePath}>
        <NavListItemAvatar as={AppsIcon} />
        <NavListItemText>
          {t(`${I18nVariant.Navigation}.${AppsName}`)}
        </NavListItemText>
        <NavListItemSmall>Ctrl+A</NavListItemSmall>
      </NavListItemLink>
    </NavListItem>
  );
};

const NotificationsNavListItem = (props: CustomNavListItemProps) => {
  const { t } = useTranslation();
  return (
    <NavListItem {...props}>
      <NavListItemLink to={NotificationsRoutePath}>
        <NavListItemAvatar as={NotificationsIcon} />
        <NavListItemText>
          {t(`${I18nVariant.Navigation}.${NotificationsName}`)}
        </NavListItemText>
      </NavListItemLink>
    </NavListItem>
  );
};

const SettingsNavListItem = (props: CustomNavListItemProps) => {
  const { t } = useTranslation();
  return (
    <NavListItem {...props}>
      <NavListItemLink to={SettingsRoutePath}>
        <NavListItemAvatar as={SettingsIcon} />
        <NavListItemText>
          {t(`${I18nVariant.Navigation}.${SettingsName}`)}
        </NavListItemText>
      </NavListItemLink>
    </NavListItem>
  );
};
