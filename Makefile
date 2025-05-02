SOURCE_FILES ?=./...
TEST_OPTIONS := -v -failfast -race
TEST_PATTERN ?=.
BENCH_OPTIONS ?= -v -bench=. -benchmem
CLEAN_OPTIONS ?=-modcache -testcache
TEST_TIMEOUT ?=1m

.PHONY: all
all: help

.PHONY: help
help:
	@echo "make format - use gofmt, goimports, golines"
	@echo "make lint - run golangci-lint"
	@echo "make test - run go test including race detection"
	@echo "make bench - run go test including benchmarking"


.PHONY: build
build:
	$(info: Make: Build)
	go build -gcflags='-m -m' examples/compiler/*.go

.PHONY: fmt
fmt:
	$(info: Make: Format)
	gofmt -w ./**/*
	goimports -w ./**/*
	golines -w ./**/*

.PHONY: lint
lint:
	$(info: Make: Lint)
	@golangci-lint run --tests=false


.PHONY: test
test:
	CGO_ENABLED=1 go test ${TEST_OPTIONS} ${SOURCE_FILES} -run ${TEST_PATTERN} -timeout=${TEST_TIMEOUT}

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