build:
	@go build -gcflags "all=-N -l"

debug: build
	@protoc --plugin=protoc-gen-dict=./debug.sh -I ../../../proto --dict_out . ../../../proto/resources/resource.proto