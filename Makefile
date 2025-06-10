# VibeRS Makefile (Optional - prefer direct go commands)

.PHONY: help build test lint clean db-init run dev

help: ## Show this help message
	@echo "VibeRS - Parallel Recall → Dedup → Three-Stage Ranking"
	@echo ""
	@echo "Available commands:"
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | sort | awk 'BEGIN {FS = ":.*?## "}; {printf "  \033[36m%-15s\033[0m %s\n", $$1, $$2}'

build: ## Build the API server
	go build -o bin/api ./cmd/api

test: ## Run all tests
	go test ./...

lint: ## Run code linting
	./scripts/lint.sh

clean: ## Clean build artifacts
	rm -rf bin/
	rm -f data/vibers.db

db-init: ## Initialize database
	./scripts/db_init.sh

run: build ## Build and run API server
	./bin/api

dev: ## Run API server in development mode
	go run ./cmd/api

# Python model training
model-env: ## Set up Python environment for model training
	cd model-training && python -m venv .venv && source .venv/bin/activate && pip install -r requirements.txt

# Default target
.DEFAULT_GOAL := help 