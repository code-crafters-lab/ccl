package main

import (
	"context"
	"log"
	"net/http"

	"connectrpc.com/connect"
	v1 "github.com/code-crafters-lab/ccl/pkg/grpc/category/v1"
	"github.com/code-crafters-lab/ccl/pkg/grpc/category/v1/v1connect"
)

func main() {
	client := v1connect.NewCategoryUserServiceClient(
		http.DefaultClient,
		"http://localhost:8080",
		connect.WithGRPC(),
	)
	res, err := client.CreateCategory(
		context.Background(),
		connect.NewRequest(&v1.CreateCategoryRequest{Name: "Jane"}),
	)
	if err != nil {
		log.Println(err)
		return
	}
	log.Println(res.Msg)
}
