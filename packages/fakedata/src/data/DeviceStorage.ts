import { faker } from "@faker-js/faker";
import type { Device, DeviceStorage, DeviceStorageFile } from "@pango/data";
import type { SeedOptions } from "./Common";
import { newMimeType } from "./MimeType";

export function seedDeviceStorage(devices: Device[], opts?: SeedOptions) {
  const onlineDevices = devices.filter((device) => device.online);
  if (onlineDevices.length <= 0) return [];
  return faker.helpers.multiple(() => newDeviceStorage(onlineDevices), opts);
}

function newDeviceStorage(devices: Device[]) {
  const device = faker.helpers.arrayElement(devices);
  const fileType = faker.helpers.arrayElement(["F", "D"]);
  const mimeType = fileType === "F" ? newMimeType() : "";
  const extname = fileType === "F" ? faker.system.fileExt(mimeType) : "";
  const name =
    fileType === "F"
      ? faker.system.commonFileName(extname)
      : faker.system.fileName({ extensionCount: 0 });
  return {
    id: faker.string.nanoid(),
    peerId: device.peerId,
    storageId: faker.number.int({ min: 1, max: 999999 }),
    name,
    fileType,
    mimeType,
    size: faker.number.int({ min: 1, max: 999999 }),
    available: faker.datatype.boolean(),
    createdAt: faker.date.past().toString(),
    updatedAt: faker.date.past().toString(),
  } as DeviceStorage;
}

export function toDeviceStorageFile(storage: DeviceStorage) {
  return {
    id: storage.id,
    peerId: storage.peerId,
    storageId: storage.storageId,
    name: storage.name,
    filePath: "",
    fileType: storage.fileType,
    mimeType: storage.mimeType,
    parentPath: "",
    size: storage.size,
    available: storage.available,
    createdAt: storage.createdAt,
    updatedAt: storage.updatedAt,
  } as DeviceStorageFile;
}
