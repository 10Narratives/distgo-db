PROTOC = protoc
PROTO_DIR = proto
GO_MODULE_PATH = github.com/10Narratives/distgo-db/proto
OUTPUT_DIR = pkg/proto
GOOGLEAPIS_DIR = vendor/googleapis
PGV_DIR = vendor/protoc-gen-validate

PROTOC_GEN_GO := $(shell which protoc-gen-go)
PROTOC_GEN_GO_GRPC := $(shell which protoc-gen-go-grpc)
PROTOC_GEN_VALIDATE := $(shell which protoc-gen-validate)

PROTO_FILES = $(wildcard $(PROTO_DIR)/distgodb/worker/document/v1/*.proto)

all: setup generate

setup: install clone-repos
	@echo "Setup completed."

install:
	@echo "Installing required tools..."
	@GOBIN=$(PWD)/bin go install google.golang.org/protobuf/cmd/protoc-gen-go@latest
	@GOBIN=$(PWD)/bin go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@latest
	@GOBIN=$(PWD)/bin go install github.com/envoyproxy/protoc-gen-validate@latest
	@echo "Tools installed successfully."

clone-repos:
	@echo "Cloning required repositories..."
	@git clone https://github.com/googleapis/googleapis.git  $(GOOGLEAPIS_DIR) || true
	@git clone https://github.com/envoyproxy/protoc-gen-validate.git  $(PGV_DIR) || true
	@echo "Repositories cloned successfully."

generate: $(PROTO_FILES)
	@echo "Generating Go code from .proto files..."
	@mkdir -p $(OUTPUT_DIR)/distgodb/worker/document/v1
	$(PROTOC) \
	    --proto_path=$(PROTO_DIR) \
	    --proto_path=$(GOOGLEAPIS_DIR) \
	    --proto_path=$(PGV_DIR) \
	    --go_out=paths=source_relative:$(OUTPUT_DIR) \
	    --go-grpc_out=paths=source_relative:$(OUTPUT_DIR) \
	    --validate_out=lang=go,paths=source_relative:$(OUTPUT_DIR) \
	    $(PROTO_FILES)
	@echo "Code generation completed."

clean:
	@echo "Cleaning generated files..."
	rm -rf $(OUTPUT_DIR)
	rm -rf $(GOOGLEAPIS_DIR)
	rm -rf $(PGV_DIR)
	rm -rf bin/
	@echo "Cleanup completed."

.PHONY: all setup install clone-repos generate clean