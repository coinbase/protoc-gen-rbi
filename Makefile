init:
	bundle install

install:
	go install -mod=vendor .

vendor:
	go mod vendor

# Have to be in a separate directory, so we don't get "is already defined in file" errors
testbinary/%_bin.proto: testdata/%.proto
	mkdir -p $(shell dirname $@)
	cp $< $@

testdata/%_bin-descriptor-set.proto.bin: testbinary/%_bin.proto
	$(eval GRPC_TOOLS_LOCATION := $(shell bundle show grpc-tools))
	$(eval PROTOC_BINARY := $(GRPC_TOOLS_LOCATION)/bin/grpc_tools_ruby_protoc)
	$(PROTOC_BINARY) --descriptor_set_out=$@ $<

test: init install testdata/example_bin-descriptor-set.proto.bin
	$(eval PROTOS := $(shell cd testdata && find . -name "*.proto" | sed 's|^./||'))
	$(eval GRPC_TOOLS_LOCATION := $(shell bundle show grpc-tools))
	$(eval PROTOC_BINARY := $(GRPC_TOOLS_LOCATION)/bin/grpc_tools_ruby_protoc)
	$(eval GRPC_PLUGIN := $(GRPC_TOOLS_LOCATION)/bin/grpc_tools_ruby_protoc_plugin)
	$(PROTOC_BINARY) --proto_path=testdata --ruby_out=testdata $(PROTOS)
	$(PROTOC_BINARY) --proto_path=testdata --ruby_grpc_out=testdata --plugin=protoc-gen-ruby_grpc=$(GRPC_PLUGIN) $(PROTOS)
	$(PROTOC_BINARY) --proto_path=testdata --rbi_out=grpc=true:testdata $(PROTOS)
	$(PROTOC_BINARY) --proto_path=testdata --rbi_out=hide_common_methods=true:testdata/hide_common_methods $(PROTOS)
	$(PROTOC_BINARY) --proto_path=testdata --rbi_out=use_abstract_message=true:testdata/use_abstract_message $(PROTOS)
	$(PROTOC_BINARY) --proto_path=testdata --rbi_out=use_generic_proto_containers=true:testdata/use_generic_proto_containers $(PROTOS)
	$(PROTOC_BINARY) --proto_path=testdata --rbi_out=grpc=true,hide_common_methods=true,use_abstract_message=true,use_generic_proto_containers=true:testdata/all $(PROTOS)
	$(PROTOC_BINARY) --descriptor_set_in=testdata/example_bin-descriptor-set.proto.bin --rbi_out=. testbinary/example_bin.proto
	git diff --exit-code testdata testbinary
