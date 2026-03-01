.PHONY: build vet lint test test-sb test-selector test-integration test-all report report-junit clean help

GOTESTSUM := $(shell command -v gotestsum 2>/dev/null || echo "$(shell go env GOPATH)/bin/gotestsum")
JUNIT_DIR := junit

help: ## Show this help
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | sort | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-20s\033[0m %s\n", $$1, $$2}'

build: ## Build all packages
	go build ./...

vet: ## Run go vet
	go vet ./...

lint: build vet ## Build and vet

test: test-selector test-sb ## Run all unit tests

test-selector: ## Run selector unit tests
	go test ./selector/... -v

test-sb: ## Run sb package unit tests
	go test ./sb/... -v

test-integration: ## Run integration tests (requires Playwright browsers)
	go test ./examples/... -tags integration -v -timeout 120s

test-all: test test-integration ## Run all tests (unit + integration)

report-junit: ## Generate JUnit XML reports
	@mkdir -p $(JUNIT_DIR)
	$(GOTESTSUM) --junitfile $(JUNIT_DIR)/unit-selector.xml -- ./selector/... -v
	$(GOTESTSUM) --junitfile $(JUNIT_DIR)/unit-sb.xml -- ./sb/... -v

report: report-junit ## Generate HTML report from test results
	@mkdir -p $(JUNIT_DIR)
	$(GOTESTSUM) --junitfile $(JUNIT_DIR)/integration.xml -- ./examples/... -tags integration -v -timeout 120s || true
	go run ./cmd/report $(JUNIT_DIR)/*.xml
	@echo "Report generated: report.html"
	@which open > /dev/null 2>&1 && open report.html || true

install-tools: ## Install required tools
	go install gotest.tools/gotestsum@latest

clean: ## Clean generated files
	rm -rf $(JUNIT_DIR) report.html visual_baseline coupang_*.png

.DEFAULT_GOAL := help
