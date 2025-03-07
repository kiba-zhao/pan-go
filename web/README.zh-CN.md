## Web 包说明

当前包内源码用于生成 web 控制台页面。支持 PC,mobile 等不同设备访问。

**web 控制台功能:**

- 应用设置: 节点信息,服务,广播等参数设置
- 远程节点：其他设备节点信息设置
- ExtFS: 所有节点及节点文件管理器，包含浏览及检索等功能。

## React + TypeScript + Vite

源码为 [Typescript](https://www.typescriptlang.org/) 和 [React](https://react.dev/) 前端框架，使用 [vite](https://vite.dev/) 进行开发和源码打包。

## Material UI

UI 使用 [Material UI](https://mui.com/material-ui/getting-started/) 以及 [React-Admin](https://marmelab.com/react-admin/)。

## Mock

前端开发采用 [json-server](https://github.com/typicode/json-server/tree/v0)模拟服务端数据。

## 主要包内文件简介

```
|--public   # web静态资源目录
|--src      # 页面源码目录
|   |--api          # 后端api定义文件
|   |--assets       # vite打包资源文件
|   |--components   # react组件文件
|   |--locales      # 国际化多语言文件
|   |--api.ts           # 后端api定义
|   |--App.css
|   |--App.tsx          # 应用组件
|   |--i18n.ts
|   |--index.css
|   |--main.tsx         # 打包入口文件
|   |--vite-env.d.ts    # 运行环境变量声明
|--.env     # 默认环境变量
|--.env.development # 开发用环境变量（npm run dev）
|--db.cjs   # json-server构造模拟数据脚本
|--middleware.cjs   # json-server中间件
|--package.json     # node包定义文件，包含包信息，模块依赖等
|--routes.json      # json-server路由映射
|--vite.config.ts   # vite配置文件
```

## 开发

```shell
# 依赖安装
npm install

# 启动json-server
npm run mock

# 启动vite开发
npm run dev

```

## 构建

```shell
# 依赖安装
npm install

# 启动vite打包源码到dist目录
npm run build

```
