build:
	@go build -gcflags "all=-N -l"

debug: build
	@protoc --plugin=protoc-gen-dict=./debug.sh \
		-I ../../../proto \
		--go_opt paths=source_relative --dict_out . \
		../../../proto/resources/*.proto