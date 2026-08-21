import type { ReactNode, ComponentProps } from "react";
import { lazy } from "react";
import { Route, Routes } from "./Router";

import { NotFoundOutlets } from "./PageOutlets";
import { DashboardRoutePath } from "@/components/Dashboard/meta";
const DashboardOutlets = lazy(() => import("@/components/Dashboard/Outlets"));

import { SettingsRoutePath } from "@/components/Settings/meta";
const SettingsOutlets = lazy(() => import("@/components/Settings/Outlets"));

import { ClusterRoutePath } from "@/components/DeviceCluster/meta";
const ClusterOutlets = lazy(
  () => import("@/components/DeviceCluster/ClusterOutlets"),
);

type PageRoutesProps = {
  layout: ReactNode;
} & Omit<ComponentProps<typeof Routes>, "children">;
const PageRoutes = ({ layout, ...props }: PageRoutesProps) => (
  <Routes {...props}>
    <Route element={layout}>
      <Route path={DashboardRoutePath} Component={DashboardOutlets} />
      <Route path={SettingsRoutePath} Component={SettingsOutlets} />
      <Route path={ClusterRoutePath} Component={ClusterOutlets} />
      <Route path="*" Component={NotFoundOutlets} />
    </Route>
  </Routes>
);

export default PageRoutes;
