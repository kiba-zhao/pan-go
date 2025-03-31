const { faker } = require("@faker-js/faker");
const path = require("node:path");
const mimeTypes = require('mime-types')

module.exports = () => {
  const folders = ["/"];
  const diskFiles = faker.helpers.multiple(
    generateDiskFile.bind(this, folders),
    {
      count: { min: 50, max: 100 },
    }
  );
  diskFiles.unshift(generateDiskRoot());

  const nodes = generateNodes();

  const nodeItems = faker.helpers.multiple(generateExtFSNodeItem, {
    count: { min: 1, max: 10 },
  });

  remoteNodes = faker.helpers
    .arrayElements(nodes, { min: 1 })
    .map(parseRemoteNode);
  const remoteItems = remoteNodes.reduce((arr, node) => {
    if (!node.available) return arr;
    return arr.concat(
      faker.helpers.multiple(
        generateExtFSRemoteNodeItem.bind(this, node.peerId),
        {
          count: { min: 1, max: 10 },
        }
      )
    );
  }, []);

  const nodeFiles = nodeItems.reduce((arr, item) => {
    if (!item.available || item.fileType !== "D") return arr;
    const folders = [""];
    return arr.concat(
      faker.helpers.multiple(
        generateExtFSNodeFile.bind(this, item.id, folders),
        {
          count: { min: 1, max: 10 },
        }
      )
    );
  }, []);

  const remoteFiles = remoteItems.reduce((arr, item) => {
    if (!item.available || item.fileType !== "D") return arr;
    const folders = [""];
    return arr.concat(
      faker.helpers.multiple(
        generateExtFSRemoteFile.bind(this, item.peerId, item.itemId, folders),
        {
          count: { min: 1, max: 10 },
        }
      )
    );
  }, []);

  const searchItems = faker.helpers.multiple(generateExtFSSearchItem, {
    count: { min: 1, max: 10 },
  });

  const searchFiles = [...nodeItems,...nodeFiles].map(generateNodeSearchFiles)
  const remoteSearchFiles = [...remoteItems,...remoteFiles].map(generateRemoteSearchFiles)
  return {
    "app-disk-files": diskFiles,
    "app-settings": generateAppSettings(
      diskFiles.filter((_) => _.fileType === "D").map((_) => _.filePath)
    ),
    "app-nodes": nodes,
    "extfs-remote-nodes": remoteNodes,
    "extfs-node-items": nodeItems,
    "extfs-remote-items": remoteItems,
    "extfs-node-files": nodeFiles,
    "extfs-remote-files": remoteFiles,
    "extfs-search-items": searchItems,
    "extfs-search-files": searchFiles,
    "extfs-remote-search-files": remoteSearchFiles,
  };
};

function generateDiskFile(folders) {
  const isDir = faker.datatype.boolean();
  const folder = faker.helpers.arrayElement(folders);
  const name = isDir ? faker.word.sample() : faker.system.fileName();
  const filePath = path.join(folder, name);

  if (isDir) {
    folders.push(filePath);
  }
  return {
    id: faker.string.nanoid(),
    name,
    filePath,
    parentPath: folder,
    fileType: isDir ? "D" : "F",
    updatedAt: faker.date.past(),
  };
}

function generateDiskRoot() {
  return {
    id: faker.string.nanoid(),
    name: "/",
    filePath: "/",
    fileType: "D",
    parentPath: "",
    updatedAt: faker.date.past(),
  };
}

function generateAppSettings(rootPaths) {
  const webAddress = faker.helpers.multiple(generateAddress, {
    count: { min: 1, max: 3 },
  });
  const peerAddress = faker.helpers.multiple(generateAddress, {
    count: { min: 1, max: 3 },
  });
  const broadcastAddress = faker.helpers.multiple(generateAddress, {
    count: { min: 1, max: 3 },
  });
  const publicAddress = faker.helpers.multiple(generateAddress, {
    count: { min: 1, max: 3 },
  });
  return {
    rootPath: faker.helpers.arrayElement(rootPaths),
    name: faker.internet.domainName(),
    webAddress,
    peerAddress,
    broadcastAddress,
    publicAddress,
    peerId: faker.helpers.arrayElement(["", faker.string.nanoid()]),
    guardEnabled: faker.datatype.boolean(),
    guardAccess: faker.datatype.boolean(),
  };
}

function generateAddress() {
  return `${faker.helpers.arrayElement([
    faker.internet.ip(),
    "0.0.0.0",
  ])}:${faker.internet.port()}`;
}

function generateNodes() {
  return faker.helpers.multiple(
    () => {
      const blocked = faker.datatype.boolean();
      return {
        id: faker.number.int(),
        peerId: faker.string.nanoid(),
        name: faker.internet.domainName(),
        blocked,
        online: blocked || faker.datatype.boolean(),
        createdAt: faker.date.past(),
        updatedAt: faker.date.past(),
      };
    },
    {
      count: { min: 1, max: 10 },
    }
  );
}

function parseRemoteNode(node) {
  return {
    id: node.id,
    peerId: node.peerId,
    name: node.name,
    available: node.online,
    createdAt: node.createdAt,
    updatedAt: node.updatedAt,
  };
}

