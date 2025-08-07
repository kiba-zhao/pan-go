export type AppSettings = {
  name: string;
  rootPath: string;
  peerId: string;
  webPort: number;
  peerPort: number;
  broadcastAddress: string[];
  publicAddress: string[];
  guardEnabled: boolean;
  guardAccess: boolean;
};

export type AppSettingsFields = Omit<AppSettings, "peerId" | "rootPath">;
