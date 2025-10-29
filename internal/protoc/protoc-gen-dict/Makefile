build:
	@go build -gcflags "all=-N -l"

PROTO_BASE="../../../proto"

debug: build
	@protoc --plugin=protoc-gen-dict=./debug.sh \
		-I $(PROTO_BASE) \
		--dict_opt paths=source_relative --dict_out . \
		--dict_opt java_enum_package=net.jqsoft.cds.bid.enums \
		--dict_opt out_type="md|sql|java" \
		$(PROTO_BASE)/jqsoft/bid/*.proto