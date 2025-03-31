import React from "react";
import ReactDOM from "react-dom/client";
import App from "./App.tsx";
import { BrowserProvider } from "./components/Global/Browser.tsx";
import "./index.css";

ReactDOM.createRoot(document.getElementById("root")!).render(
  <React.StrictMode>
    <BrowserProvider window={window}>
      <App />
    </BrowserProvider>
  </React.StrictMode>
);
