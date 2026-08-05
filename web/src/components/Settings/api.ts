import type {
  Settings,
  HostSettings,
  SettingsFields,
  HostSettingsFields,
} from "@pango/data";
import { BASE_PATH, toJson, withJSONBody } from "@/lib/fetch";
import { SettingsName } from "./meta";

export type { Settings, HostSettings, SettingsFields, HostSettingsFields };

let fakeSettings: Promise<{ settings: Settings; hostSettings: HostSettings }>;
if (import.meta.env.VITE_FAKE_DATA) {
  fakeSettings = import("@pango/fakedata").then((_) => ({
    settings: _.newSettings(),
    hostSettings: _.newHostSettings(),
  }));
}

const SettingsAPIPath = `${SettingsName}/base`;
const HostSettingsAPIPath = `${SettingsName}/host`;

export async function fetchSettings() {
  if (import.meta.env.VITE_FAKE_DATA) {
    const { settings } = await fakeSettings;
    return settings;
  }
  const res = await fetch(`${BASE_PATH}/${SettingsAPIPath}`);
  return toJson<Settings>(res);
}

export async function saveSettings(fields: SettingsFields) {
  if (import.meta.env.VITE_FAKE_DATA) {
    const { settings, ...others } = await fakeSettings;
    const newSettings = { ...settings, ...fields };
    fakeSettings = Promise.resolve({
      ...others,
      settings: newSettings,
    });
    return newSettings;
  }

  const res = await fetch(
    `${BASE_PATH}/${SettingsAPIPath}`,
    withJSONBody(fields, { method: "PATCH" }),
  );
  return toJson<Settings>(res);
}

export async function fetchHostSettings() {
  if (import.meta.env.VITE_FAKE_DATA) {
    const { hostSettings } = await fakeSettings;
    return hostSettings;
  }
  const res = await fetch(`${BASE_PATH}/${HostSettingsAPIPath}`);
  return toJson<HostSettings>(res);
}

export async function saveHostSettings(fields: HostSettingsFields) {
  if (import.meta.env.VITE_FAKE_DATA) {
    const { hostSettings, ...others } = await fakeSettings;
    const newHostSettings = { ...hostSettings, ...fields };
    fakeSettings = Promise.resolve({
      ...others,
      hostSettings: newHostSettings,
    });
    return newHostSettings;
  }

  const res = await fetch(
    `${BASE_PATH}/${HostSettingsAPIPath}`,
    withJSONBody(fields, { method: "PATCH" }),
  );
  return toJson<HostSettings>(res);
}
