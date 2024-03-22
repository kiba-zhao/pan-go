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

  return {
    targets,
    "disk-files": diskFiles,
  };
};

function generateTarget() {
  return {
    id: faker.number.int({ min: 1, max: 999999 }),
    name: faker.word.sample(),
    filepath: faker.system.directoryPath(),
    enabled: faker.datatype.boolean(),
    version: faker.number.int(0, 255),
    invalid: faker.datatype.boolean(),
    createAt: faker.date.past(),
    updateAt: faker.date.past(),
  };
}

function generateDiskFile(folders) {
  const isDir = faker.datatype.boolean();
  const folder = faker.helpers.arrayElement(folders);
  const filepath = path.join(
    folder,
    isDir ? faker.word.sample() : faker.system.fileName()
  );

  if (isDir) {
    folders.push(filepath);
  }
  return {
    id: faker.string.nanoid(),
    filepath,
    parent: folder,
    fileType: isDir ? "D" : "F",
    mimeType: isDir
      ? undefined
      : faker.helpers.arrayElement([faker.system.mimeType(), undefined]),
  };
}

function generateDiskRoot() {
  return {
    id: faker.string.nanoid(),
    filepath: "/",
    fileType: "D",
    parent: "",
  };
}