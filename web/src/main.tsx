import "@fontsource/roboto/300.css";
import "@fontsource/roboto/400.css";
import "@fontsource/roboto/500.css";
import "@fontsource/roboto/700.css";
import "./scripts/colorScheme.js";

import React from "react";
import ReactDOM from "react-dom/client";
import App from "./App.tsx";
import { BrowserProvider } from "./components/Browser.tsx";
// import { ColorSchemeScript } from "./components/Master/Theme";

ReactDOM.createRoot(document.getElementById("root")!).render(
  <React.StrictMode>
    <BrowserProvider window={window}>
      {/* <ColorSchemeScript /> */}
      <App />
    </BrowserProvider>
  </React.StrictMode>
);
