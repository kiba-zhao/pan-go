import {
  AppBar as AdminAppBar,
  Layout,
  LayoutProps,
  LocalesMenuButton,
  TitlePortal,
  Toolbar,
} from "react-admin";

const AppBar = () => (
  <AdminAppBar>
    <Toolbar>
      <TitlePortal />
      <LocalesMenuButton />
    </Toolbar>
  </AdminAppBar>
);

export const AppLayout = (props: LayoutProps) => (
  <Layout {...props} appBar={AppBar} />
);
