import { lazy } from "react";
import { Outlets, MainSuspense } from "@/components/App/Outlets";

const DashboardMain = lazy(() => import("./Main"));
const DashboardOutlets = () => (
  <Outlets>
    <MainSuspense>
      <DashboardMain />
    </MainSuspense>
  </Outlets>
);

export default DashboardOutlets;
