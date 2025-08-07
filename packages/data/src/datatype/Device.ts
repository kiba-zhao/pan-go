export type Device = {
  id: number;
  name: string;
  peerId: string;
  enabled: boolean;
  online: boolean;
  networkAddrs: string[] | null;
  createdAt: string;
  updatedAt: string;
};

export type DeviceFields = Omit<
  Device,
  "id" | "createdAt" | "updatedAt" | "online" | "networkAddrs"
>;
