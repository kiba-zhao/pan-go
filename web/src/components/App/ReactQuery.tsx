export {
  QueryClient,
  QueryClientProvider,
  useQuery,
  useMutation,
  useQueryClient,
} from "@tanstack/react-query";

import type { UseQueryOptions, MutationOptions } from "@tanstack/react-query";
export type { UseQueryOptions, MutationOptions };

export type UseQueryOpts<TQueryFnData, TData> = Partial<
  Omit<UseQueryOptions<TQueryFnData, Error, TData>, "queryFn">
>;

export type UseMutationOpts<TData, TVariables = Partial<TData>> = Omit<
  MutationOptions<TData, Error, TVariables>,
  "mutationFn"
>;
