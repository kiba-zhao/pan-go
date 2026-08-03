import {
  createContext,
  type Dispatch,
  useReducer,
  useContext,
  type PropsWithChildren,
} from "react";

export enum AppSubject {
  Chat,
  Clusters,
  Apps,
}

export enum AppAsideMode {
  Collapsed = "collapsed",
  Hidden = "hidden",
}

type AppState = {
  subject?: AppSubject;
  subjectVisible?: boolean;
  asideMode?: AppAsideMode;
};

type AppHeaderState = {
  title?: string;
  breadcrumbs?: Array<[string, string]>;
};

type AppExtraState = any;

export type AppContextState = AppState & {
  header?: AppHeaderState;
  extra?: AppExtraState;
};

type AppContextAction = AppContextState;
const AppContextReducer = (
  state: AppContextState,
  action: AppContextAction,
) => {
  return { ...state, ...action };
};

const context = createContext<AppContextState>({});
const dispatchContext = createContext<Dispatch<AppContextAction> | null>(null);

export const useAppContext = () => useContext(context);
export const useAppDispatch = () => useContext(dispatchContext);
export const useAppHeader = () => useContext(context)?.header || {};
export const useAppExtra = <T extends unknown>(defaultValue: T) => {
  const extra = useContext(context)?.extra;
  if (extra === void 0 || typeof extra != typeof defaultValue) {
    return defaultValue;
  }
  return extra as T;
};

export const AppContextProvider = ({ children }: PropsWithChildren) => {
  const [state, dispatch] = useReducer(AppContextReducer, {});
  return (
    <context.Provider value={state}>
      <dispatchContext.Provider value={dispatch}>
        {children}
      </dispatchContext.Provider>
    </context.Provider>
  );
};

export function resetToBlank() {
  return { header: void 0, extra: void 0 };
}
