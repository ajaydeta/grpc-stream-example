# Golang gRPC Stream Example

This is a simple example of a gRPC stream in Golang.

## Abstract
### what is gRPC?
gRPC is an open-source remote procedure call (RPC) framework developed by Google. RPC is a communication protocol that enables applications on different systems to communicate with each other and invoke functions or procedures remotely, as if they were on the same machine.

In simple terms, gRPC allows applications to talk to each other and make remote function calls efficiently and uniformly. It leverages Protocol Buffers (protobuf) technology to define message formats and service interfaces.

### what is stream?
gRPC supports streaming semantics, where either the client or the server (or both) send a stream of messages on a single RPC call. The most common types are unary (client sends a single request and gets a single response) and server streaming (client sends a single request and gets a stream of responses) or vice versa.

## How to run
### Prerequisites
- [Go](https://golang.org/doc/install)
- [Protocol Buffers](https://grpc.io/docs/protoc-installation/)
- [gRPC](https://grpc.io/docs/languages/go/quickstart/)

### Run
1. Clone this repository
```bash
git clone https://github.com/ajaydeta/grpc-stream-example.git
```

2. Install dependencies
```bash
make dependencies
```

3. Generate protobuf
```bash
make proto
```

4. Run server
- you can run server with this command:
```bash
make server
```
- you can also run server and access it using client with this command:
```bash
make run
```

## References
- [gRPC concepts](https://grpc.io/docs/what-is-grpc/core-concepts/)
- [gRPC Go Quick Start](https://grpc.io/docs/languages/go/quickstart/)
- [gRPC Go Tutorial](https://grpc.io/docs/languages/go/basics/)