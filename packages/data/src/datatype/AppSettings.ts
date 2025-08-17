export type AppSettings = {
  name: string;
  rootPath: string;
  peerId: string;
  peerPort: number;
  broadcastAddrs: string[];
  publicAddrs: string[];
  enabled: boolean;
  webAddr: string;
};

export type AppSettingsFields = Omit<AppSettings, "peerId" | "rootPath">;
