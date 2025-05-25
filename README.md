# distgo-db

Distributed document-oriented database written in Go

## Protobuf Code Generation

This project uses Protocol Buffers (Protobuf) to define APIs and generate Go code. The `Makefile` automates the process of generating code from `.proto` files.

### Prerequisites

Before generating code, ensure the following tools and utilities are installed on your system:

1. Install `protoc` via package manager of download it for the [official releases](https://github.com/protocolbuffers/protobuf/releases?spm=a2ty_o01.29997173.0.0.176dc921maxc15).

2. Ensure `Go` is installed and configured (`go version` should return a valid version). Download Go from the [official website](https://go.dev/dl/?spm=a2ty_o01.29997173.0.0.176dc921maxc15)

3. `git` is required for cloning repositories

### Protobuf Plugins

These plugins are automatically installed by the Makefile, but they require Go to be installed:

- protoc-gen-go : Generates Go code for Protobuf messages.
- protoc-gen-go-grpc : Generates gRPC service stubs in Go.
- protoc-gen-validate : Generates validation logic for Protobuf messages.

### Cloned Repositories

The following repositories are automatically cloned by the Makefile:

- googleapis: Contains Google API definitions.
- protoc-gen-validate: Contains validation rules for Protobuf messages.

### Code generation

#### Set Up the Project

Run the setup target to install required tools and clone dependencies:

```bash
make setup
```

This will:

- Install protoc-gen-go, protoc-gen-go-grpc, and protoc-gen-validate into the bin directory.
- Clone the googleapis and protoc-gen-validate repositories into the vendor directory.

#### Generate Code

Run the generate target to generate Go code:

```bash
make generate
```

This will:

- Process all .proto files in the proto/distgodb/worker/document/v1/ directory.
- Generate Protobuf message definitions (*.pb.go), gRPC service stubs (*_grpc.pb.go), and validation logic (*.validate.pb.go) in the pkg/proto/distgodb/worker/document/v1/ directory.

#### Clean Up

To remove all generated files, cloned repositories, and installed tools, run:

```bash
make clean
```
