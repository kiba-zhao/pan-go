import type {
  InfiniteData,
  MutationFunction,
  QueryKey,
} from '@tanstack/react-query';
import {
  QueryClient,
  QueryClientProvider,
  QueryObserver,
  useInfiniteQuery,
  useIsFetching,
  useMutation,
  usePrefetchQuery,
  useQuery,
  useQueryClient,
} from '@tanstack/react-query';

export const ReactQueryProvider = ({children}: {children: React.ReactNode}) => {
  const queryClient = new QueryClient();
  return (
    <QueryClientProvider client={queryClient}>{children}</QueryClientProvider>
  );
};

export {
  QueryObserver,
  useInfiniteQuery,
  useIsFetching,
  useMutation,
  usePrefetchQuery,
  useQuery,
  useQueryClient,
};
export type {InfiniteData, MutationFunction, QueryClient, QueryKey};
