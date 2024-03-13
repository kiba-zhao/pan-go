import "@fontsource/roboto/300.css";
import "@fontsource/roboto/400.css";
import "@fontsource/roboto/500.css";
import "@fontsource/roboto/700.css";
import { useEffect } from "react";

import { I18nextProvider } from "react-i18next";
import i18n from "./i18n";

import { QueryClient, QueryClientProvider } from "@tanstack/react-query";
import { ReactQueryDevtools } from "@tanstack/react-query-devtools";
import { APIProvider } from "./api.tsx";
import Header from "./components/Header";
import RouteView from "./routes";

const queryClient = new QueryClient();
function App() {
  // const prefersDarkMode = useMediaQuery('(prefers-color-scheme: dark)');
  useEffect(() => {
    document.title = "pan-go";
    console.log(1111, import.meta.env);
  }, []);

  return (
    <I18nextProvider i18n={i18n}>
      <APIProvider>
        <QueryClientProvider client={queryClient}>
          <ReactQueryDevtools initialIsOpen={false} />
          <Header />
          <RouteView />
        </QueryClientProvider>
      </APIProvider>
    </I18nextProvider>
  );
}

export default App;
