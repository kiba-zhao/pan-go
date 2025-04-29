/**
 * App Component
 *
 * root component of the application
 */

import "./App.css";
import { QueryClient, QueryClientProvider } from "@tanstack/react-query";
import { I18nProvider } from "./components/i18n/Context";
import { Router, Route, Outlet } from "./components/Route/Router";

import { Master, MasterErrorBoundary } from "./components/Master/Layout";
// AppNode
import { AppNodeCreate, AppNodeEdit } from "./components/AppNode/Page";
import { AppNodeCreatePath, AppNodeEditPath } from "./components/AppNode/Route";
//
// ExtFSBrowseFile
import { default as ExtFSBrowseFile } from "./components/ExtFSBrowseFile/Page";
import { ExtFSBrowseFilePath } from "./components/ExtFSBrowseFile/Route";
//
// ExtFSNodeItem
import {
  ExtFSNodeItemCreate,
  ExtFSNodeItemEdit,
  ExtFSNodeItemView,
} from "./components/ExtFSNodeItem/Page";
import {
  ExtFSNodeItemCreatePath,
  ExtFSNodeItemEditPath,
  ExtFSNodeItemShowPath,
} from "./components/ExtFSNodeItem/Route";
//
// AppSettings
import { AppSettings } from "./components/AppSettings/Page";
import { AppSettingsPath } from "./components/AppSettings/Route";
//
// ExtFS
import { default as ExtFSHome } from "./components/ExtFS/Page";
import { ExtFSPath } from "./components/ExtFS/Route";
//

// Error Page
import { PageNotFound } from "./components/Master/Page";
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
const AppProvider = ({ children }: { children?: React.ReactNode }) => (
  <I18nProvider>
    <QueryClientProvider client={queryClient}>{children}</QueryClientProvider>
  </I18nProvider>
);

const AppLayout = () => (
  <Master>
    <Outlet />
  </Master>
);

export const App = () => (
  <AppProvider>
    <Router>
      <Route Component={AppLayout} ErrorBoundary={MasterErrorBoundary}>
        <Route index path={ExtFSPath} Component={ExtFSHome} />
        <Route path={AppSettingsPath} Component={AppSettings} />
        <Route path={ExtFSNodeItemCreatePath} Component={ExtFSNodeItemCreate} />
        <Route path={ExtFSNodeItemEditPath} Component={ExtFSNodeItemEdit} />
        <Route path={ExtFSNodeItemShowPath} Component={ExtFSNodeItemView} />
        <Route path={ExtFSBrowseFilePath} Component={ExtFSBrowseFile} />
        <Route path={AppNodeCreatePath} Component={AppNodeCreate} />
        <Route path={AppNodeEditPath} Component={AppNodeEdit} />
        <Route path="*" Component={PageNotFound} />
      </Route>
    </Router>
  </AppProvider>
);

export default App;
