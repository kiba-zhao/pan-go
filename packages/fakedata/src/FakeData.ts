import { newSettings, newHostSettings } from "./data/settings";
export function seed() {
  return {
    settings: newSettings(),
    hostSettings: newHostSettings(),
  };

  // const devices = seedDevice(opts);
  // const devicesStorages = seedDeviceStorage(devices, opts);
  // const devicesStorageFiles = seedDeviceStorageFile(devicesStorages, opts);

  // const storages = seedStorage(opts);
  // const storageFiles = seedStorageFile(storages, opts);

  // const storageSearchItems = seedStorageSearchItem(opts);
  // const storageSearchFiles = seedStorageSearchFile(
  //   [...storages.map(toStorageFile), ...storageFiles],
  //   opts,
  // );
  // const deviceStorageSearchFiles = seedDeviceStorageSearchFile(
  //   [...devicesStorages.map(toDeviceStorageFile), ...devicesStorageFiles],
  //   opts,
  // );
}
