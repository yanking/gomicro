# Clean rules

##@ Cleanup

.PHONY: clean
clean: ## Clean generated files
	@echo "Cleaning generated files..."
	find $(PB_DIR) -name "*.pb.go" -type f -delete
	@echo "Generated files cleaned!"

# Output directory for generated files
PB_DIR = api/helloworld

