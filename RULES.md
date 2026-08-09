# 开发规则
 项目代码包含web应用以及mobile应用两部分.

 - web技术栈: 采用前后端分离架构.,
   - 前端: react 18, typescript
   - 后端: go 1.26.5, gin
   - 数据库: sqlite3 v1.5.4
 - mobile技术栈: 采用react native并集成go mobile编译的一部分可共用的go程序代码.
   - mobile代码: react native 0.79.2, kotlin
   - 应用: go 1.26.5, go mobile
   - 数据库: sqlite3 v1.5.4

## 代码组织说明

### 共用文件
不同应用共用的文件

**文件目录列表**:
- gobase : go程序应用入口模块
- internal: go程序内部模块,主要为应用的业务逻辑
- mocks: go程序部分的mock文件,用于测试
- pkg: go程序的基础框架模块代码,包括但不限于,运行引擎模块,引导模块,配置模块,网络模块,数据仓库模块,web模块,ptp模块等.
  -  runtime: 运行引擎模块,负责实现对模块分类初始化调用.
  -  bootstrap: 引导模块,负责实现对模块初始化延迟调用,销毁调用,以及就绪调用.
  -  config: 配置模块,负责实现对配置的读取和写入,以及配置更新的监听和通知.
  -  repository: 数据仓库模块,负责实现对数据库的构建和操作,以及对数据层的结构定义.
  -  web: web模块,负责实现对web使用的服务器.
  -  ptp: ptp模块,负责实现对ptp使用的服务器,ptp客户端,设备发现(如:广播,地址手册等).
  -  servlet: 负责实现对请求响应处理的服务功能(MVC架构). 
  -  app: 应用服务模块,负责实现对请求响应处理的集成服务功能(MVC架构).
  -  log: 日志模块,负责实现对日志的记录和管理.
  -  proto: 包含用于对protobuf进行处理的功能方法.
  -  module: 包含简单的runtime模块类定义,包含base模块类以及sub子模块类定义
- packages: typescript 共用模块,包括但不限于,数据结构定义模块,假数据生成模块.

### Web应用

**文件目录列表**:
- web: 前端代码目录
- main.go: web应用入口文件
- script/build_web.go: 构建web前端文件,包括不限于,html,css,js等.
- 其他请参考[共用文件](#共用文件)部分

### Mobile应用

**文件目录列表**:
- mobile:  mobile代码目录,包含react native代码以及kotlin代码.
- gomobile: go mobile使用的代码目录,与mobile代码集成的go程序代码.
- script/build_mobile.go:  用于将go代码编译为mobile代码可使用程序文件(android为aar文件).
- 其他请参考[共用文件](#共用文件)部分 


## Go代码组织规范
- go代码按照业务或功能模块进行组织.
- pkg下模块为应用功能模块,用于应用组织基础功能或提供给internal下的业务模块使用.需要尽可能保证独立性,不依赖internal下的模块.
- internal下为业务逻辑模块,用于实现应用业务逻辑.



   