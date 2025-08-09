import { useQuery } from '../Common/ReactQuery';
import { load } from '../Spec/AppSettings';

export const QueryKey = ['settings'];

export  const useSettings = () => useQuery({
    queryKey: QueryKey,
    queryFn: load,
});