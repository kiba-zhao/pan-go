import { QueryClient } from "@tanstack/react-query";
import {
  Admin,
  CustomRoutes,
  RaThemeOptions,
  Resource,
  defaultTheme,
} from "react-admin";
import { BrowserRouter, Route } from "react-router-dom";
import { dataProvider } from "./api";
import { APIProvider } from "./components/API";
import { useI18nProvider } from "./i18n";

import {
  APPNodeEdit,
  AppNodeCreate,
  AppNodeRoutePath,
  AppNodes,
} from "./components/AppNodes";
import { AppSettings, AppSettingsRoutePath } from "./components/AppSettings";
import Dashboard from "./components/Dashboard";
import ExtFSHome, { ExtFSRoutePath } from "./components/ExtFS";
import {
  ExtFSNodeItemCreate,
  ExtFSNodeItemEdit,
  ExtFSNodeItemRoutePath,
  ExtFSNodeItemView,
} from "./components/ExtFSNodeItem";

import ExtFSBrowseFile, {
  RoutePath as ExtFSBrowseFileRoutePath,
} from "./components/ExtFSBrowseFile";
import { AppLayout } from "./components/Layout";
import NotFound from "./components/NotFound";

const queryClient = new QueryClient();
const darkTheme: RaThemeOptions = {
  ...defaultTheme,
  palette: { mode: "dark" },
};
export const App = () => {
  const i18nProvider = useI18nProvider();
  if (!i18nProvider) return null;
  return (
    <BrowserRouter basename={import.meta.env.BASE_URL}>
      <APIProvider>
        <Admin
          disableTelemetry
          theme={defaultTheme}
          darkTheme={darkTheme}
          dataProvider={dataProvider}
          i18nProvider={i18nProvider}
          catchAll={NotFound}
          dashboard={Dashboard}
          queryClient={queryClient}
          layout={AppLayout}
        >
          <CustomRoutes>
            <Route path={AppSettingsRoutePath} element={<AppSettings />} />
            <Route path={ExtFSRoutePath} element={<ExtFSHome />} />

            <Route
              path={ExtFSBrowseFileRoutePath}
              element={<ExtFSBrowseFile />}
            />
            <Route
              path={`${ExtFSNodeItemRoutePath}/create`}
              element={<ExtFSNodeItemCreate />}
            />
            <Route
              path={`${ExtFSNodeItemRoutePath}/:id`}
              element={<ExtFSNodeItemEdit />}
            />
            <Route
              path={`${ExtFSNodeItemRoutePath}/:id/show`}
              element={<ExtFSNodeItemView />}
            />
          </CustomRoutes>
          <Resource
            name={AppNodeRoutePath.substring(1)}
            list={AppNodes}
            create={AppNodeCreate}
            edit={APPNodeEdit}
          />
        </Admin>
      </APIProvider>
    </BrowserRouter>
  );
};

export default App;
