import { faker } from "@faker-js/faker";

export type SeedOptions = Parameters<typeof faker.helpers.multiple>[1];

type PastDateOptions = Parameters<typeof faker.date.past>[0];
export function newPastDateString(opts?: PastDateOptions): string {
  const date = faker.date.past(opts);
  return toDateString(date);
}

export function toDateString(date: Date): string {
  return `${date.toLocaleDateString()} ${date.toLocaleTimeString()}`;
}
