PROTO_DIR := proto
PKG_DIR := pkg
GO_MOD_NAME := github.com/10Narratives/distgo-db

PROTO_INCLUDE := $(shell go list -m -f '{{ .Dir }}' google.golang.org/protobuf)/types/known

ifeq ($(wildcard $(PROTO_INCLUDE)),)
$(error "Path to Google Protobuf types not found. Please ensure google.golang.org/protobuf is installed.")
endif

all: generate

generate:
	@echo "Generating Go code for worker package..."
	mkdir -p $(PKG_DIR)/proto/worker
	protoc -I . proto/worker/*.proto --go_out=$(PKG_DIR) --go_opt=paths=source_relative \
		--go-grpc_out=$(PKG_DIR) --go-grpc_opt=paths=source_relative
	@echo "Code generation for worker package completed."

clean:
	@echo "Cleaning generated files..."
	rm -rf $(PKG_DIR)/proto
	@echo "Cleanup completed."

help:
	@echo "Available targets:"
	@echo "  all                     - Generate Go code for all packages (default target)."
	@echo "  generate                - Generate Go code for all packages."
	@echo "  generate-common         - Generate Go code for common package."
	@echo "  generate-worker         - Generate Go code for worker package."
	@echo "  clean                   - Remove generated files and directories."
	@echo "  help                    - Show this help message."

.PHONY: all generate generate-common generate-worker clean help