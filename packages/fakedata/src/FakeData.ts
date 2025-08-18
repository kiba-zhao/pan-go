import { newAppSettings } from "./data/AppSettings";
import type { SeedOptions } from "./data/Common";
import { seedDevice } from "./data/Device";
import { seedDeviceStorage, toDeviceStorageFile } from "./data/DeviceStorage";
import { seedDeviceStorageFile } from "./data/DeviceStorageFile";
import { seedDeviceStorageSearchFile } from "./data/DeviceStorageSearchFile";
import { seedStorage, toStorageFile } from "./data/Storage";
import { seedStorageFile } from "./data/StorageFile";
import { seedStorageSearchFile } from "./data/StorageSearchFile";
import { seedStorageSearchItem } from "./data/StorageSearchItem";

export function seed(opts?: SeedOptions, isMobile?: boolean) {
  const appSettings = newAppSettings();

  const devices = seedDevice(opts);
  if (isMobile)
    return {
      appSettings,
      devices,
      devicesStorages: [],
      devicesStorageFiles: [],
      storages: [],
      storageFiles: [],
      storageSearchItems: [],
      storageSearchFiles: [],
      deviceStorageSearchFiles: [],
    };

  const devicesStorages = seedDeviceStorage(devices, opts);
  const devicesStorageFiles = seedDeviceStorageFile(devicesStorages, opts);

  const storages = seedStorage(opts);
  const storageFiles = seedStorageFile(storages, opts);

  const storageSearchItems = seedStorageSearchItem(opts);
  const storageSearchFiles = seedStorageSearchFile(
    [...storages.map(toStorageFile), ...storageFiles],
    opts
  );
  const deviceStorageSearchFiles = seedDeviceStorageSearchFile(
    [...devicesStorages.map(toDeviceStorageFile), ...devicesStorageFiles],
    opts
  );

  return {
    appSettings,
    devices,
    devicesStorages,
    devicesStorageFiles,
    storages,
    storageFiles,
    storageSearchItems,
    storageSearchFiles,
    deviceStorageSearchFiles,
  };
}
