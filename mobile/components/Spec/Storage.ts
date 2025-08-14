import type {CursorFields, FetchResults, QFields, Storage} from '@pango/data';
import {exec} from '../Native/AgentModule';

export type FetchStoragesCondition = CursorFields<string> & QFields<string>;
export type FetchStoragesResults = FetchResults<Storage, string, string>;
export async function fetchStorages(
  condition: FetchStoragesCondition,
): Promise<FetchStoragesResults> {
  if (__DEV__) {
    const {mockEnabled, sort, readData, fetch} = require('../Common/FakeData');
    if (mockEnabled()) {
      const storages = readData('storages').filter((entity: Storage) => {
        if (
          condition.q !== void 0 &&
          condition.q.length > 0 &&
          entity.name.indexOf(condition.q) < 0
        )
          return false;
        return true;
      });
      const sorted: Storage[] = sort(
        {sortField: 'updatedAt', order: 'desc'},
        storages,
      );
      return fetch(condition, sorted, 'updatedAt', '') as FetchStoragesResults;
    }
  }
  return await exec('fetch:storages');
}
