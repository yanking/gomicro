package main

import (
	"context"
	"log"
	"log/slog"
	"os"
	"time"

	"github.com/yanking/gomicro/api/helloworld"
	"github.com/yanking/gomicro/pkg/transport/rpc"
	"github.com/yanking/gomicro/pkg/transport/rpc/clientinterceptors"
)

func main() {
	// 创建日志记录器
	logger := slog.New(slog.NewTextHandler(os.Stdout, nil))

	// 创建RPC客户端并添加拦截器
	client, err := rpc.NewClient(logger, "localhost:50051",
		rpc.WithInsecure(),
		rpc.WithClientUnaryInterceptors(
			clientinterceptors.LoggingInterceptor(logger),
			clientinterceptors.TimeoutInterceptor(5*time.Second),
			clientinterceptors.RetryInterceptor(3, time.Second),
			clientinterceptors.MetricsInterceptor(),
		),
		rpc.WithClientStreamInterceptors(
			clientinterceptors.LoggingStreamInterceptor(logger),
		),
	)
	if err != nil {
		log.Fatal("Failed to create RPC client:", err)
	}
	defer client.Close()

	// 使用客户端
	greeterClient := helloworld.NewGreeterClient(client.GetConn())

	// 调用远程方法
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	resp, err := greeterClient.SayHello(ctx, &helloworld.HelloRequest{Name: "world"})
	if err != nil {
		log.Fatal("Failed to call SayHello:", err)
	}

	log.Printf("Response: %s", resp.GetMessage())
}
