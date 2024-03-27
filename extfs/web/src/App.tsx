import { Admin, RaThemeOptions, Resource, defaultTheme } from "react-admin";
import { QueryClient } from "react-query";
import { BrowserRouter } from "react-router-dom";
import { dataProvider } from "./api";
import { useI18nProvider } from "./i18n";

import Dashboard from "./components/Dashboard";
import { AppLayout } from "./components/Layout";
import NotFound from "./components/NotFound";
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
  console.log(1111, import.meta.env.BASE_URL);
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
          name="targets"
          list={Targets}
          edit={TargetEdit}
          show={TargetShow}
          create={TargetCreate}
        />
      </Admin>
    </BrowserRouter>
  );
};

export default App;
