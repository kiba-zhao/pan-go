import React from "react";
import ReactDOM from "react-dom/client";
import App from "@/App.tsx";
import { BrowserProvider } from "@/components/App/Browser";

ReactDOM.createRoot(document.getElementById("root")!).render(
  <React.StrictMode>
    <BrowserProvider window={window}>
      <App />
    </BrowserProvider>
  </React.StrictMode>,
);
