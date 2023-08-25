.DEFAULT_GOAL := all
PROJECT_ROOT ?= .
PROJECT_NAME ?= $(notdir $(CURDIR))
GOOS=linux
GOARCH=amd64

GIT_BRANCH ?= $(shell git rev-parse --abbrev-ref HEAD)
ifeq (release,$(firstword $(MAKECMDGOALS)))
  RUN_ARGS := $(wordlist 2,$(words $(MAKECMDGOALS)),$(MAKECMDGOALS))
  $(eval $(RUN_ARGS):;@:)
endif

gci:
	go install github.com/daixiang0/gci@latest

gofumpt:
	go install mvdan.cc/gofumpt@latest

fmt: gci gofumpt ## format code and imports
	gci write -s standard -s default -s "prefix(github.com)" . --skip-generated
	gofumpt -e -w .

golangci-lint:
	go install github.com/golangci/golangci-lint/cmd/golangci-lint@v1.52.2

lint: golangci-lint ## run linters
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 golangci-lint run -E revive --timeout 5m


help: ## displays this help page
	@awk 'BEGIN {FS = ":.*?## "} /^[a-zA-Z_-]+:.*?## / {sub("\\\\n",sprintf("\n%22c"," "), $$2);printf "\033[36m%-20s\033[0m %s\n", $$1, $$2}' $(MAKEFILE_LIST)

installmigrate: ## install sql-migrate
	go install github.com/rubenv/sql-migrate/...@latest

installutils: gci gofumpt installmigrate ## install utilities
	go install github.com/dmarkham/enumer@latest

migrate: installutils  ## apply migrations
	sql-migrate up -config=dbconfig.yml -env="development"

all: installutils binary ## (default) make binaries + utils

binary: ## build binary
	go build -o skel cmd/*.go

test: ## run tests
	go test -race ./...

run: fmt lint test ## lint+run
	go run cmd/*.go

