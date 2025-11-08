package main

import (
	"context"
	"log"
	"log/slog"
	"os"
	"time"

	"github.com/yanking/gomicro/api/helloworld"
	"github.com/yanking/gomicro/pkg/transport/rpc"
	"github.com/yanking/gomicro/pkg/transport/rpc/serverinterceptors"
)

// server is used to implement helloworld.GreeterServer.
type server struct {
	helloworld.UnimplementedGreeterServer
}

// SayHello implements helloworld.GreeterServer
func (s *server) SayHello(ctx context.Context, in *helloworld.HelloRequest) (*helloworld.HelloReply, error) {
	log.Printf("Received: %v", in.GetName())

	// 模拟处理时间
	time.Sleep(100 * time.Millisecond)

	return &helloworld.HelloReply{Message: "Hello " + in.GetName()}, nil
}

func main() {
	// 创建日志记录器
	logger := slog.New(slog.NewTextHandler(os.Stdout, nil))

	// 创建RPC服务器并添加拦截器
	rpcServer := rpc.NewServer(logger,
		rpc.WithAddress(":50051"),
		rpc.WithHealthz(true),
		rpc.WithReflection(true),
		rpc.WithUnaryInterceptors(
			serverinterceptors.RecoveryInterceptor(logger),
			serverinterceptors.LoggingInterceptor(logger),
			serverinterceptors.MetricsInterceptor(),
			serverinterceptors.TimeoutInterceptor(5*time.Second),
		),
		rpc.WithStreamInterceptors(
			serverinterceptors.LoggingStreamInterceptor(logger),
		),
	)

	// 注册服务
	helloworld.RegisterGreeterServer(rpcServer.Server, &server{})

	// 启动服务器
	log.Printf("server listening at %v", ":50051")
	if err := rpcServer.Start(context.Background()); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
