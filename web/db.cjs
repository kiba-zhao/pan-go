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
  const remoteNodeItems = remoteNodes.reduce((arr, node) => {
    if (!node.available) return arr;
    return arr.concat(
      faker.helpers.multiple(
        generateExtFSRemoteNodeItem.bind(this, node.nodeId),
        {
          count: { min: 1, max: 10 },
        }
      )
    );
  }, []);

  const fileItems = nodeItems.reduce((arr, item) => {
    if (!item.available || item.filetype !== "D") return arr;
    const folders = ["/"];
    return arr.concat(
      faker.helpers.multiple(
        generateExtFSFileItem.bind(this, item.id, folders),
        {
          count: { min: 1, max: 10 },
        }
      )
    );
  }, []);

  const remoteFileItems = remoteNodeItems.reduce((arr, item) => {
    if (!item.available || item.filetype !== "D") return arr;
    const folders = ["/"];
    return arr.concat(
      faker.helpers.multiple(
        generateExtFSRemoteFileItem.bind(
          this,
          item.nodeId,
          item.itemId,
          folders
        ),
        {
          count: { min: 1, max: 10 },
        }
      )
    );
  }, []);

  return {
    "app-disk-files": diskFiles,
    "app-settings": generateAppSettings(
      diskFiles.filter((_) => _.fileType === "D").map((_) => _.filepath)
    ),
    "app-nodes": nodes,
    "extfs-remote-nodes": remoteNodes,
    "extfs-node-items": nodeItems,
    "extfs-remote-node-items": remoteNodeItems,
    "extfs-file-items": fileItems,
    "extfs-remote-file-items": remoteFileItems,
  };
};

function generateDiskFile(folders) {
  const isDir = faker.datatype.boolean();
  const folder = faker.helpers.arrayElement(folders);
  const name = isDir ? faker.word.sample() : faker.system.fileName();
  const filepath = path.join(folder, name);

  if (isDir) {
    folders.push(filepath);
  }
  return {
    id: faker.string.nanoid(),
    name,
    filepath,
    parent: folder,
    fileType: isDir ? "D" : "F",
    updatedAt: faker.date.past(),
  };
}

function generateDiskRoot() {
  return {
    id: faker.string.nanoid(),
    name: "/",
    filepath: "/",
    fileType: "D",
    parent: "",
    updatedAt: faker.date.past(),
  };
}

function generateAppSettings(rootPaths) {
  const webAddress = faker.helpers.multiple(generateAddress, {
    count: { min: 1, max: 3 },
  });
  const nodeAddress = faker.helpers.multiple(generateAddress, {
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
    nodeAddress,
    broadcastAddress,
    publicAddress,
    nodeId: faker.helpers.arrayElement(["", faker.string.nanoid()]),
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
        nodeId: faker.string.nanoid(),
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
    nodeId: node.nodeId,
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
    filepath: faker.system.directoryPath(),
    filetype: faker.helpers.arrayElement(["F", "D"]),
    size: faker.number.int({ min: 1, max: 999999 }),
    enabled: faker.datatype.boolean(),
    available: faker.datatype.boolean(),
    createdAt: faker.date.past(),
    updatedAt: faker.date.past(),
    tagQuantity: 0,
    pendingTagQuantity: 0,
  };
}

function generateExtFSRemoteNodeItem(nodeId) {
  return {
    id: faker.string.nanoid(),
    nodeId,
    itemId: faker.number.int({ min: 1, max: 999999 }),
    name: faker.word.sample(),
    filetype: faker.helpers.arrayElement(["F", "D"]),
    size: faker.number.int({ min: 1, max: 999999 }),
    available: faker.datatype.boolean(),
    createdAt: faker.date.past(),
    updatedAt: faker.date.past(),
    tagQuantity: 0,
    pendingTagQuantity: 0,
  };
}

function generateExtFSFileItem(itemId, folders = ["/"]) {
  const isDir = faker.datatype.boolean();
  const folder = faker.helpers.arrayElement(folders);
  const filetype = isDir ? "D" : "F";
  const name = isDir ? faker.word.sample() : faker.system.fileName();
  const filepath = path.join(folder, name);
  if (isDir) {
    folders.push(filepath);
  }
  return {
    id: faker.string.nanoid(),
    itemId,
    name,
    filepath,
    filetype,
    parentPath: folder,
    size: faker.number.int({ min: 1, max: 999999 }),
    available: true,
    createdAt: faker.date.past(),
    updatedAt: faker.date.past(),
    tagQuantity: 0,
    pendingTagQuantity: 0,
  };
}

function generateExtFSRemoteFileItem(nodeId, itemId, folders = ["/"]) {
  const isDir = faker.datatype.boolean();
  const folder = faker.helpers.arrayElement(folders);
  const filetype = isDir ? "D" : "F";
  const name = isDir ? faker.word.sample() : faker.system.fileName();
  const filepath = path.join(folder, name);
  if (isDir) {
    folders.push(filepath);
  }
  return {
    id: faker.string.nanoid(),
    nodeId,
    itemId,
    name,
    parentPath: folder,
    filepath,
    filetype,
    size: faker.number.int({ min: 1, max: 999999 }),
    available: true,
    createdAt: faker.date.past(),
    updatedAt: faker.date.past(),
    tagQuantity: 0,
    pendingTagQuantity: 0,
  };
}
