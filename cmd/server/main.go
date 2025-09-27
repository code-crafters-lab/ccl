package main

import (
	"context"
	"log"
	"net/http"
	"time"

	"connectrpc.com/connect"
	"connectrpc.com/validate"
	categoryv1 "github.com/code-crafters-lab/ccl/internal/gen/category/v1"
	categoryv1connect "github.com/code-crafters-lab/ccl/internal/gen/category/v1/v1connect"
	"golang.org/x/net/http2"
	"golang.org/x/net/http2/h2c"
)

type CategoryServer struct{}

func (c *CategoryServer) CreateCategory(ctx context.Context, req *categoryv1.CreateCategoryRequest) (*categoryv1.CreateCategoryResponse, error) {
	log.Println("Request headers: ", req)
	ca := &categoryv1.Category{
		Id:   1,
		Pid:  req.Pid,
		Name: req.Name,
		Sort: req.Sort,
	}
	res := &categoryv1.CreateCategoryResponse{
		Category: ca,
	}
	return res, nil
}

// 服务核心配置（集中管理，便于修改）
const (
	ServiceName     = "user-service"
	ServerPort      = 8080
	APIPrefix       = "/api/v1" // 统一路径前缀
	NacosAddr       = "127.0.0.1:8848"
	NacosNamespace  = ""
	ReadTimeout     = 5 * time.Second  // 读取请求超时
	WriteTimeout    = 10 * time.Second // 写入响应超时
	ShutdownTimeout = 5 * time.Second  // 优雅关闭超时
)

func main() {
	server := &CategoryServer{}
	mux := http.NewServeMux()

	// 暴露健康检查接口（Nacos/Envoy 探测）
	mux.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/plain")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("ok"))
	})

	interceptor, err := validate.NewInterceptor()
	if err != nil {
		log.Fatal(err)
	}
	interceptors := connect.WithInterceptors(interceptor)

	svcPrefix, handler := categoryv1connect.NewCategoryServiceHandler(server, interceptors)
	// 支持原始 grpc 调用
	mux.Handle(svcPrefix, handler)
	// 4. 配置 HTTP 路由（带路径前缀）

	// 统一添加路径前缀 /api/v1
	mux.Handle(APIPrefix+"/", http.StripPrefix(APIPrefix, handler))

	http.ListenAndServe(
		":8080",
		// Use h2c so we can serve HTTP/2 without TLS.
		h2c.NewHandler(mux, &http2.Server{}),
	)

}
