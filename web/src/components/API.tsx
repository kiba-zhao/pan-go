import type { ReactNode } from "react";
import { createContext, useContext } from "react";

import type { API } from "../api";
import * as api from "../api";

const APIContext = createContext<API | null>(null);

export const useAPI = () => useContext(APIContext);

export const APIProvider = ({ children }: { children: ReactNode }) => (
  <APIContext.Provider value={api}>{children}</APIContext.Provider>
);