function generateExtFSNodeItem() {
  const fileType = faker.helpers.arrayElement(["F", "D"]);
  const mimeType = fileType === "F"?generateMimeType():"";
  const extname = fileType === "F"?faker.system.fileExt(mimeType):"";
  const name = fileType === "F" ? faker.system.commonFileName(extname): faker.system.fileName({extensionCount:0});

  return {
    id: faker.number.int({ min: 1, max: 999999 }),
    name,
    filePath: faker.system.directoryPath(),
    fileType,
    mimeType,
    size: faker.number.int({ min: 1, max: 999999 }),
    enabled: faker.datatype.boolean(),
    available: faker.datatype.boolean(),
    createdAt: faker.date.past(),
    updatedAt: faker.date.past(),
  };
}

function generateExtFSRemoteNodeItem(peerId) {
  const fileType = faker.helpers.arrayElement(["F", "D"]);
  const mimeType = fileType === "F"?generateMimeType():"";
  const extname = fileType === "F"?faker.system.fileExt(mimeType):"";
  const name = fileType === "F" ? faker.system.commonFileName(extname): faker.system.fileName({extensionCount:0});
  return {
    id: faker.string.nanoid(),
    peerId,
    itemId: faker.number.int({ min: 1, max: 999999 }),
    name,
    fileType,
    mimeType,
    size: faker.number.int({ min: 1, max: 999999 }),
    available: faker.datatype.boolean(),
    createdAt: faker.date.past(),
    updatedAt: faker.date.past(),
  };
}

function generateExtFSNodeFile(itemId, folders = [""]) {
  const isDir = faker.datatype.boolean();
  const mimeType = isDir ? "" : generateMimeType();
  const folder = faker.helpers.arrayElement(folders);
  const fileType = isDir ? "D" : "F";
  const extname = isDir?"":faker.system.fileExt(mimeType);
  const name = isDir ? faker.system.fileName({extensionCount:0}):faker.system.commonFileName(extname);
  const filePath = folder.length > 0 ? path.join(folder, name):name;

  if (isDir) {
    folders.push(filePath);
  }
  return {
    id: faker.string.nanoid(),
    itemId,
    name,
    filePath,
    fileType,
    mimeType,
    parentPath: folder,
    size: faker.number.int({ min: 1, max: 999999 }),
    available: true,
    createdAt: faker.date.past(),
    updatedAt: faker.date.past(),
  };
}

function generateExtFSRemoteFile(peerId, itemId, folders = [""]) {
  const isDir = faker.datatype.boolean();
  const mimeType = isDir ? "" : generateMimeType();
  const folder = faker.helpers.arrayElement(folders);
  const fileType = isDir ? "D" : "F";
  const extname = isDir?"":faker.system.fileExt(mimeType);
  const name = isDir ? faker.system.fileName({extensionCount:0}):faker.system.commonFileName(extname);
  const filePath = folder.length > 0 ? path.join(folder, name):name;

  if (isDir) {
    folders.push(filePath);
  }
  return {
    id: faker.string.nanoid(),
    peerId,
    itemId,
    name,
    parentPath: folder,
    filePath,
    fileType,
    mimeType,
    size: faker.number.int({ min: 1, max: 999999 }),
    available: true,
    createdAt: faker.date.past(),
    updatedAt: faker.date.past(),
  };
}

function generateExtFSSearchItem() {
  return {
    id: faker.number.int({ min: 1, max: 999999 }),
    query: faker.word.words(),
    createdAt: faker.date.past(),
    updatedAt: faker.date.past(),
  };
}

function generateNodeSearchFiles(nodefile){
  return {
    id: faker.number.int({ min: 1, max: 999999 }),
    itemId: nodefile.itemId === void 0?nodefile.id:nodefile.itemId,
    name: nodefile.name,
    fileType: nodefile.fileType,
    size: nodefile.size,
    available: nodefile.available,
    filePath: nodefile.itemId === void 0?"":nodefile.filePath,
    mimeType: nodefile.mimeType,
    available: nodefile.available,
    score: faker.number.int({ min: 1, max: 10 }),
    tokens: [nodefile.name],
    createdAt: faker.date.past(),
    updatedAt: faker.date.past(),
  }
}

function generateRemoteSearchFiles(remotefile){
  return {
    id: faker.number.int({ min: 1, max: 999999 }),
    peerId: remotefile.peerId,
    itemId: remotefile.itemId,
    name: remotefile.name,
    fileType: remotefile.fileType,
    size: remotefile.size,
    available: remotefile.available,
    filePath: remotefile.filePath,
    mimeType: remotefile.mimeType,
    available: remotefile.available,
    score: faker.number.int({ min: 1, max: 10 }),
    tokens: [remotefile.name],
    createdAt: faker.date.past(),
    updatedAt: faker.date.past(),
  }
}

function generateMimeType() {
  return faker.helpers.arrayElement([
    "text/plain",
    "text/html",
    "application/pdf",
    "application/json",
    "application/xml",
    "application/octet-stream",
  ])
}
