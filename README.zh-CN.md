## 目标

PAN-GO 是一个的多设备文件实时浏览程序。它使得单个设备能够实时访问和协同管理多个设备上的文件。

我将努力逐一实现以下目标：

1. **即时获取多设备文件**
   在任一设备实时获取多个设备上的文件。

1. **易于使用**
   不进行过多的设置，只进行必要的用户交互。

## 入门

请查看[入门手册](Guide.md)

## 手动构建

构建需要以下软件，请安装到执行构建的系统。

- [golang v1.23.4+](https://go.dev/doc/install)
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
