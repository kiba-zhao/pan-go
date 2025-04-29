import { withJSONBody, withMethod, withPath, withQuery } from "fetch-utils";
import { fetchOne, fetchMany } from "../../utils/api";

export type AppNode = {
  id: number;
  peerId: string;
  name: string;
  blocked: boolean;
  online: boolean;
  networkAddrs: string[] | null;
  createdAt: Date;
  updatedAt: Date;
};

export type AppNodeFields = Pick<
  AppNode,
  "peerId" | "name" | "blocked" | "networkAddrs"
>;

export type AppNodeSearchCondition = {
  q?: string;
  blocked?: boolean;
  online?: boolean;
};

export async function selectAllAppNodes(
  condition?: AppNodeSearchCondition
): Promise<AppNode[]> {
  const query = new URLSearchParams({
    q: condition?.q || "",
  });
  if (condition?.blocked !== void 0) {
    query.set("blocked", condition.blocked.toString());
  }
  if (condition?.online !== void 0) {
    query.set("online", condition.online.toString());
  }

  const [_, nodes] = await fetchMany(
    withPath("app/nodes", "merge"),
    withQuery(query, "merge")
  );

  return nodes;
}

export async function selectAppNode(id: AppNode["id"]): Promise<AppNode> {
  return await fetchOne(withPath(`app/nodes/${id}`, "merge"));
}

export async function saveAppNode(
  fields: AppNodeFields,
  id?: AppNode["id"]
): Promise<AppNode> {
  return await fetchOne(
    withPath(`app/nodes${id ? `/${id}` : ""}`, "merge"),
    withMethod(id ? "PATCH" : "POST"),
    withJSONBody(fields)
  );
}

export async function deleteAppNode(id: AppNode["id"]) {
  return await fetchOne(
    withPath(`app/nodes/${id}`, "merge"),
    withMethod("DELETE")
  );
}
