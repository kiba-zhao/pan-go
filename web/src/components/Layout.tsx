import {
  AppBar as AdminAppBar,
  Layout,
  LayoutProps,
  LocalesMenuButton,
  Menu,
  TitlePortal,
  ToggleThemeButton,
  useTranslate,
} from "react-admin";

import { AppNodeIcon, AppNodeRoutePath } from "./AppNodes";
import { AppSettingsIcon, AppSettingsRoutePath } from "./AppSettings";
import { ExtFSIcon, ExtFSRoutePath } from "./ExtFS";

const AppMenu = () => {
  const t = useTranslate();
  return (
    <Menu>
      <Menu.DashboardItem />
      <Menu.Item
        to={ExtFSRoutePath}
        primaryText={t("custom.extfs.name")}
        leftIcon={<ExtFSIcon />}
      />
      <Menu.Item
        to={AppNodeRoutePath}
        primaryText={t("resources.app/nodes.name")}
        leftIcon={<AppNodeIcon />}
      />
      <Menu.Item
        to={AppSettingsRoutePath}
        primaryText={t("custom.app/settings.name")}
        leftIcon={<AppSettingsIcon />}
      />
    </Menu>
  );
};

const AppBar = () => (
  <AdminAppBar
    toolbar={
      <>
        <TitlePortal />
        <LocalesMenuButton />
        <ToggleThemeButton />
      </>
    }
  ></AdminAppBar>
);

export const AppLayout = (props: LayoutProps) => (
  <Layout {...props} appBar={AppBar} menu={AppMenu} />
);
