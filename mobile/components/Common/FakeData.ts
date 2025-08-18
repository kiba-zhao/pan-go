import {
  create,
  destroy,
  fetch,
  findById,
  findOne,
  newNanoID,
  nextID,
  replace,
  search,
  seed,
  sort,
  update,
} from '@pango/fakedata';
import config from 'react-native-config';

type FakeData = ReturnType<typeof seed>;
const fakeData = seed({count: 20}, true);

export function readData<
  Key extends keyof FakeData,
  Value extends FakeData[Key],
>(name: Key): Value {
  return fakeData[name] as Value;
}

export function writeData<
  Key extends keyof FakeData,
  Value extends FakeData[Key],
>(name: Key, value: Value): void {
  fakeData[name] = value;
}

export {
  create,
  destroy,
  fetch,
  findById,
  findOne,
  newNanoID,
  nextID,
  replace,
  search,
  sort,
  update,
};

export function mockEnabled(): boolean {
  return config.MOCK_ENABLED === 'true';
}
