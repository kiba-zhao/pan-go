import { faker } from "@faker-js/faker";

export function newMimeType() {
  return faker.helpers.arrayElement([
    "text/plain",
    "text/html",
    "application/pdf",
    "application/json",
    "application/xml",
    "application/octet-stream",
  ]);
}
