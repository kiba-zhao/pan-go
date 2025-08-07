export type ActionFields<Action extends unknown> = {
  _action?: Action;
};

export type RangeFields = {
  _start?: number;
  _end?: number;
};

export type SortFields<Fields extends string> = {
  _sort?: NoInfer<Fields>;
  _order?: "desc" | "asc";
};

export type QFields<Q extends unknown> = {
  q?: Q;
};

export type SearchResults<Entity extends unknown> = [number, Entity[]];

export type CursorFields<Cursor extends unknown> = {
  _cursor?: Cursor;
  _limit?: number;
};

export type FetchResultsMeta<Tag extends unknown, Cursor extends unknown> = {
  offset: number;
  tag: Tag;
  prev?: Cursor;
  next?: Cursor;
};

export type FetchResults<
  Entity extends unknown,
  Tag extends unknown,
  Cursor extends unknown
> = [FetchResultsMeta<Tag, Cursor>, Entity[]];
