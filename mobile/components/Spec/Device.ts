import type {
  ActionFields,
  CursorFields,
  Device,
  DeviceFields,
  FetchResults,
  FetchResultsMeta,
  QFields,
  RangeFields,
  SearchResults,
  SortFields,
} from '@pango/data';
import {exec} from '../Native/AgentModule';

export type SearchCondition = Partial<Pick<Device, 'enabled' | 'online'>> &
  QFields<string> &
  RangeFields &
  SortFields<'updatedAt' | 'createdAt' | 'enabled'>;
export async function searchDevices(
  condition: SearchCondition,
): Promise<SearchResults<Device>> {
  if (__DEV__) {
    const {mockEnabled, search, readData} = require('../Common/FakeData');

    if (mockEnabled()) {
      return search(
        (entity: Device): boolean => {
          if (
            condition.enabled !== void 0 &&
            entity.enabled !== condition.enabled
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

  return await exec<SearchCondition, SearchResults<Device>>(
    'search:devices',
    condition,
  );
}

export async function selectDevice(id: Device['id']): Promise<Device> {
  if (__DEV__) {
    const {mockEnabled, findById, readData} = require('../Common/FakeData');

    if (mockEnabled()) {
      return findById(id, readData('devices'));
    }
  }
  return await exec('select:devices', {id});
}

export async function updateDevice(
  fields: Partial<DeviceFields>,
  id: Device['id'],
): Promise<Device> {
  if (__DEV__) {
    const {mockEnabled, update, readData} = require('../Common/FakeData');

    if (mockEnabled()) {
      return update(id, fields, readData('devices'));
    }
  }
  return await exec(`update:devices`, {...fields, id});
}

export async function destroyDevice(id: Device['id']) {
  if (__DEV__) {
    const {mockEnabled, destroy, readData} = require('../Common/FakeData');

    if (mockEnabled()) {
      return destroy(id, readData('devices'));
    }
  }
  return await exec('destroy:devices', {id});
}

export async function selectDeviceByName(name: string): Promise<Device> {
  if (__DEV__) {
    const {mockEnabled, findOne, readData} = require('../Common/FakeData');

    if (mockEnabled()) {
      return findOne(
        (entity: Device): boolean => entity.name === name,
        readData('devices'),
      );
    }
  }
  return await exec('select:devices?name', {name});
}

export async function createDevice(fields: DeviceFields): Promise<Device> {
  if (__DEV__) {
    const {
      mockEnabled,
      nextID,
      create,
      readData,
    } = require('../Common/FakeData');

    if (mockEnabled()) {
      return create(fields, readData('devices'), nextID);
    }
  }
  return await exec('create:devices', fields);
}

type FetchBaseCondition = Partial<Pick<Device, 'enabled' | 'online'>> &
  QFields<string> &
  ActionFields<'search'>;
type FetchCursorCondition = CursorFields<string> & FetchBaseCondition;
export type FetchCondition = FetchCursorCondition;

export type FetchDevicesResults = FetchResults<Device, string, string>;

export async function fetchDevices(
  condition: FetchCondition,
): Promise<FetchDevicesResults> {
  if (__DEV__) {
    const {mockEnabled, sort, readData} = require('../Common/FakeData');
    if (mockEnabled()) {
      const data = readData('devices').filter((entity: Device) => {
        if (
          condition.enabled !== void 0 &&
          entity.enabled !== condition.enabled
        )
          return false;
        if (condition.online !== void 0 && entity.online !== condition.online)
          return false;
        if (condition.q !== void 0 && entity.name.indexOf(condition.q) < 0)
          return false;
        return true;
      });
      const sorted: Device[] = sort(
        {sortField: 'updatedAt', order: 'desc'},
        data,
      );

      if (sorted.length <= 0) {
        return [{tag: '', offset: 0}, []];
      }

      const tag = (sorted.at(0) as Device).updatedAt;
      const {_cursor, _limit} = condition as FetchCursorCondition;
      if (_limit === 0) {
        return [{tag, offset: 0}, []];
      }

      let start_ = -1;
      let end_ = sorted.length;

      if (_cursor !== void 0) {
        start_ = sorted.findIndex(entity => entity.updatedAt === _cursor);
        if (start_ < 0) {
          return [{tag, offset: start_}, []];
        }
      }

      if (_limit !== void 0) {
        if (_limit > 0) {
          if (start_ >= end_ - 1) return [{tag, offset: end_}, []];
          start_++;
          end_ = start_ + _limit;
          if (end_ >= sorted.length) end_ = sorted.length;
        } else {
          if (start_ === 0) return [{tag, offset: start_}, []];
          if (start_ > 0) end_ = start_;
          start_ = end_ + _limit;
          if (start_ < 0) start_ = 0;
        }
      }

      const entities = sorted.slice(start_, end_);
      const meta = {tag, offset: start_} as FetchResultsMeta<string, string>;
      meta.prev = entities.at(0)?.updatedAt;
      meta.next = entities.at(-1)?.updatedAt;
      if (start_ === 0) {
        meta.prev = void 0;
      }
      if (end_ === sorted.length) {
        meta.next = void 0;
      }
      return [meta, entities];
    }
  }
  return await exec<FetchCondition, FetchResults<Device, string, string>>(
    'fetch:devices',
    condition,
  );
}
