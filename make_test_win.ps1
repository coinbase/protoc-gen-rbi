go build
protoc --proto_path=testdata --rbi_out=grpc=true:testdata --plugin=protoc-gen-rbi=C:\personal\protoc-gen-rbi\protoc-gen-rbi.exe .\testdata\*.proto
rbs parse testdata/*.rbs