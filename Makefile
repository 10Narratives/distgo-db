who ?= worker
config ?= examples/worker_example.yaml

run:
	@go run cmd/$(who)/main.go --config config/$(config)

help:
	@echo "Usage: make [target]"
	@echo ""
	@echo "Targets:"
	@echo "  run    - Runs the application with parameters who and config"
	@echo "  help   - Displays this help message"
	@echo ""
	@echo "Variables:"
	@echo "  who    - Specifies the role or component (default: worker)"
	@echo "  config - Specifies the configuration file (default: examples/worker_example.yaml)"