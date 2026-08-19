import {
  createContext,
  type Dispatch,
  useReducer,
  useContext,
  type PropsWithChildren,
} from "react";

export enum AppSubject {
  Chat,
  Apps,
}

export enum AppAsideMode {
  Collapsed = "collapsed",
  Hidden = "hidden",
}

export type AppState = {
  subject?: AppSubject;
  subjectVisible?: boolean;
  asideMode?: AppAsideMode;
};

type AppHeaderState = {
  title?: string;
  breadcrumbs?: Array<[string, string]>;
};

type AppExtraState = Record<string, unknown>;

type AppContextState = {
  app: AppState;
  header: AppHeaderState;
  extra: AppExtraState;
};

type AppContextStateKey = keyof AppContextState;
type AppContextAction =
  | {
      type: AppContextStateKey;
      payload?: AppContextState[AppContextStateKey];
    }
  | AppContextState;

function AppContextReducer(state: AppContextState, action: AppContextAction) {
  const { type, payload } = action as Exclude<
    AppContextAction,
    AppContextState
  >;
  if (type !== void 0) {
    return {
      ...state,
      [type]: payload === void 0 ? payload : { ...state[type], ...payload },
    };
  }
  return { ...state, ...action };
}

const appContext = createContext<AppState>({});
const appHeaderContext = createContext<AppHeaderState>({});
const appExtraContext = createContext<AppExtraState>({});
const dispatchContext = createContext<Dispatch<AppContextAction> | null>(null);

export const useAppContext = () => useContext(appContext);
export const useAppDispatch = () => useContext(dispatchContext);
export const useAppHeader = () => useContext(appHeaderContext);
export const useAppExtra = <T extends AppExtraState>() =>
  useContext(appExtraContext) as T;

export function withAppAction(payload?: AppState): AppContextAction {
  return { type: "app", payload };
}
export function withAppHeaderAction(
  payload?: AppHeaderState,
): AppContextAction {
  return { type: "header", payload };
}
export function withAppExtraAction(payload?: AppExtraState): AppContextAction {
  return { type: "extra", payload };
}

export const AppContextProvider = ({ children }: PropsWithChildren) => {
  const [state, dispatch] = useReducer(AppContextReducer, {
    app: {},
    header: {},
    extra: {},
  });
  return (
    <appContext.Provider value={state.app}>
      <appHeaderContext.Provider value={state.header}>
        <appExtraContext.Provider value={state.extra}>
          <dispatchContext.Provider value={dispatch}>
            {children}
          </dispatchContext.Provider>
        </appExtraContext.Provider>
      </appHeaderContext.Provider>
    </appContext.Provider>
  );
};

export function withResetAction() {
  return { header: {}, extra: {} } as AppContextState;
}
