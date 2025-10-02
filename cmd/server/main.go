package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"time"

	"connectrpc.com/connect"
	"connectrpc.com/validate"
	categoryv1 "github.com/code-crafters-lab/ccl/pkg/grpc/category/v1"
	categoryv1connect "github.com/code-crafters-lab/ccl/pkg/grpc/category/v1/v1connect"
	"golang.org/x/net/http2"
	"golang.org/x/net/http2/h2c"
)

type categoryServer struct {
	categoryv1connect.UnimplementedCategoryServiceHandler
}

func (c *categoryServer) CreateCategory(ctx context.Context, req *categoryv1.CreateCategoryRequest) (*categoryv1.CreateCategoryResponse, error) {
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

func (c *categoryServer) ListCategory(ctx context.Context, req *categoryv1.ListCategoryRequest) (*categoryv1.ListCategoryResponse, error) {
	log.Println("Request headers: ", req)
	res := &categoryv1.ListCategoryResponse{}
	res.Total = 3
	for i := 0; i < 3; i++ {
		res.Categories = append(res.Categories, &categoryv1.Category{
			Id:   uint64(i + 1),
			Pid:  uint64(i + 10001),
			Name: fmt.Sprintf("name_%2d", i+1),

			Sort: uint32(i),
		})
	}

	return res, nil
}

// 服务核心配置（集中管理，便于修改）
const (
	ServiceName     = "user-service"
	ServerPort      = 8090
	APIPrefix       = "/api/v1" // 统一路径前缀
	NacosAddr       = "127.0.0.1:8848"
	NacosNamespace  = ""
	ReadTimeout     = 5 * time.Second  // 读取请求超时
	WriteTimeout    = 10 * time.Second // 写入响应超时
	ShutdownTimeout = 5 * time.Second  // 优雅关闭超时
)

func main() {
	server := &categoryServer{}
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

	handler = CORSInterceptor(handler)

	// 支持原始 grpc 调用
	mux.Handle(svcPrefix, handler)
	// 4. 配置 HTTP 路由（带路径前缀）

	// 统一添加路径前缀 /api/v1
	mux.Handle(APIPrefix+"/", http.StripPrefix(APIPrefix, handler))

	// 4. 使用 CORS 中间件包装你的 mux
	// 中间件的执行顺序很重要，CORS 中间件应该在最外层
	//wrappedMux := corsMiddleware(mux)

	http.ListenAndServe(
		fmt.Sprintf(":%d", ServerPort),
		// Use h2c so we can serve HTTP/2 without TLS.
		h2c.NewHandler(mux, &http2.Server{}),
	)

}

func corsUnaryInterceptor() connect.UnaryInterceptorFunc {
	return func(next connect.UnaryFunc) connect.UnaryFunc {
		return func(ctx context.Context, req connect.AnyRequest) (connect.AnyResponse, error) {
			// --------------------------
			// 阶段 1：处理 OPTIONS 预检请求
			// --------------------------
			if req.HTTPMethod() == http.MethodOptions {
				// 构建 CORS 预检响应
				respHeaders := http.Header{}
				setCORSHeaders(respHeaders, req.Header())

				// 返回 204 No Content（预检请求无需响应体）
				response := connect.NewResponse[string](nil)
				return response, nil
			}
			// --------------------------
			// 阶段 2：处理正常 RPC 请求（如 POST）
			// --------------------------
			// 先执行业务逻辑（如 Greet 方法）
			resp, err := next(ctx, req)
			if err != nil {

				return resp, err
			}
			// 为 RPC 响应注入 CORS 头
			setCORSHeaders(resp.Header(), req.Header())
			return resp, nil
		}
	}
}

// 1. 定义 CORS 配置（生产环境需替换为实际域名）
var corsConfig = struct {
	AllowedOrigins   []string // 允许的跨域来源
	AllowedMethods   []string // 允许的 HTTP 方法（Connect 基于 POST）
	AllowedHeaders   []string // 允许的请求头
	ExposedHeaders   []string // 暴露给前端的响应头
	AllowCredentials bool     // 是否允许跨域携带 Cookie
	MaxAge           int      // 预检请求缓存时间（秒）
}{
	AllowedOrigins:   []string{"https://your-frontend.com", "http://localhost:8090"}, // 开发+生产域名
	AllowedMethods:   []string{"POST", "OPTIONS"},                                    // Connect 核心用 POST，预检用 OPTIONS
	AllowedHeaders:   []string{"Content-Type", "Authorization"},
	ExposedHeaders:   []string{"X-RPC-ID"},
	AllowCredentials: true,
	MaxAge:           86400, // 1 天
}

// setCORSHeaders：根据请求头动态设置 CORS 响应头
func setCORSHeaders(respHeaders, reqHeaders http.Header) {
	// 1. 处理 Allow-Origin（支持Credentials时，不能用 *）
	//origin := reqHeaders.Get("Origin")
	//if origin != "" && isAllowedOrigin(origin) {
	respHeaders.Set("Access-Control-Allow-Origin", "*")
	//}

	// 2. 处理 Allow-Credentials
	//if corsConfig.AllowCredentials {
	//respHeaders.Set("Access-Control-Allow-Credentials", "true")
	//}

	// 3. 处理预检请求的 Allow-Methods/Allow-Headers
	//accessControlRequestMethod := reqHeaders.Get("Access-Control-Request-Method")
	//if accessControlRequestMethod != "" {
	respHeaders.Set("Access-Control-Allow-Methods", "*")
	//}

	//accessControlRequestHeaders := reqHeaders.Get("Access-Control-Request-Headers")
	//if accessControlRequestHeaders != "" {
	respHeaders.Set("Access-Control-Allow-Headers", "*")
	//}

	// 4. 暴露自定义响应头
	//if len(corsConfig.ExposedHeaders) > 0 {
	respHeaders.Set("Access-Control-Expose-Headers", "*")
	//}

	// 5. 预检请求缓存时间
	if corsConfig.MaxAge > 0 {
		respHeaders.Set("Access-Control-Max-Age", string(rune(corsConfig.MaxAge)))
	}
}

// isAllowedOrigin：检查请求 Origin 是否在允许列表中
func isAllowedOrigin(origin string) bool {
	for _, allowed := range corsConfig.AllowedOrigins {
		if allowed == "*" || allowed == origin {
			return true
		}
	}
	return false
}
