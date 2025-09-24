GO_BIN := "$$(go env GOPATH)/bin"
#gen_authopt_path := "$(go_bin)/protoc-gen-authoption"
#gen_zitadel_path := "$(go_bin)/protoc-gen-zitadel"

NOW := $(shell date '+%Y-%m-%dT%T%z' | sed -E 's/.([0-9]{2})([0-9]{2})$$/-\1:\2/')
VERSION ?= development-$(now)
COMMIT_SHA ?= $(shell git rev-parse HEAD)

# Include versions of tools we build on-demand
include tools/env.mk
# This provides the "help" target.
include tools/help.mk

.PHONY: core_api
core_api:
	buf generate
	mkdir -p pkg/grpc
	cp -r .artifacts/grpc/github.com/code-crafters-lab/ccl/pkg/grpc/** pkg/grpc/
#	mkdir -p openapi/v2/zitadel
#	cp -r .artifacts/grpc/zitadel/ openapi/v2/zitadel

.PHONY: core_grpc_dependencies
core_grpc_dependencies:
	go install google.golang.org/protobuf/cmd/protoc-gen-go@v1.36.9 						# https://pkg.go.dev/google.golang.org/protobuf/cmd/protoc-gen-go?tab=versions
	go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@v1.5.1 						# https://pkg.go.dev/google.golang.org/grpc/cmd/protoc-gen-go-grpc?tab=versions
	go install github.com/grpc-ecosystem/grpc-gateway/v2/protoc-gen-grpc-gateway@v2.27.0	# https://pkg.go.dev/github.com/grpc-ecosystem/grpc-gateway/v2/protoc-gen-grpc-gateway?tab=versions
	go install github.com/grpc-ecosystem/grpc-gateway/v2/protoc-gen-openapiv2@v2.22.0 		# https://pkg.go.dev/github.com/grpc-ecosystem/grpc-gateway/v2/protoc-gen-openapiv2?tab=versions
	#go install github.com/envoyproxy/protoc-gen-validate@v1.2.1								# https://pkg.go.dev/github.com/envoyproxy/protoc-gen-validate?tab=versions
	#go install github.com/bufbuild/buf/cmd/buf@v1.45.0										# https://pkg.go.dev/github.com/bufbuild/buf/cmd/buf?tab=versions
	go install connectrpc.com/connect/cmd/protoc-gen-connect-go@v1.18.1						# https://pkg.go.dev/connectrpc.com/connect/cmd/protoc-gen-connect-go?tab=versions

lint:
	@buf lint

protoc:
	@protoc -I resources \
    	--go_out . \
    	--go_opt module=github.com/code-crafters-lab/idl \
    	resources/app.proto

db-test:
	mysqldump -h 10.1.83.26 -P 3306 -u teamwork --password='jqkj5350**)' -v \
	teamwork \
	t_user t_user_bind t_role t_user_role \
	t_system t_menu t_operation t_data_limit t_api_resource t_external_link \
	t_authority t_role_authority_rel \
	t_dict t_file \
	> teamwork-test.sql

db-prod:
	mysqldump -h 192.168.44.82 -P 3306 -u teamwork --password='teamwork_jqkj5350**)123' -v \
	teamwork \
	t_user t_user_bind t_role t_user_role \
	t_system t_menu t_operation t_data_limit t_api_resource t_external_link \
	t_authority t_role_authority_rel \
	t_dict t_file \
	> teamwork-prod.sql

init22:
	buf config init buf.build/ccl/dict -o dict

