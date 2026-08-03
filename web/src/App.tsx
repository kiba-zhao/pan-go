import "@/App.css";
import { Provider as I18NextProvider } from "@/components/App/I18Next";
import { QueryClient, QueryClientProvider } from "@/components/App/ReactQuery";

import AppPage from "@/components/App/Page";

const queryClient = new QueryClient();
const App = () => (
  <QueryClientProvider client={queryClient}>
    <I18NextProvider>
      <AppPage />
    </I18NextProvider>
  </QueryClientProvider>
);

export default App;
