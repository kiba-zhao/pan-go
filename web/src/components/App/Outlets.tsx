import {
  ComponentProps,
  createContext,
  useContext,
  Suspense,
  type PropsWithChildren,
  type ReactNode,
} from "react";
import { useOutlet, Outlet as RouteOutlet } from "./Route";
import { AppLoading } from "./Loading";

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
  }
>;
export const Outlet = ({ name, context, children }: OutletProps) => {
  const outlet = useOutlet(context);
  return (
    <OutletsContext.Provider value={{ name, children }}>
      {outlet}
    </OutletsContext.Provider>
  );
};

export const MainSuspense = ({
  children,
  fallback = <AppLoading />,
  ...props
}: ComponentProps<typeof Suspense>) => (
  <Suspense {...props} fallback={fallback}>
    {children}
  </Suspense>
);

export { Suspense as ExtraSuspense };
