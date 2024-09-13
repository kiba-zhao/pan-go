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
import { APIProvider } from "./API";
import { useI18nProvider } from "./i18n";

import {
  APPNodeEdit,
  AppNodeCreate,
  AppNodeRoutePath,
  AppNodes,
} from "./components/AppNodes";
import { AppSettings, AppSettingsRoutePath } from "./components/AppSettings";
import Dashboard from "./components/Dashboard";
import ExtFSHome, { ExtFSProvider, ExtFSRoutePath } from "./components/ExtFS";
import {
  ExtFSNodeItemCreate,
  ExtFSNodeItemEdit,
  ExtFSNodeItemRoutePath,
} from "./components/ExtFSNodeItem";
import {
  ExtFSTagCreate,
  ExtFSTagEdit,
  ExtFSTagList,
  ExtFSTagRoutePath,
  ExtFSTagShow,
} from "./components/ExtFSTag";
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
        <ExtFSProvider>
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
                path={`${ExtFSNodeItemRoutePath}/create`}
                element={<ExtFSNodeItemCreate />}
              />
              <Route
                path={`${ExtFSNodeItemRoutePath}/:id`}
                element={<ExtFSNodeItemEdit />}
              />
            </CustomRoutes>
            <Resource
              name={AppNodeRoutePath.substring(1)}
              list={AppNodes}
              create={AppNodeCreate}
              edit={APPNodeEdit}
            />
            <Resource
              name={ExtFSTagRoutePath.substring(1)}
              list={ExtFSTagList}
              create={ExtFSTagCreate}
              edit={ExtFSTagEdit}
              show={ExtFSTagShow}
            />
          </Admin>
        </ExtFSProvider>
      </APIProvider>
    </BrowserRouter>
  );
};

export default App;
