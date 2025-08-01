import {Device} from '@pango/data';
import type {QueryClient} from '../Common/ReactQuery';

export const QueryKey = ['devices'];

export const RecentlyQueryKey = [
  ...QueryKey,
  {
    _sort: 'updatedAt',
    _order: 'desc',
    _start: 0,
    _end: 10,
  },
];

export const OnlineQueryKey = [
  ...QueryKey,
  {
    _sort: 'updatedAt',
    _order: 'desc',
    _start: 0,
    _end: 10,
    online: true,
  },
];

export const DisabledQueryKey = [
  ...QueryKey,
  {
    _sort: 'updatedAt',
    _order: 'desc',
    _start: 0,
    _end: 10,
    enabled: false,
  },
];

export function invalidateListQueryCache(
  queryClient: QueryClient,
  device: Device,
) {
  queryClient.invalidateQueries({
    queryKey: RecentlyQueryKey,
  });

  if (device.online) {
    queryClient.invalidateQueries({
      queryKey: OnlineQueryKey,
    });
  }

  if (!device.enabled) {
    queryClient.invalidateQueries({
      queryKey: DisabledQueryKey,
    });
  }
}
