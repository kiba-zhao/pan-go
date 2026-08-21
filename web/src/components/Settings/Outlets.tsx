import Main from "./Main";
import Extra from "./Extra";

import { Outlets } from "@/components/App/Outlets";

const SettingsOutlets = () => (
  <Outlets extra={<Extra />}>
    <Main />
  </Outlets>
);

export default SettingsOutlets;
