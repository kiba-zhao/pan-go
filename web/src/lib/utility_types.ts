export type DecorateType<
  Type,
  RequiredKey extends keyof Type = never,
  PartialKey extends keyof Type = never,
  ExcludeKey extends keyof Type = never,
> = Omit<Type, RequiredKey | PartialKey | ExcludeKey> &
  Required<Pick<Type, RequiredKey>> &
  Partial<Pick<Type, PartialKey>>;

export type RequiredType<
  Type,
  RequiredKey extends keyof Type = never,
  ExcludeKey extends keyof Type = never,
> = DecorateType<Type, RequiredKey, never, ExcludeKey>;

export type PartialType<
  Type,
  PartialKey extends keyof Type = never,
  ExcludeKey extends keyof Type = never,
> = DecorateType<Type, never, PartialKey, ExcludeKey>;
