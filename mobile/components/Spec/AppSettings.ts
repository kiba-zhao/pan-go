import type {AppSettings, AppSettingsFields} from '@pango/data';
import {exec} from '../Native/AgentModule';

export async function load(): Promise<AppSettings> {
  if (__DEV__) {
    const {mockEnabled, readData} = require('../Common/FakeData');

    if (mockEnabled()) {
      return readData('appSettings');
    }
  }
  return await exec('load:appSettings');
}

export async function save(fields: AppSettingsFields): Promise<AppSettings> {
  if (__DEV__) {
    const {mockEnabled, readData, writeData} = require('../Common/FakeData');

    if (mockEnabled()) {
      const settings = {...readData('appSettings'), ...fields};
      writeData('appSettings', settings);
      return settings;
    }
  }
  return await exec('save:appSettings', fields);
}
