import { Admin, RaThemeOptions, Resource, defaultTheme } from "react-admin";
import { QueryClient } from "react-query";
import { BrowserRouter } from "react-router-dom";
import { dataProvider } from "./api";
import { useI18nProvider } from "./i18n";

import Dashboard from "./components/Dashboard";
import { AppLayout } from "./components/Layout";
import NotFound from "./components/NotFound";
import {
  TargetFileIcon,
  TargetFileShow,
  TargetFiles,
} from "./components/TargetFiles";
import {
  TargetCreate,
  TargetEdit,
  TargetShow,
  Targets,
} from "./components/Targets";

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
        <Resource
          name="extfs/targets"
          list={Targets}
          edit={TargetEdit}
          show={TargetShow}
          create={TargetCreate}
        />
        <Resource
          name="extfs/target-files"
          icon={TargetFileIcon}
          list={TargetFiles}
          show={TargetFileShow}
          hasEdit={false}
          hasCreate={false}
        />
      </Admin>
    </BrowserRouter>
  );
};

export default App;
