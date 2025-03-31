/**
 * Layout Component Definition File
 */
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

/**
 * AppMenu
 *
 * Menu component for AppLayout
 *
 * This menu is composed of all the main sections of the application.
 *
 *
 * @returns {ReactElement} The menu element
 */
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

/**
 * AppBar
 *
 * Top Bar for AppLayout
 *
 *
 * @param {AppBarProps} props The AppBar props
 * @returns {ReactElement} The AppBar element
 */
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

/**
 * AppLayout
 *
 * Layout for App
 *
 * @example
 * import { AppLayout } from "./Layout";
 *
 * const App = () => (
 *   <Admin dataProvider={dataProvider} layout={AppLayout}>
 *     ...
 *   </Admin>
 * );
 *
 * @param {LayoutProps} props The layout props
 * @returns {ReactElement} The layout element
 */
export const AppLayout = (props: LayoutProps) => (
  <Layout {...props} appBar={AppBar} menu={AppMenu} />
);
