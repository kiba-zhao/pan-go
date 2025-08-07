import type { Device } from "./datatype/Device";

const DeviceQRCodeURI = "pango://app/devices";

export type DeviceQRCodeValue = Pick<Device, "name" | "peerId">;
export function MarshalDeviceQRCode({
  name,
  peerId,
}: DeviceQRCodeValue): string {
  const query = new URLSearchParams({ peerId, name });
  return `${DeviceQRCodeURI}?${query.toString()}`;
}

export function UnmarshalDeviceQRCode(
  value: string
): DeviceQRCodeValue | undefined {
  if (
    !value.startsWith(DeviceQRCodeURI) ||
    value.length <= DeviceQRCodeURI.length + 1
  )
    return;
  const query = new URLSearchParams(value.slice(DeviceQRCodeURI.length + 1));
  const peerId = query.get("peerId");
  const name = query.get("name");
  if (!peerId || peerId.length <= 0 || !name || name.length <= 0) return;
  return { name, peerId };
}

console.log("");
