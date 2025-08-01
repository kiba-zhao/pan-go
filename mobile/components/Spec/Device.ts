import type {
  Device,
  DeviceFields,
  SearchQueryFields,
  SearchRangeFields,
  SearchResults,
  SearchSortFields,
} from '@pango/data';
import {exec} from '../Native/AgentModule';

export type SearchCondition = Partial<Pick<Device, 'enabled' | 'online'>> &
  SearchQueryFields<string> &
  SearchRangeFields &
  SearchSortFields<'updatedAt' | 'createdAt' | 'enabled'>;
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
