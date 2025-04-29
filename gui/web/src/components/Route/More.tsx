import { ReactNode, useMemo } from "react";
import { More } from "../Common/MoreButton";
import { AppSettingsMore } from "../AppSettings/More";
import { AppSettingsPath } from "../AppSettings/Route";
import { AppNodeCreatePath } from "../AppNode/Route";
import { NewAppNodeMore } from "../AppNode/More";
import { ExtFSMore } from "../ExtFS/More";
import { ExtFSPath } from "../ExtFS/Route";
import { NewExtFSNodeItemMore } from "../ExtFSNodeItem/More";
import { ExtFSNodeItemCreatePath } from "../ExtFSNodeItem/Route";

export const RouteMore = ({
  path,
  children,
}: {
  path?: string;
  children?: ReactNode;
}) => {
  const routes = useMemo(() => {
    const routes_: ReactNode[] = [];
    if (path !== AppSettingsPath)
      routes_.push(<AppSettingsMore key="app-settings" />);
    if (path !== ExtFSPath) routes_.push(<ExtFSMore key="extfs" />);
    if (path !== AppNodeCreatePath)
      routes_.push(<NewAppNodeMore key="app-node-create" />);
    if (path !== ExtFSNodeItemCreatePath)
      routes_.push(<NewExtFSNodeItemMore key="extfs-node-item-create" />);
    return routes_;
  }, [path]);
  return (
    <More>
      {routes}
      {children}
    </More>
  );
};
