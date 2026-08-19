import {
  BrowserRouter,
  Routes,
  Route,
  Outlet,
  useNavigate,
  NavLink,
  useParams,
  useSearchParams,
  useOutlet,
  generatePath,
  useLocation,
} from "react-router";
import type { BrowserRouterProps } from "react-router";

export type RouterProps = Pick<BrowserRouterProps, "basename" | "children">;
export const Router = ({ basename, children }: RouterProps) => {
  return <BrowserRouter basename={basename}>{children}</BrowserRouter>;
};

export {
  Route,
  Routes,
  Outlet,
  useNavigate,
  NavLink as Link,
  useParams,
  useSearchParams,
  useOutlet,
  generatePath,
  useLocation,
};
