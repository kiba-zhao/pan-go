import "./Theme.css";
import "./Page.css";

import { useMedia } from "@/lib/hooks";
import { useTranslation } from "./I18Next";
import { Router, Route } from "./Route";
import { Outlet, Outlets } from "./Outlets";
import { useMemo } from "react";
import { cn } from "@/lib/utils";
import {
  AppContextProvider,
  useAppContext,
  useAppDispatch,
  AppSubject,
  AppAsideMode,
} from "./Context";
import { AsideModeControl, Separator } from "./Tools";
import AppBreadcrumb from "./Breadcrumb";

import { AppNavigation, AppSecondaryNavigation } from "./Navigation";
import ChatNavigation from "@/components/Chat/Navigation";
import ClustersNavigation from "@/components/Clusters/Navigation";
import AppsNavigation from "@/components/Apps/Navigation";

import { BlockOutlineIcon } from "./Icon";
import { PageInfoSection } from "./Layout";

import { DashboardRoutePath } from "@/components/Dashboard/meta";
import DashboardOutlets from "@/components/Dashboard/Outlets";
import { SettingsRoutePath } from "@/components/Settings/meta";
import { default as SettingsOutlets } from "@/components/Settings/Outlets";
import { Toaster } from "./Toast";

const AppPage = () => {
  return (
    <AppContextProvider>
      <Router>
        <Route element={<AppLayout />}>
          <Route path={DashboardRoutePath} Component={DashboardOutlets} />
          <Route path={SettingsRoutePath} Component={SettingsOutlets} />
          <Route path="*" Component={NotFoundOutlets} />
        </Route>
      </Router>
    </AppContextProvider>
  );
};

export default AppPage;

const AppLayout = () => {
  return (
    <>
      <AppHeader />
      <AppMain />
      <AppFooter />
      <AppAside />
      <Outlet name="extra" />
      <Toaster />
    </>
  );
};

const useLayoutElementClassName = () => {
  const { asideMode } = useAppContext();
  if (asideMode === AppAsideMode.Collapsed) {
    return "md:pl-14";
  }

  if (asideMode === AppAsideMode.Hidden) {
    return "";
  }

  return "md:pl-(--container-3xs)";
};

const AppContainerClassName = "max-w-7xl w-full";
const AppHeader = () => {
  const layoutClassName = useLayoutElementClassName();

  return (
    <header
      className={cn(
        "bg-background h-16 border-b z-10 w-full fixed flex justify-center",
        layoutClassName,
      )}
    >
      <section
        className={cn(
          "px-4 h-full flex flex-row items-center gap-2",
          AppContainerClassName,
        )}
      >
        <AsideModeControl />
        <Separator />
        <AppBreadcrumb className="grow" />
        <Outlet name="headerExtra" />
      </section>
    </header>
  );
};

const AppFooter = () => {
  const layoutClassName = useLayoutElementClassName();

  return (
    <footer className={cn("flex justify-center", layoutClassName)}>
      <Outlet name="footer">
        <div className={AppContainerClassName}>
          <AppInfoSection />
        </div>
      </Outlet>
    </footer>
  );
};

const AppMain = () => {
  const layoutClassName = useLayoutElementClassName();

  return (
    <main className={cn("pt-16 grow flex justify-center", layoutClassName)}>
      <div className={AppContainerClassName}>
        <Outlet />
      </div>
    </main>
  );
};

const AppAside = () => {
  const { asideMode } = useAppContext();
  const dispatch = useAppDispatch();

  const isWide = useMedia("(min-width: 48rem)");

  return (
    <aside
      className={cn(
        "fixed h-full z-20 inset-s-0 flex",
        "max-md:w-full max-md:bg-black/10 max-md:backdrop-blur-xs",
        asideMode === AppAsideMode.Hidden && "hidden",
        asideMode === AppAsideMode.Collapsed && "max-md:hidden",
      )}
      onClick={() =>
        dispatch?.({
          asideMode: isWide ? asideMode : AppAsideMode.Collapsed,
          subjectVisible: false,
        })
      }
    >
      <div
        className="bg-sidebar text-sidebar-foreground border-sidebar-border border-0 border-r h-full flex flex-row overflow-hidden"
        onMouseLeave={() => dispatch?.({ subjectVisible: false })}
      >
        <AppNavigationSection />
        <AppSubjectNavigationSection />
      </div>
    </aside>
  );
};

const AppNavigationSection = () => {
  const { asideMode } = useAppContext();

  const collapsedClassName = useMemo(
    () => (asideMode === AppAsideMode.Collapsed ? "collapsed" : ""),
    [asideMode],
  );

  return (
    <section
      className={cn(
        // "max-md:hidden",
        asideMode === AppAsideMode.Collapsed ? "w-14" : "w-3xs",
        "h-full flex flex-col justify-between overflow-auto scrollbar-none p-2",
      )}
    >
      <AppNavigation className={collapsedClassName} />
      <AppSecondaryNavigation className={collapsedClassName} />
    </section>
  );
};

const AppSubjectNavigationSection = () => {
  const { subject, subjectVisible } = useAppContext();

  const subjectNavigation = useMemo(() => {
    if (subject === AppSubject.Chat) {
      return <ChatNavigation />;
    }
    if (subject === AppSubject.Clusters) {
      return <ClustersNavigation />;
    }
    if (subject === AppSubject.Apps) {
      return <AppsNavigation />;
    }
    return null;
  }, [subject]);
  return (
    <section
      className={cn(
        "flex flex-col w-56 overflow-y-auto border-sidebar-primary border-0 border-l-2 p-2",
        "max-md:hidden",
        subjectVisible ? "block" : "hidden",
      )}
    >
      {subjectNavigation}
    </section>
  );
};

const AppInfoSection = () => {
  return (
    <section className="py-3 font-medium text-sm flex items-center justify-center gap-1 text-muted-foreground">
      <p>
        Power by
        <span className="pl-1 uppercase">{import.meta.env.VITE_APP_NAME}</span>
      </p>
      <p>|</p>
      <p>Version: {import.meta.env.VITE_APP_VERSION}</p>
    </section>
  );
};

export const AppNotFound = () => {
  const { t } = useTranslation();

  return (
    <PageInfoSection
      logo={<BlockOutlineIcon className="w-26" />}
      title={t("pages.notFound.title")}
      description={t("pages.notFound.description")}
    />
  );
};

const NotFoundOutlets = () => (
  <Outlets>
    <AppNotFound />
  </Outlets>
);
