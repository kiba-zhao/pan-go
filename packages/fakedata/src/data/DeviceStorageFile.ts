import { faker } from "@faker-js/faker";
import type { DeviceStorage, DeviceStorageFile } from "@pango/data";
import type { SeedOptions } from "./Common";
import { newMimeType } from "./MimeType";

export function seedDeviceStorageFile(
  storages: DeviceStorage[],
  opts?: SeedOptions
) {
  const folders: DeviceStorageFolder[] = storages
    .filter((storage) => storage.fileType === "D" && storage.available)
    .map((_) => ({ peerId: _.peerId, storageId: _.storageId, filePath: "" }));
  if (folders.length <= 0) return [];

  return faker.helpers.multiple(() => newDeviceStorageFile(folders), opts);
}

type DeviceStorageFolder = Pick<
  DeviceStorageFile,
  "peerId" | "storageId" | "filePath"
>;
function newDeviceStorageFile(folders: DeviceStorageFolder[]) {
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
    folders.push({
      peerId: folder.peerId,
      storageId: folder.storageId,
      filePath,
    });
  }

  return {
    id: faker.string.nanoid(),
    peerId: folder.peerId,
    storageId: folder.storageId,
    name,
    parentPath: folder.filePath,
    filePath,
    fileType,
    mimeType,
    size: faker.number.int({ min: 1, max: 999999 }),
    available: true,
    createdAt: faker.date.past().toString(),
    updatedAt: faker.date.past().toString(),
  } as DeviceStorageFile;
}
