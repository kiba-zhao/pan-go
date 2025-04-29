import {
  BrowserRouter,
  Routes,
  Route,
  Outlet,
  useNavigate,
  NavLink,
  useParams,
  useSearchParams,
} from "react-router";
import type { BrowserRouterProps, To } from "react-router";

export type RouterProps = Pick<BrowserRouterProps, "basename" | "children">;
export const Router = ({ basename, children }: RouterProps) => {
  return (
    <BrowserRouter basename={basename || import.meta.env.BASE_URL}>
      <Routes>{children}</Routes>
    </BrowserRouter>
  );
};

export {
  Route,
  Outlet,
  useNavigate,
  NavLink as Link,
  useParams,
  useSearchParams,
};
export type { To };
