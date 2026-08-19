import { lazy } from "react";
import {
  Outlets as AppOutlets,
  MainSuspense,
  ExtraSuspense,
} from "@/components/App/Outlets";

const ClusterMain = lazy(() =>
  import("./Main").then((m) => ({ default: m.ClusterMain })),
);
const ClusterExtra = lazy(() =>
  import("./Extra").then((m) => ({ default: m.ClusterExtra })),
);
const ClusterHeaderExtra = lazy(() =>
  import("./HeaderExtra").then((m) => ({ default: m.ClusterHeaderExtra })),
);

export const ClusterOutlets = () => (
  <AppOutlets
    extra={
      <ExtraSuspense>
        <ClusterExtra />
      </ExtraSuspense>
    }
    headerExtra={
      <ExtraSuspense>
        <ClusterHeaderExtra />
      </ExtraSuspense>
    }
  >
    <MainSuspense>
      <ClusterMain />
    </MainSuspense>
  </AppOutlets>
);
