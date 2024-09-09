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

import { AppNodeIcon } from "./AppNodes";
import { AppSettingsIcon } from "./AppSettings";
import { ExfFSIcon } from "./ExtFS";

const AppMenu = () => {
  const t = useTranslate();
  return (
    <Menu>
      <Menu.DashboardItem />
      <Menu.Item
        to="/extfs"
        primaryText={t("custom.extfs.name")}
        leftIcon={<ExfFSIcon />}
      />
      <Menu.Item
        to="/app/nodes"
        primaryText={t("resources.app/nodes.name")}
        leftIcon={<AppNodeIcon />}
      />
      <Menu.Item
        to="/app/settings"
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
