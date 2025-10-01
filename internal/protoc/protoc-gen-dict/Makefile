build:
	@go build -gcflags "all=-N -l"

PROTO_BASE="../../../proto"

debug: build
	@protoc --plugin=protoc-gen-dict=./debug.sh \
		-I $(PROTO_BASE) \
		--dict_opt paths=source_relative --dict_out . \
		$(PROTO_BASE)/ccl/resources/*.proto