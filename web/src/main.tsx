import React from "react";
import ReactDOM from "react-dom/client";
import App from "./App.tsx";
import { BrowerProvider } from "./components/Global/Brower.tsx";
// import "./index.css";

ReactDOM.createRoot(document.getElementById("root")!).render(
  <React.StrictMode>
    <BrowerProvider window={window}>
      <App />
    </BrowerProvider>
  </React.StrictMode>
);
