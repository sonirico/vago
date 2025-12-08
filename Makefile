# Dynamically get workspace modules
WORKSPACES := $(shell go list -f '{{.Dir}}' -m 2>/dev/null)
SOURCE_FILES := $(shell go list -f '{{.Dir}}/...' -m 2>/dev/null)

TEST_OPTIONS := -json -v -failfast -race
TEST_PATTERN ?=.
BENCH_OPTIONS ?= -v -bench=. -benchmem
CLEAN_OPTIONS ?=-modcache -testcache
TEST_TIMEOUT ?=1m

.PHONY: all
all: help

.PHONY: help
help:
	@echo "make fmt - use gofmt, goimports, golines"
	@echo "make lint - run golangci-lint"
	@echo "make test - run go test including race detection"
	@echo "make bench - run go test including benchmarking"
	@echo "make tidy - run go mod tidy on all workspace modules"
	@echo "make update-deps - update all dependencies in all workspace modules"


.PHONY: fmt
fmt:
	$(info: Make: Format)
	gofmt -w ./**/*.go
	goimports -w ./**/*.go
	golines -w ./**/*.go

.PHONY: lint
lint:
	$(info: Make: Lint)
	@golangci-lint run --tests=false


.PHONY: test
test:
	CGO_ENABLED=1 go test ${TEST_OPTIONS} ${SOURCE_FILES} -run ${TEST_PATTERN} -timeout=${TEST_TIMEOUT} | tparse --all --progress

.PHONY: bench
bench:
	CGO_ENABLED=1 go test ${BENCH_OPTIONS} ${SOURCE_FILES} -run ${TEST_PATTERN} -timeout=${TEST_TIMEOUT}

.PHONY: docs
docs:
	go run readme.go 

.PHONY: setup
setup:
	go install golang.org/x/tools/cmd/goimports@latest
	go install github.com/segmentio/golines@latest
	go install github.com/mfridman/tparse@latest

.PHONY: tidy
tidy:
	$(info: Make: Tidy all workspace modules)
	@for dir in $(WORKSPACES); do \
		echo "Running go mod tidy in $$dir..."; \
		(cd $$dir && go mod tidy) || exit 1; \
	done

.PHONY: update-deps
update-deps:
	$(info: Make: Update dependencies in all workspace modules)
	@for dir in $(WORKSPACES); do \
		echo "Updating dependencies in $$dir..."; \
		(cd $$dir && go get -u ./... && go mod tidy) || exit 1; \
	done
	@echo "All dependencies updated successfully!"