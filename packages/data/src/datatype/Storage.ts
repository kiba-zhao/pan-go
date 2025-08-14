enum FileType {
  File = "F",
  Directory = "D",
}

export type Storage = {
  id: number;
  name: string;
  filePath: string;
  fileType: FileType;
  mimeType: string;
  size: number;
  enabled: boolean;
  available: boolean;
  createdAt: string;
  updatedAt: string;
};

export type StorageFields = Omit<
  Storage,
  | "id"
  | "createdAt"
  | "updatedAt"
  | "available"
  | "fileType"
  | "mimeType"
  | "size"
>;
