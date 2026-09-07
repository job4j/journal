GOOSE_VERSION := v3.27.3
TOOLS_DIR := $(CURDIR)/bin
GOOSE := $(TOOLS_DIR)/goose

ifeq ($(OS),Windows_NT)
GOOSE := $(TOOLS_DIR)/goose.exe
endif

MIGRATIONS_DIR := $(CURDIR)/migrations
GOOSE_DRIVER ?= postgres
GOOSE_DBSTRING ?=

.PHONY: deps migration-create migrate-up migrate-down migrate-status migrate-validate

deps:
	mkdir -p "$(TOOLS_DIR)"
	GOBIN="$(TOOLS_DIR)" go install github.com/pressly/goose/v3/cmd/goose@$(GOOSE_VERSION)

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
