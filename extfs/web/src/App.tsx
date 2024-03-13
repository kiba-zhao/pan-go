import { Admin } from "react-admin";

import { dataProvider } from "./api";

export const App = () => <Admin dataProvider={dataProvider}></Admin>;

export default App;
