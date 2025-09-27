package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"time"

	"connectrpc.com/connect"
	categoryv1 "github.com/code-crafters-lab/ccl/internal/gen/category/v1"
	categoryv1connect "github.com/code-crafters-lab/ccl/internal/gen/category/v1/v1connect"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func main() {
	callGrpc()
}

func connectRpc() {
	client := categoryv1connect.NewCategoryServiceClient(
		http.DefaultClient,
		"http://localhost:8080/api/v1",
		connect.WithGRPC(),
	)
	res, err := client.CreateCategory(
		context.Background(),
		&categoryv1.CreateCategoryRequest{Name: "Jane"},
	)
	if err != nil {
		log.Println(err)
		return
	}
	log.Println(res)
}

func callGrpc() {
	// --- 1. 创建 gRPC 连接 ---
	// 目标服务器地址
	const address = "localhost:8080"

	// 创建连接选项。对于开发环境，可以使用不安全的连接。
	// 在生产环境中，应使用 WithTransportCredentials 配置 TLS。
	conn, err := grpc.NewClient(address, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("无法连接到 gRPC 服务器: %v", err)
	}
	// 确保在程序退出时关闭连接
	defer conn.Close()

	// --- 2. 创建客户端实例 ---
	client := categoryv1.NewCategoryServiceClient(conn)

	// --- 3. 调用 RPC 方法 ---

	// --- 示例 A: 调用 Unary RPC ---
	callUnaryRPC(client)

	//// --- 示例 B: 调用 Server Streaming RPC ---
	//callServerStreamingRPC(client)
	//
	//// --- 示例 C: 调用 Client Streaming RPC ---
	//callClientStreamingRPC(client)
	//
	//// --- 示例 D: 调用 Bidirectional Streaming RPC ---
	//callBidirectionalStreamingRPC(client)
}

// callUnaryRPC 演示如何调用一个简单的请求-响应 RPC。
func callUnaryRPC(client categoryv1.CategoryServiceClient) {
	fmt.Println("--- 调用 Unary RPC: SayHello ---")
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
	defer cancel()

	response, err := client.CreateCategory(ctx, &categoryv1.CreateCategoryRequest{Name: "Alice"})
	if err != nil {
		log.Printf("Unary RPC 调用失败: %v", err)
		return
	}
	log.Printf("收到响应: %s", response)
}
