# protoc-gen-rbs
Protobuf compiler plugin that generates `.rbs` "Ruby Type Specification" files.

### Installation

```
go get github.com/fundthmcalculus/protoc-gen-rbs
```

### Usage

```
protoc --rbs_out=. example.proto
```

To disable generation of gRPC `.rbs` files, use the `grpc=false` option:

```
protoc --rbs_out=grpc=false:. example.proto
```

### Example

For the input [example.proto](testdata/example.proto):
 - [example_pb.rbi](testdata/example_pb.rbi) contains the message(s) interface
 - [example_services_pb.rbi](testdata/example_services_pb.rbi) contains the service(s) interface
