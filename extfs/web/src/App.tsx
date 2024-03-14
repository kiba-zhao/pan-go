import { Admin, Resource } from "react-admin";
import { dataProvider } from "./api";

import polyglotI18nProvider from "ra-i18n-polyglot";
import englishMessages from "ra-language-english";

import { TargetEdit, Targets } from "./components/Targets";

const i18nProvider = polyglotI18nProvider(
  (locale) =>
    !locale || locale == "en"
      ? englishMessages
      : import(`./locales/${locale}/translation.json`),
  "en", // Default locale
  [
    { locale: "en", name: "English" },
    // { locale: "fr", name: "Français" },
  ]
);

export const App = () => (
  <Admin
    dataProvider={dataProvider}
    i18nProvider={i18nProvider}
    // layout={AppLayout}
  >
    <Resource name="targets" list={Targets} edit={TargetEdit} />
  </Admin>
);

export default App;
