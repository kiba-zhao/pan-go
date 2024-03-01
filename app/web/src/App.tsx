import { useEffect } from "react";
import "@fontsource/roboto/300.css";
import "@fontsource/roboto/400.css";
import "@fontsource/roboto/500.css";
import "@fontsource/roboto/700.css";

import { I18nextProvider } from "react-i18next";
import i18n from "./i18n";

import * as api from "./api";
import { Context } from "./api.tsx";
import RouteView from "./routes";
import Header from "./components/Header";

export function APIProvider({ children }: { children: React.ReactNode }) {
  return <Context.Provider value={api}>{children}</Context.Provider>;
}
function App() {
  useEffect(() => {
    document.title = "pan-go";
    console.log(1111, import.meta.env);
  }, []);

  return (
    <I18nextProvider i18n={i18n}>
      <APIProvider>
        <Header />
        <RouteView />
      </APIProvider>
    </I18nextProvider>
  );
}

export default App;
