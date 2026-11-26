# Common tasks. Run `make help` for the list.

SHELL := /bin/bash
.DEFAULT_GOAL := help

# Code generators are pinned in tools/go.mod so every machine and CI use the
# same versions without installing anything globally.
GO_TOOL := go tool -modfile=tools/go.mod

# Loads .env into the environment of a recipe line, for commands run on the
# host (go run ...) rather than in docker compose.
WITH_ENV := set -a && source ./.env && set +a &&

.PHONY: help dev down dev-db dev-api dev-web test test-go test-web lint lint-go lint-web \
	gen gen-check migrate-up migrate-down migrate-status migrate-new

help: ## Show this help
	@grep -hE '^[a-zA-Z_-]+:.*## ' $(MAKEFILE_LIST) | \
		awk 'BEGIN {FS = ":.*## "}; {printf "  \033[36m%-15s\033[0m %s\n", $$1, $$2}'

.env:
	cp .env.example .env

web/node_modules: web/package-lock.json
	cd web && npm ci
	@touch web/node_modules

# --- Run -------------------------------------------------------------------

dev: .env ## Build and run the whole stack on http://localhost:8080
	docker compose up --build

down: ## Stop the stack (data is kept in the pgdata volume)
	docker compose down

dev-db: .env ## Start only Postgres and apply migrations (for dev-api)
	docker compose up -d postgres
	docker compose run --rm --build migrate

dev-api: .env ## Run the API on the host (port 8081), needs dev-db
	$(WITH_ENV) go run ./cmd/api

dev-web: web/node_modules ## Run Vite with hot reload, proxying /api to dev-api
	cd web && npm run dev

# --- Quality ---------------------------------------------------------------

test: test-go test-web ## Run all tests (integration tests need Docker)

test-go:
	go test -race ./...

test-web: web/node_modules
	cd web && npm run typecheck

lint: lint-go lint-web ## Run go vet, golangci-lint, ESLint and Prettier

lint-go:
	go vet ./...
	golangci-lint run ./...

lint-web: web/node_modules
	cd web && npm run lint

# --- Code generation -------------------------------------------------------

gen: web/node_modules ## Regenerate sqlc, oapi-codegen and the TypeScript client
	$(GO_TOOL) sqlc generate
	$(GO_TOOL) oapi-codegen -config api/oapi-codegen.yaml api/openapi.yaml
	cd web && npm run gen:api

# Only generated paths are compared, so uncommitted hand-written changes do
# not make the check fail locally.
GENERATED := ':(glob)internal/**/sqlcgen/**' internal/platform/httpapi/httpapi.gen.go web/src/lib/api/schema.d.ts

gen-check: gen ## Fail if generated code is out of date (used by CI)
	git diff --exit-code -- $(GENERATED)
	@if [ -n "$$(git status --porcelain --untracked-files=all -- $(GENERATED))" ]; then \
		echo "generated files are not committed:"; git status --short -- $(GENERATED); exit 1; fi

# --- Database --------------------------------------------------------------

migrate-up: .env ## Apply pending migrations to the dev database
	$(WITH_ENV) go run ./cmd/migrate up

migrate-down: .env ## Roll back the latest migration
	$(WITH_ENV) go run ./cmd/migrate down

migrate-status: .env ## Show applied and pending migrations
	$(WITH_ENV) go run ./cmd/migrate status

migrate-new: ## Create db/migrations/NNNN_<name>.sql, e.g. make migrate-new name=add_products
	@test -n "$(name)" || (echo "usage: make migrate-new name=add_products" && exit 1)
	go run ./cmd/migrate create $(name)
