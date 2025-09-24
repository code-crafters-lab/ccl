package main

import (
	"buf.build/go/protovalidate"
	"github.com/code-crafters-lab/ccl/pkg/grpc/category/v1"
	"github.com/golang/protobuf/proto"
)

func main() {

	category := v1.Category{
		Id:          1,
		Pid:         2,
		Name:        "",
		Sort:        3,
		Description: "test",
	}

	v2 := proto.MessageV2(&category)
	if err := protovalidate.Validate(v2); err != nil {
		// Handle failure.
		print(err.Error())
	}

}
