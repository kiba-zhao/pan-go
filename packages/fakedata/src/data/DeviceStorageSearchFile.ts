import { faker } from "@faker-js/faker";
import type { DeviceStorageFile, DeviceStorageSearchFile } from "@pango/data";
import type { SeedOptions } from "./Common";

export function seedDeviceStorageSearchFile(
  storageFiles: DeviceStorageFile[],
  opts?: SeedOptions
) {
  return faker.helpers.multiple(
    () => newDeviceStorageSearchFile(storageFiles),
    opts
  );
}

function newDeviceStorageSearchFile(storageFiles: DeviceStorageFile[]) {
  const storageFile = faker.helpers.arrayElement(storageFiles);
  return {
    id: faker.number.int({ min: 1, max: 999999 }),
    peerId: storageFile.peerId,
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
  } as DeviceStorageSearchFile;
}
