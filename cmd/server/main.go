package main

import (
	"context"
	"log"
	"net/http"

	"connectrpc.com/connect"
	categoryv1 "github.com/code-crafters-lab/ccl/pkg/grpc/category/v1"
	"github.com/code-crafters-lab/ccl/pkg/grpc/category/v1/v1connect"
	"golang.org/x/net/http2"
	"golang.org/x/net/http2/h2c"
)

type CategoryServer struct{}

func (c *CategoryServer) CreateCategory(ctx context.Context, req *connect.Request[categoryv1.CreateCategoryRequest]) (*connect.Response[categoryv1.CreateCategoryResponse], error) {
	log.Println("Request headers: ", req.Header())
	msg := req.Msg
	log.Println("Received message:", msg)

	ca := &categoryv1.Category{
		Id:   1,
		Pid:  msg.Pid,
		Name: msg.Name,
		Sort: msg.Sort,
	}
	res := connect.NewResponse(&categoryv1.CreateCategoryResponse{
		Category: ca,
	})

	res.Header().Set("Greet-Version", "v1")
	return res, nil
}

func main() {
	server := &CategoryServer{}
	mux := http.NewServeMux()
	path, handler := v1connect.NewCategoryUserServiceHandler(server)
	mux.Handle(path, handler)
	http.ListenAndServe(
		":8080",
		// Use h2c so we can serve HTTP/2 without TLS.
		h2c.NewHandler(mux, &http2.Server{}),
	)

}
