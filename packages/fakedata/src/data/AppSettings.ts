// import { faker } from "@faker-js/faker";
// import type { AppSettings } from "@pango/data";

// export function newAppSettings(): AppSettings {
//   const broadcastAddrs = faker.helpers.multiple(newAddress, {
//     count: { min: 1, max: 3 },
//   });
//   const publicAddrs = faker.helpers.multiple(newAddress, {
//     count: { min: 1, max: 3 },
//   });
//   const peerPort = faker.number.int({ min: 1, max: 65535 });
//   return {
//     name: faker.internet.domainName(),
//     rootPath: faker.system.directoryPath(),
//     peerId: faker.helpers.arrayElement(["", faker.string.nanoid()]),
//     peerPort,
//     broadcastAddrs,
//     publicAddrs,
//     enabled: faker.datatype.boolean(),
//     webAddr: faker.helpers.arrayElement([
//       "127.0.0.1:" + peerPort.toString(),
//       "0.0.0.0:" + peerPort.toString(),
//       "localhost:" + peerPort.toString(),
//       "[::1]:" + peerPort.toString(),
//       "[::]:" + peerPort.toString(),
//     ]),
//   } as AppSettings;
// }

// function newAddress() {
//   return `${faker.helpers.arrayElement([
//     faker.internet.ip(),
//     "0.0.0.0",
//   ])}:${faker.internet.port()}`;
// }
