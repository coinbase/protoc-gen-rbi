# protoc-gen-rbi
Protobuf compiler plugin that generates [Sorbet](https://sorbet.org/) `.rbi` "Ruby Interface" files.

### Installation

```
go get github.com/sorbet/protoc-gen-rbi
```

### Usage

```
protoc --rbi_out=. example.proto
```

To disable generation of gRPC `.rbi` files, use the `grpc=false` option:

```
protoc --rbi_out=grpc=false:. example.proto
```

### Example

For the input [example.proto](testdata/example.proto):
 - [example_pb.rbi](testdata/example_pb.rbi) contains the message(s) interface
 - [example_services_pb.rbi](testdata/example_services_pb.rbi) contains the service(s) interface
