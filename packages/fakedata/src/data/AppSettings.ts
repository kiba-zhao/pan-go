import { faker } from "@faker-js/faker";
import type { AppSettings } from "@pango/data";

export function newAppSettings(): AppSettings {
  const broadcastAddress = faker.helpers.multiple(newAddress, {
    count: { min: 1, max: 3 },
  });
  const publicAddress = faker.helpers.multiple(newAddress, {
    count: { min: 1, max: 3 },
  });
  return {
    rootPath: faker.system.directoryPath(),
    name: faker.internet.domainName(),
    peerId: faker.helpers.arrayElement(["", faker.string.nanoid()]),
    webPort: faker.number.int({ min: 1, max: 65535 }),
    peerPort: faker.number.int({ min: 1, max: 65535 }),
    broadcastAddress,
    publicAddress,
    guardEnabled: faker.datatype.boolean(),
    guardAccess: faker.datatype.boolean(),
  } as AppSettings;
}

function newAddress() {
  return `${faker.helpers.arrayElement([
    faker.internet.ip(),
    "0.0.0.0",
  ])}:${faker.internet.port()}`;
}
