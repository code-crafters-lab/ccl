build:
	@go build -gcflags "all=-N -l"

PROTO_BASE="../../../proto"

debug: build
	@protoc --plugin=protoc-gen-dict=./debug.sh \
		-I $(PROTO_BASE) \
		-I ~/.cache/buf/v3/modules/b5/buf.build/googleapis/googleapis/61b203b9a9164be9a834f58c37be6f62/files \
		--dict_opt paths=source_relative --dict_out . \
		--dict_opt out_type="md|sql|java" \
		$(PROTO_BASE)/jqsoft/bid/*.proto