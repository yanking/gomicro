# Tools installation rules

##@ Development

.PHONY: install-tools
install-tools: ## Install gRPC development tools (protoc-gen-go and protoc-gen-go-grpc)
	@echo "Installing gRPC development tools..."
	$(GO_CMD) install $(PROTOC_GEN_GO)@latest
	$(GO_CMD) install $(PROTOC_GEN_GO_GRPC)@latest
	@echo "gRPC development tools installed successfully!"

.PHONY: install-all-tools
install-all-tools: install-tools install-swagger ## Install all development tools (gRPC + Swagger)

.PHONY: install-protoc
install-protoc: ## Install protoc compiler (package manager based)
	@echo "Installing protoc compiler..."
	@if command -v brew >/dev/null 2>&1; then \
		brew install protobuf; \
	elif command -v apt-get >/dev/null 2>&1; then \
		sudo apt-get update && sudo apt-get install -y protobuf-compiler; \
	elif command -v yum >/dev/null 2>&1; then \
		sudo yum install -y protobuf-compiler; \
	else \
		echo "Note: The protoc compiler cannot be installed directly via Go."; \
		echo "Please install protoc manually from https://github.com/protocolbuffers/protobuf/releases"; \
		echo "Alternatively, you can install it using a package manager:"; \
		echo "  - macOS: brew install protobuf"; \
		echo "  - Ubuntu/Debian: sudo apt-get install protobuf-compiler"; \
		echo "  - CentOS/RHEL: sudo yum install protobuf-compiler"; \
		exit 1; \
	fi
	@echo "protoc compiler installed successfully!"

.PHONY: install-swagger
install-swagger: ## Install swagger tool for Swagger documentation generation
	@echo "Installing swagger tool..."
	$(GO_CMD) install github.com/go-swagger/go-swagger/cmd/swagger@latest
	@echo "swagger tool installed successfully!"

.PHONY: init
init: install-tools ## Initialize development environment

.PHONY: check-tools
check-tools: ## Check if required tools are installed
	@if command -v protoc >/dev/null 2>&1; then \
		echo "protoc: $$(protoc --version)"; \
	else \
		echo "protoc: not found"; \
	fi
	@if command -v protoc-gen-go >/dev/null 2>&1; then \
		echo "protoc-gen-go: installed"; \
	else \
		echo "protoc-gen-go: not found"; \
	fi
	@if command -v protoc-gen-go-grpc >/dev/null 2>&1; then \
		echo "protoc-gen-go-grpc: installed"; \
	else \
		echo "protoc-gen-go-grpc: not found"; \
	fi
	@if command -v $(SWAGGER_CMD) >/dev/null 2>&1; then \
		echo "swagger: installed ($$($(SWAGGER_CMD) version))"; \
	else \
		echo "swagger: not found"; \
	fi