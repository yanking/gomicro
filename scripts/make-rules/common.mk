# Common variables used across Makefile rules

# Go modules
GO_CMD = go
GO_MOD = github.com/yanking/gomicro

# gRPC related tools
PROTOC_GEN_GO = google.golang.org/protobuf/cmd/protoc-gen-go
PROTOC_GEN_GO_GRPC = google.golang.org/grpc/cmd/protoc-gen-go-grpc

# Protobuf files directory
PROTO_DIR = api/proto
PROTO_FILES = $(wildcard $(PROTO_DIR)/*.proto)

# Output directory for generated files
PB_DIR = api/helloworld

# Swagger related tools
SWAGGER_CMD = swagger