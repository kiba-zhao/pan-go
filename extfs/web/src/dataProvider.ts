import simpleRestProvider from "ra-data-simple-rest";

const ROOT_PATH = `${import.meta.env.BASE_URL}${
  import.meta.env.VITE_API_PATH || ""
}`;

export const dataProvider = simpleRestProvider(ROOT_PATH);
