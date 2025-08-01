import type {MutationFunction} from '@tanstack/react-query';
import {
  QueryClient,
  QueryClientProvider,
  useIsFetching,
  useMutation,
  useQuery,
  useQueryClient,
} from '@tanstack/react-query';

export const ReactQueryProvider = ({children}: {children: React.ReactNode}) => {
  const queryClient = new QueryClient();
  return (
    <QueryClientProvider client={queryClient}>{children}</QueryClientProvider>
  );
};

export {useIsFetching, useMutation, useQuery, useQueryClient};
export type {MutationFunction, QueryClient};
