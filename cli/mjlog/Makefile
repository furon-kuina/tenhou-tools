.DEFAULT_GOAL := help

up: ## Do docker compose up
	docker compose up -d

down: ## Do docker compose down
	docker compose down

help: ## Show options
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | \
		awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-20s\033[0m %s\n", $$1, $$2}'

migrate:
	PGPASSWORD=postgres psqldef -h 127.0.0.1 -p 35432 -U postgres tenhou < schema.sql