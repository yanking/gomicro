# Swagger documentation rules

##@ Documentation

.PHONY: swagger-init
swagger-init: ## Initialize Swagger documentation
	@echo "Initializing Swagger documentation..."
	@mkdir -p ./docs/swagger
	@echo "Swagger documentation directory created successfully!"

.PHONY: swagger-generate
swagger-generate: ## Generate Swagger documentation from source code
	@echo "Generating Swagger documentation from source code..."
	$(SWAGGER_CMD) generate spec -o ./docs/swagger/swagger.json ./examples/...
	@echo "Swagger documentation generated successfully!"

.PHONY: swagger-validate
swagger-validate: ## Validate Swagger specification
	@echo "Validating Swagger specification..."
	$(SWAGGER_CMD) validate ./docs/swagger/swagger.json
	@echo "Swagger specification validated successfully!"

.PHONY: swagger-serve
swagger-serve: ## Serve Swagger documentation
	@echo "Serving Swagger documentation..."
	@echo "Documentation will be available at http://localhost:63150/docs"
	@echo "Press Ctrl+C to stop"
	$(SWAGGER_CMD) serve -F=redoc ./docs/swagger/swagger.json

.PHONY: check-swagger
check-swagger: ## Check if swagger tool is installed
	@if command -v $(SWAGGER_CMD) >/dev/null 2>&1; then \
		echo "swagger: installed ($$($(SWAGGER_CMD) version))"; \
	else \
		echo "swagger: not found"; \
		echo "Run 'make install-swagger' to install it"; \
	fi