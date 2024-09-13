const { faker } = require("@faker-js/faker");
const path = require("node:path");

module.exports = () => {
  const targets = faker.helpers.multiple(generateTarget, {
    count: { min: 1, max: 10 },
  });

  const folders = ["/"];
  const diskFiles = faker.helpers.multiple(
    generateDiskFile.bind(this, folders),
    {
      count: { min: 50, max: 100 },
    }
  );
  diskFiles.unshift(generateDiskRoot());

  const targetFiles = targets.reduce((files, target) => {
    return files.concat(
      faker.helpers.multiple(() => generateTargetFile(target), {
        count: { min: 30, max: 100 },
      })
    );
  }, []);

  const nodeItems = faker.helpers.multiple(generateExtFSNodeItem, {
    count: { min: 1, max: 10 },
  });

  return {
    "extfs-targets": targets,
    "extfs-target-files": targetFiles,
    "app-disk-files": diskFiles,
    "app-settings": generateAppSettings(
      diskFiles.filter((_) => _.fileType === "D").map((_) => _.filepath)
    ),
    "app-nodes": generateNodes(),
    "extfs-items": generateExtFSItems(),
    "extfs-node-items": nodeItems,
  };
};

function generateTarget() {
  return {
    id: faker.number.int({ min: 1, max: 999999 }),
    name: faker.word.sample(),
    filepath: faker.system.directoryPath(),
    enabled: faker.datatype.boolean(),
    version: faker.number.int(0, 255),
    available: faker.datatype.boolean(),
    createdAt: faker.date.past(),
    updatedAt: faker.date.past(),
  };
}

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

function generateTargetFile(target) {
  return {
    targetId: target.id,
    available: target.available,
    id: faker.number.int({ min: 1, max: 999999 }),
    filepath: faker.system.filePath(),
    size: faker.number.int(),
    modTime: faker.date.past(),
    checkSum: faker.string.alphanumeric(88),
    mimeType: faker.system.mimeType(),
    createdAt: faker.date.past(),
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

function generateExtFSItems() {
  let fileItems;
  let parentItems = [];
  let depth = 10;

  for (;;) {
    depth--;
    if (depth <= 0) break;
    if (parentItems.length <= 0) {
      parentItems.push(generateExtFSItem("N"));
      parentItems.push(
        ...faker.helpers.multiple(generateExtFSItem.bind(this, "RN"), {
          count: { min: 1, max: 3 },
        })
      );
      parentItems.push(
        ...faker.helpers.multiple(generateExtFSItem.bind(this, "S"), {
          count: { min: 2, max: 6 },
        })
      );

      fileItems = Array.from(parentItems);
      continue;
    }

    const [_, parentItems_] = parentItems.reduce(
      (ctx, parentItem) => {
        const fileItems_ = faker.helpers.multiple(
          generateExtFSItemByParent.bind(this, parentItem),
          {
            count: { min: 2, max: 5 },
          }
        );
        if (fileItems_.length > 0) {
          const [ctxFileItems, ctxParentItems] = ctx;
          ctxFileItems.push(...fileItems_);
          ctxParentItems.push(
            ...fileItems_.filter((_) =>
              ["S", "N", "RN", "D", "RD"].includes(_.fileType)
            )
          );
        }
        return ctx;
      },
      [fileItems, []]
    );
    if (parentItems_.length <= 0) break;

    parentItems = parentItems_;
  }

  return fileItems;
}

function generateExtFSItemByParent(parentItem) {
  const fileTypes = generateExtFSFileTypes(
    parentItem ? parentItem.fileType : void 0
  );
  if (!fileTypes.length) return;
  const fileType = faker.helpers.arrayElement(fileTypes);

  const item = generateExtFSItem(fileType);
  if (parentItem) {
    item.parentId = parentItem.id;
  }
  return item;
}

function generateExtFSItem(fileType) {
  return {
    id: faker.string.nanoid(),
    fileType,
    name: faker.system.fileName(),
    updatedAt: faker.date.past().toLocaleString(),
    tagQuantity: fileType === "N" ? 0 : faker.number.int(5),
    pendingTagQuantity: fileType === "N" ? 0 : faker.number.int(5),
    disabled: fileType === "N" || faker.datatype.boolean(),
    linkId: faker.number.int().toString(),
  };
}

function generateExtFSFileTypes(parentFileType) {
  if (parentFileType === "S") return ["S", "D", "F", "RD", "RF"];
  if (parentFileType === "D" || parentFileType === "N") return ["F", "D"];
  if (parentFileType === "RD" || parentFileType === "RN") return ["RF", "RD"];
  return [];
}

function generateExtFSNodeItem() {
  return {
    id: faker.number.int({ min: 1, max: 999999 }),
    name: faker.word.sample(),
    filepath: faker.system.directoryPath(),
    enabled: faker.datatype.boolean(),
    available: faker.datatype.boolean(),
    createdAt: faker.date.past(),
    updatedAt: faker.date.past(),
  };
}
