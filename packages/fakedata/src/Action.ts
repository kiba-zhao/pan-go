import { faker } from "@faker-js/faker";
import type {
  CursorFields,
  FetchResults,
  FetchResultsMeta,
  RangeFields,
  SearchResults,
  SortFields,
} from "@pango/data";
import { toDateString } from "./data/Common";

type EntityBase<ID = string | number> = {
  id: ID;
  updatedAt?: string;
  createdAt?: string;
};
type NextIDFunction<ID = string | number> = (entities: EntityBase<ID>[]) => ID;

export function nextID(entities: EntityBase<number>[]): number {
  return entities.length > 0
    ? Math.max(...entities.map((entity) => entity.id)) + 1
    : parseInt((Math.random() * 10000).toFixed());
}

export function newNanoID(entities: EntityBase<string>[]): string {
  let id = faker.string.nanoid();
  while (entities.some((entity) => entity.id === id)) {
    id = faker.string.nanoid();
  }

  return id;
}

export function create<
  ID extends string | number,
  Entity extends EntityBase<ID>
>(
  fields: Omit<Entity, "id" | "updatedAt" | "createdAt">,
  entities: Entity[],
  nextId: NextIDFunction<ID>
): Entity {
  const id = nextId(entities) as ID;
  const entity = {
    ...fields,
    id,
    createdAt: toDateString(new Date()),
    updatedAt: toDateString(new Date()),
  } as Entity;
  entities.push(entity as Entity);
  return entity;
}

export function replace<
  ID extends string | number,
  Entity extends EntityBase<ID>
>(
  id: ID,
  fields: Omit<Entity, "id" | "updatedAt" | "createdAt">,
  entities: Entity[]
): Entity | null {
  return update(id, fields, entities);
}

export function update<
  ID extends string | number,
  Entity extends EntityBase<ID>
>(
  id: ID,
  fields: Partial<Omit<Entity, "id" | "updatedAt" | "createdAt">>,
  entities: Entity[]
): Entity | null {
  const index = entities.findIndex((e) => e.id === id);
  if (index < 0) return null;
  entities[index] = {
    ...entities[index],
    ...fields,
    updatedAt: toDateString(new Date()),
  };
  return entities[index];
}

export function destroy<
  ID extends string | number,
  Entity extends EntityBase<ID>
>(id: ID, entities: Entity[]): Entity | null {
  const index = entities.findIndex((e) => e.id === id);
  if (index < 0) return null;
  const entitiesWithDeleted = entities.splice(index, 1);
  return entitiesWithDeleted[0];
}

export function findById<
  ID extends string | number,
  Entity extends EntityBase<ID>
>(id: ID, entities: Entity[]): Entity | null {
  return entities.find((e) => e.id === id) || null;
}

export function findOne<
  ID extends string | number,
  Entity extends EntityBase<ID>
>(
  predicate: (value: Entity, index: number, obj: Entity[]) => unknown,
  entities: Entity[]
): Entity | null {
  return entities.find(predicate) || null;
}

export function sort<ID extends string | number, Entity extends EntityBase<ID>>(
  sortField: SortFields<Extract<keyof Entity, "string">>,
  entities: Entity[]
): Entity[] {
  if (entities.length <= 1) return entities;

  const _sort = sortField._sort || "id";
  const _order = sortField._order || "asc";

  return entities.slice().sort((a, b) => {
    const fieldA = a[_sort];
    const fieldB = b[_sort];
    if (fieldA === fieldB) {
      return 0;
    }
    if (typeof fieldA === void 0) {
      return _order === "asc" ? 1 : -1;
    }
    if (typeof fieldB === void 0) {
      return _order === "asc" ? -1 : 1;
    }

    if (typeof fieldA === "string" && typeof fieldB === "string") {
      return (
        (fieldA as string).localeCompare(fieldB as string) *
        (_order === "asc" ? 1 : -1)
      );
    }

    if (typeof fieldA === "number" && typeof fieldB === "number") {
      return (
        (fieldA as number) - (fieldB as number) * (_order === "asc" ? 1 : -1)
      );
    }

    return fieldA && _order === "asc" ? 1 : -1;
  });
}

export function range<T extends unknown>(
  rangeField: RangeFields,
  entities: T[]
): T[] {
  return entities.slice(rangeField._start, rangeField._end);
}

export function search<
  ID extends string | number,
  Entity extends EntityBase<ID>
>(
  predicate: (entity: Entity) => boolean,
  opts: SortFields<Extract<keyof Entity, "string">> & RangeFields,
  entities: Entity[]
): SearchResults<Entity> {
  let entities_ = entities.filter(predicate);
  const total = entities_.length;

  if (total <= 0) return [total, []];
  entities_ = sort(opts, entities_);
  const end = opts._end ?? -1;
  if (end > 0 && total <= end) return [total, entities_];
  entities_ = range(opts, entities_);
  return [total, entities_];
}

export function fetch<
  ID extends string | number,
  Entity extends EntityBase<ID>,
  CursorFieldKey extends NoInfer<Extract<keyof Entity, "string">>
>(
  condition: CursorFields<string>,
  entities: Entity[],
  field: CursorFieldKey,
  defaultTag: Entity[CursorFieldKey]
): FetchResults<Entity, Entity[CursorFieldKey], Entity[CursorFieldKey]> {
  if (entities.length <= 0) {
    return [{ tag: defaultTag, offset: 0 }, []];
  }

  const tag = (entities.at(0) as Entity)[field];
  const offset = entities.length;
  const { _cursor, _limit } = condition;
  if (_limit === 0) {
    return [{ tag, offset }, []];
  }

  let start_ = -1;
  let end_ = entities.length;

  if (_cursor !== void 0) {
    start_ = entities.findIndex((entity) => entity[field] === _cursor);
    if (start_ < 0) {
      return [{ tag, offset: start_ }, []];
    }
  }

  if (_limit !== void 0) {
    if (_limit > 0) {
      if (start_ >= end_ - 1) return [{ tag, offset: end_ }, []];
      start_++;
      end_ = start_ + _limit;
      if (end_ >= entities.length) end_ = entities.length;
    } else {
      if (start_ === 0) return [{ tag, offset: start_ }, []];
      if (start_ > 0) end_ = start_;
      start_ = end_ + _limit;
      if (start_ < 0) start_ = 0;
    }
  }
  const entities_ = entities.slice(start_, end_);
  const meta = { tag, offset: start_ } as FetchResultsMeta<
    Entity[CursorFieldKey],
    Entity[CursorFieldKey]
  >;
  meta.prev = entities_.at(0)?.[field];
  meta.next = entities_.at(-1)?.[field];

  if (start_ === 0) {
    meta.prev = void 0;
  }
  if (end_ === entities.length) {
    meta.next = void 0;
  }

  return [meta, entities_];
}
