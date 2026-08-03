# go-service-template

A [Copier](https://copier.readthedocs.io/) template for a Go backend service, built around one layered architecture that scales from a single CRUD resource up to a Postgres + Redis + Kafka + S3 setup.

Every request — HTTP today, a Kafka consumer or cron job tomorrow — flows through the same layers:

```
handler → dto (validated) → service → repository / client / messaging / cache / storage
```

See `template/README.md` for the full rendered explanation of what each package (`handler`, `service`, `repository`, `client`, `messaging`, `cache`, `locking`, `storage`, `domain`, `apperror`, `observability`, `config`) is responsible for — that file is generated into every project this template produces.

## Usage

```bash
pip install copier   # or: pipx install copier / uv tool install copier
copier copy gh:<you>/go-service-template my-service
cd my-service
go mod tidy
make dev            # docker compose: Postgres + whichever of Redis/Kafka/MinIO you enabled
make migrate-up
make run
```

To pull in template updates later:

```bash
copier update
```

## What you're prompted for

| Prompt | Effect |
| --- | --- |
| `project_name` / `module_name` | Human-readable name and Go module path. |
| `use_redis` | Includes `internal/cache` (Redis) and `internal/locking` (distributed locks), and wires them into the example service. |
| `use_kafka` | Includes `internal/messaging` (producer + consumer) and publishes/consumes a `user.created` event. |
| `use_s3` | Includes `internal/storage`, backed by MinIO in local `docker-compose.yml`. |
| `use_docker` | Dockerfile + docker-compose.yml for local infra. |
| `use_github_actions` | CI workflow: vet, gofmt check, test, (and image build if Docker is on). |

Postgres, the composition root (`cmd/api/dependencies.go`), and one canonical resource (`User`, full CRUD end to end) are always included as the worked example — see `template/README.md`'s "Adding a new resource" section for how to add your own on top of it.

## Design decisions this template encodes

- **No `infrastructure/` catch-all.** Packages are named for the responsibility they hold (`repository`, `client`, `messaging`, `cache`, `storage`, `locking`), not for being "infra."
- **`repository/` and `client/` are first-class**, symmetric with `service/` — interface and implementation live side by side (`user.go` / `user_postgres.go`), not split into an `interfaces/` package.
- **`dto/` is not `domain/`.** Wire contracts and business models are different concerns with different reasons to change; client request/response models belong with the client that owns them, not with `dto/`.
- **`apperror/`, not `errors/`**, to avoid colliding with Go's stdlib package.
- **`observability/`, not `logging/` or `monitoring/`** — logging is one of several signals (metrics, health, tracing later) the app itself emits.
- Deliberately absent: `common/`, `shared/`, `base/`, `helpers/`, `interfaces/` — these tend to become where everything goes.

## Repo layout

```
copier.yml          prompts, conditional file inclusion, post-gen gofmt task
template/            the actual project template (Jinja2-rendered)
```
