## Goals

PAN-GO is a multi device files real-time browsing program. It enables a single device to access and collaboratively manage files across multiple devices in real-time.

Syncthing should be：

1. **Real time acquisition of files from multiple devices**
   Allow direct access to files from multiple different devices through a single device

1. **Easy to use and set up**
   Minimizes setup requirements by detecting users and only engaging in necessary interactions.

## Getting Started

Please refer to the [Getting Started Guide](Guide.md)

## Build

The following software is required for building. Please install them on the system where the build will be executed.

- [golang v1.23.4+](https://go.dev/doc/install)
- [node v20.17.0+](https://nodejs.org/en/download)

**Building the Application:**

```shell
   # Clone the source code to local host
   git clone https://github.com/kiba-zhao/pan-go.git
   # Navigate to source directory
   cd pan-go
   # Install dependency modules
   go mod tidy
   # Generate the console web page
   go generate
   # Build the application
   go build .
```

**Launch Program:**

```shell
   # Start in a Linux environment
   ./pan
```

## License

All code is licensed under the [MIT License](LICENSE).
