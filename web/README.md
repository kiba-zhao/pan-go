## Web

this package is used to generate web console pages,Support PC and mobile.

**What can web console do:**

- Application Settings: set information, service, broadcast and others
- Remote Settings：set remote node.
- ExtFS: A file manager for all nodes.You can browse and search all files.

## React + TypeScript + Vite

Using [Typescript](https://www.typescriptlang.org/) and [React](https://react.dev/).

Using [vite](https://vite.dev/) for for packaging and development.

## Material UI

Most pages use [Material UI](https://mui.com/material-ui/getting-started/) and [React-Admin](https://marmelab.com/react-admin/) components.

## Mock

Using [json-server](https://github.com/typicode/json-server/tree/v0) to Simulate Development Data

## Main Package Files

```
|--public   # Web Static Resources
|--src      # Source files
|   |--api          # api definition Files
|   |--assets       # vite resources
|   |--components   # UI components
|   |--locales      # i18n language files
|   |--api.ts           # api definition
|   |--App.css
|   |--App.tsx          # App component
|   |--i18n.ts
|   |--index.css
|   |--main.tsx         # main file
|   |--vite-env.d.ts    # environment variable declaration file
|--.env     # default environment variable
|--.env.development # environment variable for development（npm run dev）
|--db.cjs   # resource definition file for json-server
|--middleware.cjs   # middleware for json-server
|--package.json     # package file for nodejs
|--routes.json      # route definition file for json-server
|--vite.config.ts   # vite config file
```

## Development

```shell
# Dependency installation
npm install

# Start JSON server to mock data
npm run mock

# Start Vite
npm run dev

```

## Build

```shell
# Dependency installation
npm install

# pack the source code to dist folder
npm run build

```
