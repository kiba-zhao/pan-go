import { faker } from "@faker-js/faker";
import type { Device } from "@pango/data";
import { newPastDateString, type SeedOptions } from "./Common";

export function seedDevice(opts?: SeedOptions) {
  return faker.helpers.multiple(newDevice, opts);
}

function newDevice() {
  const networkAddrs = faker.helpers.arrayElements([
    faker.internet.ipv4(),
    `${faker.internet.ipv4()}:${faker.internet.port()}`,
    `[${faker.internet.ipv6()}]`,
    `[${faker.internet.ipv6()}]:${faker.internet.port()}`,
    faker.internet.domainName(),
    `${faker.internet.domainName()}:${faker.internet.port()}`,
  ]);
  const enabled = faker.datatype.boolean();
  return {
    id: faker.number.int(),
    peerId: faker.string.nanoid(),
    name: faker.internet.domainName(),
    enabled,
    online: enabled && faker.datatype.boolean(),
    networkAddrs,
    createdAt: newPastDateString(),
    updatedAt: newPastDateString(),
  } as Device;
}
