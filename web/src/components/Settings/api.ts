import type {
  DeviceInfo,
  DeviceInfoFields,
  DeviceNetwork,
  DeviceNetworkFields,
  WebHost,
  WebHostFields,
} from "@pango/data";
import { BASE_PATH, toJson, withJSONBody } from "@/lib/fetch";
import { SettingsName } from "./meta";

export type {
  DeviceInfo,
  DeviceInfoFields,
  DeviceNetwork,
  DeviceNetworkFields,
  WebHost,
  WebHostFields,
};

let fakeSettings: Promise<{
  deviceInfo: DeviceInfo;
  deviceNetwork: DeviceNetwork;
  webHost: WebHost;
}>;
if (import.meta.env.VITE_FAKE_DATA) {
  fakeSettings = import("@pango/fakedata").then((_) => ({
    deviceInfo: _.newDeviceInfo(),
    deviceNetwork: _.newDeviceNetwork(),
    webHost: _.newWebHost(),
  }));
}

const DeviceInfoAPIPath = `${SettingsName}/device-info`;
const DeviceNetworkAPIPath = `${SettingsName}/device-network`;
const WebHostAPIPath = `${SettingsName}/web-host`;

export async function fetchDeviceInfo() {
  if (import.meta.env.VITE_FAKE_DATA) {
    const { deviceInfo } = await fakeSettings;
    return deviceInfo;
  }
  const res = await fetch(`${BASE_PATH}/${DeviceInfoAPIPath}`);
  return toJson<DeviceInfo>(res);
}

export async function saveDeviceInfo(fields: DeviceInfoFields) {
  if (import.meta.env.VITE_FAKE_DATA) {
    const { deviceInfo, ...others } = await fakeSettings;
    const newDeviceInfo = { ...deviceInfo, ...fields };
    fakeSettings = Promise.resolve({
      ...others,
      deviceInfo: newDeviceInfo,
    });
    return newDeviceInfo;
  }

  const res = await fetch(
    `${BASE_PATH}/${DeviceInfoAPIPath}`,
    withJSONBody(fields, { method: "PATCH" }),
  );
  return toJson<DeviceInfo>(res);
}

export async function fetchDeviceNetwork() {
  if (import.meta.env.VITE_FAKE_DATA) {
    const { deviceNetwork } = await fakeSettings;
    return deviceNetwork;
  }
  const res = await fetch(`${BASE_PATH}/${DeviceNetworkAPIPath}`);
  return toJson<DeviceNetwork>(res);
}

export async function saveDeviceNetwork(fields: DeviceNetworkFields) {
  if (import.meta.env.VITE_FAKE_DATA) {
    const { deviceNetwork, ...others } = await fakeSettings;
    const newDeviceNetwork = { ...deviceNetwork, ...fields };
    fakeSettings = Promise.resolve({
      ...others,
      deviceNetwork: newDeviceNetwork,
    });
    return newDeviceNetwork;
  }

  const res = await fetch(
    `${BASE_PATH}/${DeviceNetworkAPIPath}`,
    withJSONBody(fields, { method: "PATCH" }),
  );
  return toJson<DeviceNetwork>(res);
}

export async function fetchWebHost() {
  if (import.meta.env.VITE_FAKE_DATA) {
    const { webHost } = await fakeSettings;
    return webHost;
  }
  const res = await fetch(`${BASE_PATH}/${WebHostAPIPath}`);
  return toJson<WebHost>(res);
}

export async function saveWebHost(fields: WebHostFields) {
  if (import.meta.env.VITE_FAKE_DATA) {
    const { webHost, ...others } = await fakeSettings;
    const newWebHost = { ...webHost, ...fields };
    fakeSettings = Promise.resolve({
      ...others,
      webHost: newWebHost,
    });
    return newWebHost;
  }

  const res = await fetch(
    `${BASE_PATH}/${WebHostAPIPath}`,
    withJSONBody(fields, { method: "PATCH" }),
  );
  return toJson<WebHost>(res);
}
