import { faker } from "@faker-js/faker";
import type { StorageFile, StorageSearchFile } from "@pango/data";
import type { SeedOptions } from "./Common";

export function seedStorageSearchFile(
  storageFiles: StorageFile[],
  opts?: SeedOptions
) {
  return faker.helpers.multiple(() => newStorageSearchFile(storageFiles), opts);
}

function newStorageSearchFile(storageFiles: StorageFile[]) {
  const storageFile = faker.helpers.arrayElement(storageFiles);
  return {
    id: faker.number.int({ min: 1, max: 999999 }),
    storageId: storageFile.storageId,
    name: storageFile.name,
    fileType: storageFile.fileType,
    size: storageFile.size,
    available: storageFile.available,
    filePath: storageFile.filePath,
    mimeType: storageFile.mimeType,
    score: faker.number.int({ min: 1, max: 10 }),
    tokens: [storageFile.name],
    createdAt: faker.date.past().toString(),
    updatedAt: faker.date.past().toString(),
  } as StorageSearchFile;
}
