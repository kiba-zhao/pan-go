import { ClusterMain } from "./Main";
import { ClusterExtra } from "./Extra";
import { ClusterHeaderExtra } from "./HeaderExtra";

import { Outlets as AppOutlets } from "@/components/App/Outlets";

const ClusterOutlets = () => (
  <AppOutlets extra={<ClusterExtra />} headerExtra={<ClusterHeaderExtra />}>
    <ClusterMain />
  </AppOutlets>
);
export default ClusterOutlets;
