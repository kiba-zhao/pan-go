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
  return {
    targets,
    "disk-files": diskFiles,
    "target-files": targetFiles,
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
    createAt: faker.date.past(),
    updateAt: faker.date.past(),
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
    updateAt: faker.date.past(),
  };
}

function generateDiskRoot() {
  return {
    id: faker.string.nanoid(),
    name: "/",
    filepath: "/",
    fileType: "D",
    parent: "",
    updateAt: faker.date.past(),
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
    createAt: faker.date.past(),
    updateAt: faker.date.past(),
  };
}
