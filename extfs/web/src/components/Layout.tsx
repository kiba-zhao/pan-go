import {
  AppBar as AdminAppBar,
  Layout,
  LayoutProps,
  LocalesMenuButton,
  RefreshIconButton,
  TitlePortal,
  ToggleThemeButton,
} from "react-admin";

const AppBar = () => (
  <AdminAppBar
    toolbar={
      <>
        <TitlePortal />
        <ToggleThemeButton />
        <LocalesMenuButton />
        <RefreshIconButton />
      </>
    }
  ></AdminAppBar>
);

export const AppLayout = (props: LayoutProps) => (
  <Layout {...props} appBar={AppBar} />
);
