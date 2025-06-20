PROTO_DIR=proto/
PROTO_OUT_DIR=pkg/proto
GOOGLEAPIS=/tmp/googleapis
PGV_DIR=/tmp/protoc-gen-validate

clean:
	@echo "Cleaning generated protobuf files..."
	@rm -rf $(PROTO_OUT_DIR)

install:
	@echo "Installing dependencies..."
	@go install google.golang.org/protobuf/cmd/protoc-gen-go@latest
	@go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@latest
	@go install github.com/envoyproxy/protoc-gen-validate@latest

deps:
	@[ -d $(GOOGLEAPIS) ] || git clone https://github.com/googleapis/googleapis.git     $(GOOGLEAPIS)
	@[ -d $(PGV_DIR) ] || git clone https://github.com/envoyproxy/protoc-gen-validate.git     $(PGV_DIR)

generate: clean install deps
	@echo "Generating files..."
	@mkdir -p $(PROTO_OUT_DIR)
	protoc -I$(PROTO_DIR) \
			-I$(GOOGLEAPIS) \
			-I$(PGV_DIR) \
		$(PROTO_DIR)worker/database/v1/*.proto $(PROTO_DIR)worker/replication/v1/*.proto \
		$(PROTO_DIR)master/cluster/v1/*.proto \
		--go_out=$(PROTO_OUT_DIR) --go_opt=paths=source_relative \
		--go-grpc_out=$(PROTO_OUT_DIR) --go-grpc_opt=paths=source_relative \
		--validate_out="lang=go,paths=source_relative:$(PROTO_OUT_DIR)"

run:
	@go run cmd/worker/main.go --config config/worker.yaml

help:
	@echo "Usage: make [target]"
	@echo ""
	@echo "Available targets:"
	@echo "  clean      - Remove all generated protobuf files."
	@echo "  install    - Install required Go tools (protoc-gen-go, protoc-gen-go-grpc, protoc-gen-validate)."
	@echo "  deps       - Clone required repositories (googleapis and protoc-gen-validate) if they don't exist."
	@echo "  generate   - Generate Go protobuf, gRPC, and validation code from .proto files."
	@echo "  run        - Run the worker application using the specified configuration file."
	@echo "  help       - Display this help message."

.PHONY: clean install deps generate run help