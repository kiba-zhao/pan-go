## 目标

本项目旨在实现根据与用户文字聊天或语音指令,对设备进行自动操作以及提示说明.支持多设备协作.

## 入门

请查看[入门手册](Guide.md)

## 手动构建

构建需要以下软件，请安装到执行构建的系统。

- [golang v1.26.5+](https://go.dev/doc/install)
- [node v20.17.0+](https://nodejs.org/zh-cn/download)

**构建应用：**

```shell
   # 克隆源码到本地
   git clone https://github.com/kiba-zhao/pan-go.git
   # 进入源码目录
   cd pan-go
   # 安装依赖模块
   go mod tidy
   # 生成控制台页面
   go generate
   # 构建应用
   go build .
```

**启动应用：**

```shell
   # linux环境启动
   ./pan
```

## 许可

这是 [MIT 许可证](LICENSE)条款下的自由软件。
