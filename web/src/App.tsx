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

import { APPNodeEdit, AppNodeCreate, AppNodes } from "./components/AppNodes";
import { AppSettings } from "./components/AppSettings";
import Dashboard from "./components/Dashboard";
import ExtFSHome from "./components/ExtFS";
import {
  ExtFSLocalFileCreate,
  ExtFSLocalFileEdit,
} from "./components/ExtFSLocalFile";
import {
  ExtFSLocalNodeSettings,
  ExtFSRemoteNodeSettings,
} from "./components/ExtFSNode";
import { ExtFSRemoteFileEdit } from "./components/ExtFSRemoteFile";
import {
  ExtFSTagCreate,
  ExtFSTagEdit,
  ExtFSTagList,
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
            <Route path="/app/settings/*" element={<AppSettings />} />
            <Route path="/extfs" element={<ExtFSHome />} />
            <Route
              path="/extfs/local-node"
              element={<ExtFSLocalNodeSettings />}
            />
            <Route
              path="/extfs/remote-node/:id"
              element={<ExtFSRemoteNodeSettings />}
            />
            <Route
              path="/extfs/local-file/create"
              element={<ExtFSLocalFileCreate />}
            />
            <Route
              path="/extfs/local-file/:id"
              element={<ExtFSLocalFileEdit />}
            />
            <Route
              path="/extfs/remote-file/:id"
              element={<ExtFSRemoteFileEdit />}
            />
          </CustomRoutes>
          <Resource
            name="app/nodes"
            list={AppNodes}
            create={AppNodeCreate}
            edit={APPNodeEdit}
          />
          <Resource
            name="extfs/tags"
            list={ExtFSTagList}
            create={ExtFSTagCreate}
            edit={ExtFSTagEdit}
            show={ExtFSTagShow}
          />
        </Admin>
      </APIProvider>
    </BrowserRouter>
  );
};

export default App;
