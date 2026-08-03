export type ExtraState<Type, State extends any> = {
  type?: Type;
} & State;
