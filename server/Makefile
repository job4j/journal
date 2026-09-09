GOOSE_VERSION := v3.27.3
OAPI_CODEGEN_VERSION := v2.8.0
TOOLS_DIR := $(CURDIR)/bin
GOOSE := $(TOOLS_DIR)/goose
OAPI_CODEGEN := $(TOOLS_DIR)/oapi-codegen

ifeq ($(OS),Windows_NT)
GOOSE := $(TOOLS_DIR)/goose.exe
OAPI_CODEGEN := $(TOOLS_DIR)/oapi-codegen.exe
endif

MIGRATIONS_DIR := $(CURDIR)/migrations
OPENAPI_SPEC := $(CURDIR)/api/openapi.yml
OPENAPI_CONFIG := $(CURDIR)/api/oapi-codegen.yml
OPENAPI_OUTPUT := $(CURDIR)/gen/openapi.gen.go
GOOSE_DRIVER ?= postgres
GOOSE_DBSTRING ?= postgres://postgres:password@127.0.0.1:5433/journal?sslmode=disable

.PHONY: deps generate migration-create migrate-up migrate-down migrate-status migrate-validate test integration-test build lint

deps:
	mkdir -p "$(TOOLS_DIR)"
	GOBIN="$(TOOLS_DIR)" go install github.com/pressly/goose/v3/cmd/goose@$(GOOSE_VERSION)
	GOBIN="$(TOOLS_DIR)" go install github.com/oapi-codegen/oapi-codegen/v2/cmd/oapi-codegen@$(OAPI_CODEGEN_VERSION)

generate:
	mkdir -p "$(dir $(OPENAPI_OUTPUT))"
	"$(OAPI_CODEGEN)" --config "$(OPENAPI_CONFIG)" "$(OPENAPI_SPEC)"
	gofmt -w "$(OPENAPI_OUTPUT)"

migration-create:
	@test -n "$(name)" || (echo "Usage: make migration-create name=migration_name" && exit 1)
	@mkdir -p "$(MIGRATIONS_DIR)"
	"$(GOOSE)" -dir "$(MIGRATIONS_DIR)" create "$(name)" sql

migrate-up:
	@test -n "$(GOOSE_DBSTRING)" || (echo "GOOSE_DBSTRING is required" && exit 1)
	"$(GOOSE)" -dir "$(MIGRATIONS_DIR)" "$(GOOSE_DRIVER)" "$(GOOSE_DBSTRING)" up

migrate-down:
	@test -n "$(GOOSE_DBSTRING)" || (echo "GOOSE_DBSTRING is required" && exit 1)
	"$(GOOSE)" -dir "$(MIGRATIONS_DIR)" "$(GOOSE_DRIVER)" "$(GOOSE_DBSTRING)" down

migrate-status:
	@test -n "$(GOOSE_DBSTRING)" || (echo "GOOSE_DBSTRING is required" && exit 1)
	"$(GOOSE)" -dir "$(MIGRATIONS_DIR)" "$(GOOSE_DRIVER)" "$(GOOSE_DBSTRING)" status

migrate-validate:
	"$(GOOSE)" -dir "$(MIGRATIONS_DIR)" validate

test:
	go test ./...

integration-test: migrate-up
	INTEGRATION_DATABASE_URL="$(GOOSE_DBSTRING)" go test -tags=integration ./internal/integration -count=1

build:
	go build ./...

lint:
	go vet ./...

run: generate
	go run ./cmd
