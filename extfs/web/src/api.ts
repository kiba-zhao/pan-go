import { simple } from "fetch-utils";
import { simpleDataProvider } from "fetch-utils/raProvider";

const ROOT_PATH = `${import.meta.env.BASE_URL}${
  import.meta.env.VITE_API_PATH || ""
}`;

const { fetchOne, fetchMany } = simple(ROOT_PATH);

export const dataProvider = simpleDataProvider({ fetchOne, fetchMany });
