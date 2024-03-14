const { faker } = require("@faker-js/faker");

module.exports = () => {
  const targets = faker.helpers.multiple(generateTarget, {
    count: { min: 1, max: 10 },
  });
  return {
    targets,
  };
};

function generateTarget() {
  return {
    id: faker.number.int({ min: 1, max: 999999 }),
    name: faker.word.sample(),
    filepath: faker.system.directoryPath(),
    enabled: faker.datatype.boolean(),
    version: faker.number.int(0, 255),
    createAt: faker.date.past(),
    updateAt: faker.date.past(),
  };
}
