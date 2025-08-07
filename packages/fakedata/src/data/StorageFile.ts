import { faker } from "@faker-js/faker";
import type { Storage, StorageFile } from "@pango/data";
import type { SeedOptions } from "./Common";
import { newMimeType } from "./MimeType";

export function seedStorageFile(storages: Storage[], opts?: SeedOptions) {
  const folders: StorageFolder[] = storages
    .filter((storage) => storage.fileType === "D" && storage.available)
    .map((_) => ({ storageId: _.id, filePath: "" }));
  if (folders.length <= 0) return [];

  return faker.helpers.multiple(() => newStorageFile(folders), opts);
}

type StorageFolder = Pick<StorageFile, "storageId" | "filePath">;
function newStorageFile(folders: StorageFolder[]) {
  const isDir = faker.datatype.boolean();
  const mimeType = isDir ? "" : newMimeType();
  const folder = faker.helpers.arrayElement(folders);
  const fileType = isDir ? "D" : "F";
  const extname = isDir ? "" : faker.system.fileExt(mimeType);
  const name = isDir
    ? faker.system.fileName({ extensionCount: 0 })
    : faker.system.commonFileName(extname);
  const filePath =
    folder.filePath.length > 0 ? `${folder.filePath}/${name}` : name;

  if (isDir) {
    folders.push({ storageId: folder.storageId, filePath } as StorageFolder);
  }

  return {
    id: faker.string.nanoid(),
    storageId: folder.storageId,
    name,
    filePath,
    fileType,
    mimeType,
    parentPath: folder.filePath,
    size: faker.number.int({ min: 1, max: 999999 }),
    available: true,
    createdAt: faker.date.past().toString(),
    updatedAt: faker.date.past().toString(),
  } as StorageFile;
}
