# {{ project_name }}

{% if author_name %}Maintained by {{ author_name }}.{% endif %}

A Go backend service scaffolded from [go-service-template](https://github.com/{{ author_name | default('your-org', true) }}/go-service-template), using a layered architecture that stays consistent whether you're reading this in Go, Spring Boot, or FastAPI.

## Getting started

```bash
go mod tidy
{% if use_docker %}
make dev            # start Postgres{% if use_redis %} + Redis{% endif %}{% if use_kafka %} + Kafka{% endif %}{% if use_s3 %} + MinIO{% endif %}
{% endif %}
make migrate-up      # apply database migrations
make run             # start the API on :{{ http_port }}
```

Health check: `curl localhost:{{ http_port }}/health`

## Architecture

Every request flows through the same layers, regardless of the entry point (HTTP today; a Kafka consumer or cron job tomorrow would look identical):

```
Entry point (handler/)
      │
Request DTO (dto/request) — validated
      │
Service (service/)  ── business logic & orchestration
      │
      ├── repository/   → domain persistence (Postgres via pgx)
      ├── client/       → outbound calls to other services
{% if use_kafka %}      ├── messaging/    → publish/consume events (Kafka)
{% endif %}{% if use_redis %}      ├── cache/        → read-through cache (Redis)
      ├── locking/      → distributed locks (Redis)
{% endif %}{% if use_s3 %}      └── storage/      → object storage (S3 / MinIO)
{% endif %}
      │
Response DTO (dto/response)
```

| Package | Responsibility |
| --- | --- |
| `handler/` | HTTP transport: parses requests, calls services, writes responses. Knows nothing about business rules. |
| `middleware/` | Cross-cutting HTTP concerns: request ID, recovery, logging, CORS. |
| `dto/` | Wire-format request/response models plus their validation. Never reused as domain models. |
| `service/` | Business logic and orchestration. Depends only on interfaces (`repository`, `client`{% if use_kafka %}, `messaging`{% endif %}{% if use_redis %}, `cache`{% endif %}), never on concrete infrastructure. |
| `repository/` | Persistence abstraction. Interface + Postgres implementation live side by side (`user.go` / `user_postgres.go`). |
| `client/` | Abstractions for calling other services/APIs. Each external service owns its own request/response models — they never leak into `dto/`. |
{% if use_kafka %}| `messaging/` | Kafka producers/consumers. Entry points here are as thin as HTTP handlers — they deserialize, validate, and call a service. |
{% endif %}{% if use_redis %}| `cache/` | Cache abstraction (Redis). Not a source of truth — always backed by a repository. |
| `locking/` | Distributed locking (Redis) for coordinating across instances. |
{% endif %}{% if use_s3 %}| `storage/` | Object storage abstraction (S3-compatible; MinIO for local dev). |
{% endif %}| `domain/` | Core business models and invariants. No framework or infrastructure imports. |
| `apperror/` | Typed application errors (not Go's stdlib `errors`, to avoid import collisions). Handlers translate these into HTTP status codes. |
| `observability/` | Structured logging, `/health` and `/ready`, Prometheus metrics. |
| `config/` | Environment-based configuration loading. |

Deliberately **not** present: a catch-all `infrastructure/`, `common/`, `util`-as-junk-drawer, or per-service interfaces created "just in case." Every package answers *what responsibility does this code have*, not *what layer is this technically in*.

## Composition root

`cmd/api/main.go` stays tiny — it just calls `BuildApplication()` and starts the server. All wiring (config → connections → repositories → services → handlers) happens in `cmd/api/dependencies.go`, so business code never knows how it was constructed.

## Adding a new resource

Follow the `user` resource as a template end to end:

1. `domain/user.go` — the business model.
2. `dto/request` / `dto/response` — the wire contracts.
3. `repository/user.go` (interface) + `repository/user_postgres.go` (implementation) — plus a migration in `db/migrations/`.
4. `service/user.go` — orchestration and business rules.
5. `handler/user.go` — HTTP glue, wired into the router in `dependencies.go`.

## Project layout

```
cmd/api/                 entry point + composition root
internal/
  handler/                HTTP handlers
  middleware/              HTTP middleware
  dto/                    request/response models + validation
  service/                business logic
  repository/             persistence (Postgres)
  client/                 outbound API clients
{% if use_kafka %}  messaging/               Kafka producer/consumer
{% endif %}{% if use_redis %}  cache/                   Redis cache
  locking/                 Redis distributed locks
{% endif %}{% if use_s3 %}  storage/                 S3/MinIO object storage
{% endif %}  domain/                  business models
  apperror/                typed application errors
  observability/           logging, metrics, health
  config/                  configuration loading
  util/                    small shared helpers (pagination, HTTP responses)
db/migrations/            SQL migrations
{% if use_docker %}docker-compose.yml       local infra
Dockerfile
{% endif %}Makefile
```
