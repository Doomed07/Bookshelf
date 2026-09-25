include .env
export 

export PROJECT_ROOT=$(shell pwd)
n ?= 1

env-up:
	@docker compose up -d bookshelfapp-postgres

env-down:
	@docker compose down bookshelfapp-postgres

env-cleanup:
	@read -p "Do you really want to clear volume? Risk of data loss. [y/N]: " ans; \
	if [ "$$ans" = "y" ]; then \
		docker compose down bookshelfapp-postgres port-forwarder && \
		rm -rf ${PROJECT_ROOT}/out/pgdata && \
		echo "Environment files have been removed"; \
	else \
		echo "Environment cleanup cancelled"; \
	fi

env-port-forward:
	@docker compose up -d port-forwarder

env-port-close:
	@docker compose down port-forwarder

migrate-create:
	@if [ -z "$(seq)" ]; then \
		echo "Missing parameter: seq. Example: make migrate-create seq=value"; \
		exit 1; \
	fi; \
	docker compose run --rm bookshelfapp-postgres-migrate \
		create \
		-ext sql \
		-dir /migrations \
		-seq "$(seq)"

migrate-up:
	@make migrate-action action="up ${n}"

migrate-up-all:
	@make migrate-action action="up"

migrate-down:
	@make migrate-action action="down ${n}"

migrate-down-all:
	@read -p "Rollback ALL migrations? Schema and data will be dropped. [y/N]: " ans; \
	if [ "$$ans" = "y" ]; then \
		make migrate-action action="down -all"; \
	else \
		echo "Rollback cancelled"; \
	fi

migrate-action:
	@if [ -z "$(action)" ]; then \
		echo "Missing parameter: action. Example: make migrate-action action=value"; \
		exit 1; \
	fi; \
	docker compose run --rm bookshelfapp-postgres-migrate \
		-path /migrations \
		-database postgres://${POSTGRES_USER}:${POSTGRES_PASSWORD}@bookshelfapp-postgres:5432/${POSTGRES_DB}?sslmode=disable \
		${action}

migrate-force:
	@if [ -z "$(version)" ]; then \
		echo "Missing parameter: version. Example: make migrate-force version=number"; \
		exit 1; \
	fi; \
	docker compose run --rm bookshelfapp-postgres-migrate \
	-path /migrations \
	-database postgres://${POSTGRES_USER}:${POSTGRES_PASSWORD}@bookshelfapp-postgres:5432/${POSTGRES_DB}?sslmode=disable \
	force "${version}"

app-run:
	@export LOGGER_FOLDER=${PROJECT_ROOT}/out/logs && \
	export POSTGRES_HOST=localhost && \
	go mod tidy && \
	go run ${PROJECT_ROOT}/cmd/bookshelfapp/main.go