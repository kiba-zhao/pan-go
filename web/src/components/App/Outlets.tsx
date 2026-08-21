import {
  ComponentProps,
  createContext,
  useContext,
  Suspense,
  type PropsWithChildren,
  type ReactNode,
} from "react";
import { Outlet as RouteOutlet } from "./Router";

type OutletsProps = PropsWithChildren<{
  extra?: ReactNode;
  headerExtra?: ReactNode;
  footer?: ReactNode;
}>;

type OutletName = Exclude<keyof OutletsProps, "children">;
const OutletsContext = createContext<{
  name?: OutletName;
  children?: ReactNode;
}>({});

export const Outlets = (props: OutletsProps) => {
  const { name, children } = useContext(OutletsContext);
  return props[name || "children"] || children;
};

type OutletProps = PropsWithChildren<
  Pick<ComponentProps<typeof RouteOutlet>, "context"> & {
    name?: OutletName;
  } & Pick<ComponentProps<typeof Suspense>, "fallback">
>;
export const Outlet = ({ name, context, children, ...props }: OutletProps) => (
  <OutletsContext.Provider value={{ name, children }}>
    <Suspense {...props}>
      <RouteOutlet context={context} />
    </Suspense>
  </OutletsContext.Provider>
);
