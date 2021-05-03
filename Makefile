CMD=n26

VENDOR_DIR = vendor
BUILD_DIR = build

GO ?= go
GOLANGCI_LINT ?= golangci-lint

.PHONY: $(VENDOR_DIR) lint test test-unit build build-linux

$(VENDOR_DIR):
	@mkdir -p $(VENDOR_DIR)
	@$(GO) mod vendor
	@$(GO) mod tidy

lint:
	@$(GOLANGCI_LINT) run

test: test-unit

## Run unit tests
test-unit:
	@echo ">> unit test"
	@$(GO) test -gcflags=-l -coverprofile=unit.coverprofile -covermode=atomic -race ./...

#test-integration:
#	@echo ">> integration test"
#	@$(GO) test ./features/... -gcflags=-l -coverprofile=features.coverprofile -coverpkg ./... -godog -race

build-linux:
	@echo ">> building binary, GOFLAGS: $(GOFLAGS)"
	@GOOS=linux GOARCH=amd64 $(GO) build -ldflags "$(shell ./resources/scripts/build_args.sh)" -o $(BUILD_DIR)/$(CMD)-linux cmd/$(CMD)/*

## Build binary
build:
	@echo ">> building binary, GOFLAGS: $(GOFLAGS)"
	@rm -f $(BUILD_DIR)/*
	@$(GO) build -ldflags "$(shell ./resources/scripts/build_args.sh)" -o $(BUILD_DIR)/$(CMD) cmd/$(CMD)/*
