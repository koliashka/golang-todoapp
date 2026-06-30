include .env
export


export PROJECT_ROOT=$(CURDIR)

env-up:
	@docker compose up -d todoapp-postgres

env-down:
	@docker compose down todoapp-postgres

env-cleanup:
	@read -p "eboshim? [y/N]: " ans; \
	if [ "$$ans" = "y" ]; then \
		docker compose down && \
		rm -rf out/pgdata && \
		echo "Vse pohereli"; \
	else \
		echo "fuf"; \
	fi

env-port-forward:
	@docker compose up -d port-forwarder
	
env-port-close:
	@docker compose down port-forwarder

migrate-create:
	$(if $(strip $(seq)),,$(error ne peredano nazvanie migracii. Primer: make migrate-create seq=create_todos))
	env MSYS_NO_PATHCONV=1 docker compose run --rm todoapp-postgres-migrate create \
		-ext sql \
		-dir /migrations \
		-seq "$(seq)"

migrate-up:
	make migrate-action action=up

migrate-down:
	make migrate-action action=down

migrate-action:
	$(if $(strip $(action)),,$(error ne peredano peremenaya action))
	env MSYS_NO_PATHCONV=1 docker compose run --rm todoapp-postgres-migrate \
		-path /migrations \
		-database postgres://${POSTGRES_USER}:${POSTGRES_PASSWORD}@todoapp-postgres:5432/${POSTGRES_DB}?sslmode=disable \
		"${action}"

todoapp-run:
	export LOGGER_FOLDER=${PROJECT_ROOT}/out/log && \
	export POSTGRES_HOST=localhost && \
	go mod tidy && \
	go run cmd/todoapp/main.go
