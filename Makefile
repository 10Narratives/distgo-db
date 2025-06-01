PROTO_DIR=proto
PROTO_OUT_DIR=pkg/proto
PROTOC_GEN_GO_VERSION=v1.28.1
PROTOC_GEN_GO_GRPC_VERSION=v1.2.0

# Путь к .proto файлу относительно корня проекта
PROTO_FILE=worker/database/v1/database_service.proto

.PHONY: proto
proto: ## Generate protobuf code
	@echo "Generating protobuf code..."
	@mkdir -p $(PROTO_OUT_DIR)
	protoc \
		-I=$(PROTO_DIR) \
		--go_out=$(PROTO_OUT_DIR) \
		--go_opt=paths=source_relative \
		--go-grpc_out=$(PROTO_OUT_DIR) \
		--go-grpc_opt=paths=source_relative,require_unimplemented_servers=false \
		$(PROTO_DIR)/$(PROTO_FILE)

.PHONY: proto-deps
proto-deps: ## Install protobuf dependencies
	@echo "Installing protoc-gen-go..."
	@go install google.golang.org/protobuf/cmd/protoc-gen-go@$(PROTOC_GEN_GO_VERSION)
	@echo "Installing protoc-gen-go-grpc..."
	@go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@$(PROTOC_GEN_GO_GRPC_VERSION)
	@echo "Dependencies installed. Make sure protoc is installed on your system."

.PHONY: proto-clean
proto-clean: ## Clean generated protobuf files
	@echo "Cleaning generated protobuf files..."
	@rm -rf $(PROTO_OUT_DIR)/$(dir $(PROTO_FILE))