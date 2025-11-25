build:
	@go build -gcflags "all=-N -l"
	@mkdir -p ".out"

PROTO_BASE="../../../proto"

debug: build
	@protoc --plugin=protoc-gen-dict=./debug.sh \
		-I $(PROTO_BASE) \
		-I ~/.cache/buf/v3/modules/b5/buf.build/googleapis/googleapis/61b203b9a9164be9a834f58c37be6f62/files \
		--dict_opt module=github.com/code-crafters-lab/ccl/pkg/grpc --dict_out .out \
		--dict_opt out_type="md" \
		$(PROTO_BASE)/ccl/license/*.proto