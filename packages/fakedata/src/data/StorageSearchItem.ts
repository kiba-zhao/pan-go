import { faker } from "@faker-js/faker";
import type { StorageSearchItem } from "@pango/data";
import type { SeedOptions } from "./Common";

export function seedStorageSearchItem(opts?: SeedOptions) {
  return faker.helpers.multiple(newStorageSearchItem, opts);
}

function newStorageSearchItem() {
  return {
    id: faker.number.int({ min: 1, max: 999999 }),
    query: faker.word.words(),
    createdAt: faker.date.past().toString(),
    updatedAt: faker.date.past().toString(),
  } as StorageSearchItem;
}
