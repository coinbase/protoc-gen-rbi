module github.com/coinbase/protoc-gen-rbi

go 1.13

require (
	github.com/lyft/protoc-gen-star/v2 v2.0.3
	github.com/spf13/afero v1.11.0 // indirect
	golang.org/x/tools v0.24.0 // indirect
	google.golang.org/protobuf v1.34.2
)

replace github.com/lyft/protoc-gen-star/v2 => ./protoc-gen-star
