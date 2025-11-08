package main

import (
	"context"
	"log"
	"log/slog"
	"os"

	"github.com/yanking/gomicro/api/helloworld"
	"github.com/yanking/gomicro/pkg/transport/rpc"
)

// server is used to implement helloworld.GreeterServer.
type server struct {
	helloworld.UnimplementedGreeterServer
}

// SayHello implements helloworld.GreeterServer
func (s *server) SayHello(ctx context.Context, in *helloworld.HelloRequest) (*helloworld.HelloReply, error) {
	log.Printf("Received: %v", in.GetName())
	return &helloworld.HelloReply{Message: "Hello " + in.GetName()}, nil
}

func main() {
	// 创建日志记录器
	logger := slog.New(slog.NewTextHandler(os.Stdout, nil))

	// 创建gRPC服务器
	grpcServer := rpc.NewServer(logger,
		rpc.WithAddress(":50051"),
		rpc.WithHealthz(true),
		rpc.WithReflection(true),
	)

	// 注册服务
	helloworld.RegisterGreeterServer(grpcServer.Server, &server{})

	// 启动服务器
	log.Printf("server listening at %v", ":50051")
	if err := grpcServer.Start(context.Background()); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
