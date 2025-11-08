# Proto generation rules

##@ Generation

.PHONY: gen-proto
gen-proto: ## Generate gRPC code from proto files
	@echo "Generating gRPC code from proto files..."
	@for proto in $(PROTO_FILES); do \
		echo "Generating code for $$proto"; \
		protoc --go_out=$(PB_DIR) --go_opt=paths=source_relative \
		       --go-grpc_out=$(PB_DIR) --go-grpc_opt=paths=source_relative \
		       $$proto; \
	done
	@echo "gRPC code generated successfully!"

.PHONY: gen-helloworld
gen-helloworld: ## Generate gRPC code for helloworld.proto
	@echo "Generating gRPC code for helloworld.proto..."
	protoc --go_out=$(PB_DIR) --go_opt=paths=source_relative \
	       --go-grpc_out=$(PB_DIR) --go-grpc_opt=paths=source_relative \
	       $(PROTO_DIR)/helloworld.proto
	@echo "helloworld gRPC code generated successfully!"

.PHONY: check-proto
check-proto: ## Check if proto files exist
	@if [ -z "$(PROTO_FILES)" ]; then \
		echo "No proto files found in $(PROTO_DIR)"; \
		exit 1; \
	else \
		echo "Found proto files:"; \
		for proto in $(PROTO_FILES); do \
			echo "  - $$proto"; \
		done; \
	fi