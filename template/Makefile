BINARY       := {{ project_slug }}
MODULE       := {{ module_name }}
DATABASE_URL ?= postgres://{{ project_slug | replace('-', '_') }}:{{ project_slug | replace('-', '_') }}@localhost:5432/{{ project_slug | replace('-', '_') }}?sslmode=disable

.PHONY: run build test lint tidy migrate-up migrate-down{% if use_docker %} dev dev-down{% endif %}

run: ## Run the API locally
	go run ./cmd/api

build: ## Build the API binary
	go build -o bin/$(BINARY) ./cmd/api

test: ## Run tests
	go test ./... -race -cover

lint: ## Run go vet + gofmt check
	go vet ./...
	@test -z "$$(gofmt -l .)" || (echo "gofmt needs to be run on the files above" && exit 1)

tidy: ## Tidy go.mod
	go mod tidy

migrate-up: ## Apply database migrations
	@for f in db/migrations/*.up.sql; do \
		echo "applying $$f"; \
		psql "$(DATABASE_URL)" -f $$f || exit 1; \
	done

migrate-down: ## Roll back the last database migration
	@f=$$(ls -r db/migrations/*.down.sql | head -n1); \
	echo "reverting $$f"; \
	psql "$(DATABASE_URL)" -f $$f
{% if use_docker %}

dev: ## Start local infrastructure (Postgres{% if use_redis %}, Redis{% endif %}{% if use_kafka %}, Kafka{% endif %}{% if use_s3 %}, MinIO{% endif %})
	docker compose up -d

dev-down: ## Stop local infrastructure
	docker compose down -v
{% endif %}
