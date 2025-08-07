export type StorageSearchItem = {
  id: number;
  query: string;
  createdAt: string;
  updatedAt: string;
};

export type StorageSearchItemFields = Omit<
  StorageSearchItem,
  "id" | "createdAt" | "updatedAt"
>;
