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

type AppExtraState = Exclude<any, undefined>;

export type AppContextState = AppState & {
  header?: AppHeaderState;
  extra?: AppExtraState;
};

type AppContextBlockKey = keyof Omit<AppContextState, keyof AppState>;
type AppContextBlockAction<
  BlockKey extends AppContextBlockKey = AppContextBlockKey,
> = {
  type: "block";
  blockKey: BlockKey;
} & Pick<AppContextState, BlockKey>;
type AppContextAction =
  | AppContextState
  | AppContextBlockAction<AppContextBlockKey>;
const AppContextReducer = (
  state: AppContextState,
  action: AppContextAction,
) => {
  const blockAction = action as AppContextBlockAction<AppContextBlockKey>;
  if (!blockAction.type) {
    return { ...state, ...action };
  }

  if (blockAction.type === "block") {
    if (blockAction.blockKey === "header") {
      return {
        ...state,
        header:
          blockAction.header === void 0
            ? void 0
            : { ...state.header, ...blockAction.header },
      };
    }
    if (blockAction.blockKey === "extra") {
      return {
        ...state,
        extra:
          blockAction.extra === void 0
            ? void 0
            : { ...state.extra, ...blockAction.extra },
      };
    }
  }

  return state;
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

export function withAppHeaderAction(
  payload?: AppHeaderState,
): AppContextBlockAction<"header"> {
  return { type: "block", blockKey: "header", header: payload };
}
export function withAppExtraAction(
  payload?: AppExtraState,
): AppContextBlockAction<"extra"> {
  return { type: "block", blockKey: "extra", extra: payload };
}

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
