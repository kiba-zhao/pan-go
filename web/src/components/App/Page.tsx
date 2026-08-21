import "./Theme.css";
import "./Page.css";

import { useMedia } from "@/lib/hooks";
import { Router } from "./Router";
import { Outlet } from "./Outlets";
import { useMemo } from "react";
import { cn } from "@/lib/utils";
import {
  AppContextProvider,
  useAppContext,
  useAppDispatch,
  AppSubject,
  AppAsideMode,
  withAppAction,
} from "./Context";
import { AsideModeControl, Separator } from "./Tools";
import AppBreadcrumb from "./Breadcrumb";

import { AppNavigation, AppSecondaryNavigation } from "./PageNavigation";
import ChatNavigation from "@/components/Chat/Navigation";
import AppsNavigation from "@/components/Apps/Navigation";

import { Toaster } from "./Toast";
import { default as PageRoutes } from "./PageRoutes";

import { PageInfoSection } from "./Layout";
import { BlocksShuffle3Icon } from "./Icon";
import { useTranslation, I18nVariant } from "./I18Next";

const AppPage = () => {
  return (
    <AppContextProvider>
      <Router>
        <PageRoutes layout={<AppLayout />} />
      </Router>
    </AppContextProvider>
  );
};

export default AppPage;

type AppLayoutProps = {
  headerClassName?: string;
  mainClassName?: string;
  footerClassName?: string;
};
const AppLayout = ({
  headerClassName,
  mainClassName,
  footerClassName,
}: AppLayoutProps) => {
  return (
    <>
      <AppHeader className={headerClassName} />
      <AppMain className={mainClassName} />
      <AppFooter className={footerClassName} />
      <AppAside />
      <Outlet name="extra" />
      <Toaster />
    </>
  );
};

const useLayoutElementClassName = (className?: string) => {
  const { asideMode } = useAppContext();
  if (asideMode === AppAsideMode.Collapsed) {
    return cn("md:pl-14", className);
  }

  if (asideMode === AppAsideMode.Hidden) {
    return className || "";
  }

  return cn("md:pl-(--container-3xs)", className);
};

const AppContainerClassName = "max-w-7xl w-full";
const AppHeader = ({ className }: { className?: string }) => {
  const layoutClassName = useLayoutElementClassName(className);

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

const AppFooter = ({ className }: { className?: string }) => {
  const layoutClassName = useLayoutElementClassName(className);
  return (
    <footer className={cn("flex justify-center", layoutClassName)}>
      <Outlet name="footer" fallback={<AppFooterContent />}>
        <AppFooterContent />
      </Outlet>
    </footer>
  );
};

const AppFooterContent = () => {
  return (
    <div className={AppContainerClassName}>
      <AppInfoSection />
    </div>
  );
};

const AppMain = ({ className }: { className?: string }) => {
  const layoutClassName = useLayoutElementClassName(className);

  return (
    <main className={cn("pt-16 grow flex justify-center", layoutClassName)}>
      <div className={AppContainerClassName}>
        <Outlet fallback={<AppMainLoading />} />
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
        dispatch?.(
          withAppAction({
            asideMode: isWide ? asideMode : AppAsideMode.Collapsed,
            subjectVisible: false,
          }),
        )
      }
    >
      <div
        className="bg-sidebar text-sidebar-foreground border-sidebar-border border-0 border-r h-full flex flex-row overflow-hidden"
        onMouseLeave={() =>
          dispatch?.(withAppAction({ subjectVisible: false }))
        }
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

const AppMainLoading = () => {
  const { t } = useTranslation();
  return (
    <PageInfoSection
      logo={<BlocksShuffle3Icon className="w-26" />}
      title={t(`${I18nVariant.Main}.loading.title`)}
      description={t(`${I18nVariant.Main}.loading.description`)}
    />
  );
};
