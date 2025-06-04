PROTO_DIR=proto
PROTO_OUT_DIR=pkg/proto

# Versions for Protobuf plugins
PROTOC_GEN_GO_VERSION=v1.28.1
PROTOC_GEN_GO_GRPC_VERSION=v1.3.0
PROTOC_GEN_GO_VALIDATE_VERSION=v1.2.1

# Path to your .proto file relative to project root
PROTO_FILE=worker/database/v1/database_service.proto

.PHONY: proto
proto: ## Generate protobuf code with validation support
	@echo "Generating protobuf code with validation..."
	@mkdir -p $(PROTO_OUT_DIR)

    # Determine path to protoc-gen-validate .proto files
	VALIDATE_PATH=$(shell go list -f '{{$$src := .Dir}}{{range .GoModReplace}}{{if eq .New.Name "github.com/envoyproxy/protoc-gen-validate"}}{{$src}}/{{.Path}}{{end}}{{end}}' -m github.com/envoyproxy/protoc-gen-validate)
	protoc \
	    -I=$(PROTO_DIR) \
	    -I=$(VALIDATE_PATH) \
	    --go_out=$(PROTO_OUT_DIR) \
	    --go_opt=paths=source_relative \
	    --go-grpc_out=$(PROTO_OUT_DIR) \
		--go-grpc_opt=paths=source_relative \
	    --validate_out="lang=go:$(PROTO_OUT_DIR)" \
	    $(PROTO_DIR)/$(PROTO_FILE)

.PHONY: proto-deps
proto-deps: ## Install required protobuf tools including validate plugin
	@echo "Installing protoc-gen-go..."
	@go install google.golang.org/protobuf/cmd/protoc-gen-go@$(PROTOC_GEN_GO_VERSION)
	@echo "Installing protoc-gen-go-grpc..."
	@go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@$(PROTOC_GEN_GO_GRPC_VERSION)
	@echo "Installing protoc-gen-go-validate from envoyproxy fork..."
	@go install github.com/envoyproxy/protoc-gen-validate@$(PROTOC_GEN_GO_VALIDATE_VERSION)
	@echo "Dependencies installed. Ensure 'protoc' is installed on your system."

.PHONY: proto-clean
proto-clean: ## Clean generated protobuf files
	@echo "Cleaning generated protobuf files..."
	@rm -rf $(PROTO_OUT_DIR)/$(dir $(PROTO_FILE))