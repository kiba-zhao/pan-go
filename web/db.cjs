const { faker } = require("@faker-js/faker");
const path = require("node:path");

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
    const folders = ["/"];
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
    const folders = ["/"];
    return arr.concat(
      faker.helpers.multiple(
        generateExtFSRemoteFile.bind(this, item.peerId, item.itemId, folders),
        {
          count: { min: 1, max: 10 },
        }
      )
    );
  }, []);

  const searchItems = faker.helpers.multiple(generateSearchItem, {
    count: { min: 1, max: 10 },
  });

  const searchFiles = searchItems.reduce((files, item) => {
    const files_ = generateSearchFiles(
      { nodeFiles, remoteFiles, remoteItems, nodeItems },
      item
    );
    return files.concat(files_);
  }, []);

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
    tagQuantity: 0,
    pendingTagQuantity: 0,
  };
}

function generateExtFSNodeItem() {
  return {
    id: faker.number.int({ min: 1, max: 999999 }),
    name: faker.word.sample(),
    filePath: faker.system.directoryPath(),
    fileType: faker.helpers.arrayElement(["F", "D"]),
    size: faker.number.int({ min: 1, max: 999999 }),
    enabled: faker.datatype.boolean(),
    available: faker.datatype.boolean(),
    createdAt: faker.date.past(),
    updatedAt: faker.date.past(),
    tagQuantity: 0,
    pendingTagQuantity: 0,
  };
}

function generateExtFSRemoteNodeItem(peerId) {
  return {
    id: faker.string.nanoid(),
    peerId,
    itemId: faker.number.int({ min: 1, max: 999999 }),
    name: faker.word.sample(),
    fileType: faker.helpers.arrayElement(["F", "D"]),
    size: faker.number.int({ min: 1, max: 999999 }),
    available: faker.datatype.boolean(),
    createdAt: faker.date.past(),
    updatedAt: faker.date.past(),
    tagQuantity: 0,
    pendingTagQuantity: 0,
  };
}

function generateExtFSNodeFile(itemId, folders = ["/"]) {
  const isDir = faker.datatype.boolean();
  const folder = faker.helpers.arrayElement(folders);
  const fileType = isDir ? "D" : "F";
  const name = isDir ? faker.word.sample() : faker.system.fileName();
  const filePath = path.join(folder, name);
  if (isDir) {
    folders.push(filePath);
  }
  return {
    id: faker.string.nanoid(),
    itemId,
    name,
    filePath,
    fileType,
    parentPath: folder,
    size: faker.number.int({ min: 1, max: 999999 }),
    available: true,
    createdAt: faker.date.past(),
    updatedAt: faker.date.past(),
    tagQuantity: 0,
    pendingTagQuantity: 0,
  };
}

function generateExtFSRemoteFile(peerId, itemId, folders = ["/"]) {
  const isDir = faker.datatype.boolean();
  const folder = faker.helpers.arrayElement(folders);
  const fileType = isDir ? "D" : "F";
  const name = isDir ? faker.word.sample() : faker.system.fileName();
  const filePath = path.join(folder, name);
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
    size: faker.number.int({ min: 1, max: 999999 }),
    available: true,
    createdAt: faker.date.past(),
    updatedAt: faker.date.past(),
    tagQuantity: 0,
    pendingTagQuantity: 0,
  };
}

function generateSearchItem() {
  return {
    id: faker.string.nanoid(),
    query: faker.word.words(),
    createdAt: faker.date.past(),
    updatedAt: faker.date.past(),
  };
}

function generateSearchFiles(ctx, item) {
  const { nodeFiles, remoteFiles, remoteItems, nodeItems } = ctx;

  return [].concat(
    faker.helpers
      .arrayElements(nodeFiles)
      .map((_) => generateSearchFileWithNodeFile(_, item)),
    faker.helpers
      .arrayElements(remoteFiles)
      .map((_) => generateSearchFileWithRemoteFile(_, item)),
    faker.helpers
      .arrayElements(nodeItems)
      .map((_) => generateSearchFileWithNodeItem(_, item)),
    faker.helpers
      .arrayElements(remoteItems)
      .map((_) => generateSearchFileWithRemoteItem(_, item))
  );
}

function generateSearchFileWithNodeFile(nodeFile, item) {
  return {
    id: faker.string.nanoid(),
    name: nodeFile.name,
    fileType: nodeFile.fileType,
    size: nodeFile.size,
    available: nodeFile.available,
    reason: item.q,
    referId: nodeFile.id,
    referType: "NF",
    createdAt: faker.date.past(),
    updatedAt: faker.date.past(),
  };
}

function generateSearchFileWithRemoteFile(remoteFile, item) {
  return {
    id: faker.string.nanoid(),
    name: remoteFile.name,
    fileType: remoteFile.fileType,
    size: remoteFile.size,
    available: remoteFile.available,
    reason: item.q,
    referId: remoteFile.id,
    referType: "RF",
    createdAt: faker.date.past(),
    updatedAt: faker.date.past(),
  };
}

function generateSearchFileWithNodeItem(nodeItem, item) {
  return {
    id: faker.string.nanoid(),
    name: nodeItem.name,
    fileType: nodeItem.fileType,
    size: nodeItem.size,
    available: nodeItem.available,
    reason: item.q,
    referId: nodeItem.id.toString(),
    referType: "N",
    createdAt: faker.date.past(),
    updatedAt: faker.date.past(),
  };
}

function generateSearchFileWithRemoteItem(remoteItem, item) {
  return {
    id: faker.string.nanoid(),
    name: remoteItem.name,
    fileType: remoteItem.fileType,
    size: remoteItem.size,
    available: remoteItem.available,
    reason: item.q,
    referId: remoteItem.id.toString(),
    referType: "R",
    createdAt: faker.date.past(),
    updatedAt: faker.date.past(),
  };
}
