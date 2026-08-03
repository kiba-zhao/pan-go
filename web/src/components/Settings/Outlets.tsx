import { lazy } from "react";
import { Outlets, MainSuspense, ExtraSuspense } from "@/components/App/Outlets";

const SettingsMain = lazy(() => import("./Main"));
const SettingsExtra = lazy(() => import("./Extra"));
const SettingsOutlets = () => (
  <Outlets
    extra={
      <ExtraSuspense>
        <SettingsExtra />
      </ExtraSuspense>
    }
  >
    <MainSuspense>
      <SettingsMain />
    </MainSuspense>
  </Outlets>
);

export default SettingsOutlets;
