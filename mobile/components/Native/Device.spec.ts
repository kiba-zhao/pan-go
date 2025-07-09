import type {
  Device,
  SearchQueryFields,
  SearchRangeFields,
  SearchResults,
  SearchSortFields,
} from '@pango/datatype';
import {exec} from './Agent.module';

type SearchCondition = Partial<Pick<Device, 'blocked' | 'online'>> &
  SearchQueryFields<string> &
  SearchRangeFields &
  SearchSortFields<'updatedAt' | 'createdAt' | 'blocked'>;
export function searchDevices(
  condition: SearchCondition,
): Promise<SearchResults<Device>> {
  if (__DEV__) {
    const {mockEnabled, search, readData} = require('../Common/FakeData');

    if (mockEnabled()) {
      return search(
        (entity: Device): boolean => {
          if (
            condition.blocked !== void 0 &&
            entity.blocked !== condition.blocked
          )
            return false;
          if (condition.online !== void 0 && entity.online !== condition.online)
            return false;
          if (condition.q !== void 0 && entity.name.indexOf(condition.q) < 0)
            return false;
          return true;
        },
        condition,
        readData('devices'),
      );
    }
  }

  return exec<SearchCondition, SearchResults<Device>>(
    'search_devices',
    condition,
  );
}
