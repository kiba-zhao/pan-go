import { Routes, Route } from "react-router-dom";
import Home from "./components/Home";
import Modules from "./components/Modules";
import NotFound from "./components/NotFound";

function RouteView() {
  return (
    <Routes>
      <Route path="/" element={<Home />}></Route>
      <Route path="/modules" element={<Modules />}></Route>
      <Route path="*" element={<NotFound />}></Route>
    </Routes>
  );
}

export default RouteView;
