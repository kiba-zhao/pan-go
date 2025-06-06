/**
 * App Component
 *
 * root component of the application
 */

import "./App.css";

import { QueryClient, QueryClientProvider } from "@tanstack/react-query";
import { Provider as I18NextProvider } from "./components/I18Next/Context";
import { ThemeProvider } from "./components/Master/Theme";
import { Router, Route, Outlet } from "./components/Route/Router";

import { Master } from "./components/Master/Layout";
// AppNode
import { AppNodeCreatePage, AppNodeEditPage } from "./components/AppNode/Page";
import { AppNodeCreatePath, AppNodeEditPath } from "./components/AppNode/Route";
//
// ExtFSBrowseFile
import { default as ExtFSBrowseFilePage } from "./components/ExtFSBrowseFile/Page";
import { ExtFSBrowseFilePath } from "./components/ExtFSBrowseFile/Route";
//
// ExtFSNodeItem
import {
  ExtFSNodeItemCreatePage,
  ExtFSNodeItemEditPage,
  ExtFSNodeItemViewPage,
} from "./components/ExtFSNodeItem/Page";
import {
  ExtFSNodeItemCreatePath,
  ExtFSNodeItemEditPath,
  ExtFSNodeItemShowPath,
} from "./components/ExtFSNodeItem/Route";
//
// AppSettings
import { AppSettingsPage } from "./components/AppSettings/Page";
import { AppSettingsPath } from "./components/AppSettings/Route";
//
// ExtFS
import { default as ExtFSPage } from "./components/ExtFS/Page";
import { ExtFSPath } from "./components/ExtFS/Route";
//

// Error Page
import { NotFoundPage } from "./components/Master/Page";
//

// import { Component, lazy } from "react";

// async function lazy(path: string) {
//   const { default: Component_, Component } = await import(path);
//   return { Component: Component || Component_ };
// }
// type LazyLoad = () => Promise<Component>;
// const Lazy = ({ component }: { component: Component }) => {
//   const Component = lazy(component);

//   return (
//     <Suspense fallback={<LoadingSpinner />}>
//       <Component />
//     </Suspense>
//   );
// };

const queryClient = new QueryClient();
const AppLayout = () => (
  <QueryClientProvider client={queryClient}>
    <I18NextProvider>
      <ThemeProvider>
        <Master>
          <Outlet />
        </Master>
      </ThemeProvider>
    </I18NextProvider>
  </QueryClientProvider>
);

export const App = () => (
  <Router>
    <Route Component={AppLayout}>
      <Route index path={ExtFSPath} Component={ExtFSPage} />
      <Route path={AppSettingsPath} Component={AppSettingsPage} />
      <Route
        path={ExtFSNodeItemCreatePath}
        Component={ExtFSNodeItemCreatePage}
      />
      <Route path={ExtFSNodeItemEditPath} Component={ExtFSNodeItemEditPage} />
      <Route path={ExtFSNodeItemShowPath} Component={ExtFSNodeItemViewPage} />
      <Route path={ExtFSBrowseFilePath} Component={ExtFSBrowseFilePage} />
      <Route path={AppNodeCreatePath} Component={AppNodeCreatePage} />
      <Route path={AppNodeEditPath} Component={AppNodeEditPage} />
      <Route path="*" Component={NotFoundPage} />
    </Route>
  </Router>
);

export default App;
