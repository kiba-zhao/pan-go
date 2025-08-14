import { faker } from "@faker-js/faker";
import type { Storage, StorageFile } from "@pango/data";
import type { SeedOptions } from "./Common";
import { newMimeType } from "./MimeType";

export function seedStorage(opts?: SeedOptions) {
  return faker.helpers.multiple(newStorage, opts);
}

function newStorage() {
  const fileType = faker.helpers.arrayElement(["F", "D"]);
  const mimeType = fileType === "F" ? newMimeType() : "";
  const extname = fileType === "F" ? faker.system.fileExt(mimeType) : "";
  const name =
    fileType === "F"
      ? faker.system.commonFileName(extname)
      : faker.system.fileName({ extensionCount: 0 });

  return {
    id: faker.number.int({ min: 1, max: 999999 }),
    name,
    filePath: faker.system.directoryPath(),
    fileType,
    mimeType,
    size: faker.number.int({ min: 1, max: 999999 }),
    enabled: faker.datatype.boolean(),
    available: faker.datatype.boolean(),
    createdAt: faker.date.past().toString(),
    updatedAt: faker.date.past().toString(),
  } as Storage;
}

export function toStorageFile(storage: Storage) {
  return {
    id: faker.string.nanoid(),
    storageId: storage.id,
    name: storage.name,
    filePath: "",
    fileType: storage.fileType,
    mimeType: storage.mimeType,
    parentPath: "",
    size: storage.size,
    available: storage.available,
    createdAt: storage.createdAt,
    updatedAt: storage.updatedAt,
  } as StorageFile;
}
