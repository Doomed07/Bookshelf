# Shelfmate

A small social network for readers, in the spirit of Letterboxd but for books. Keep a personal bookshelf, mark books as read, rate them from 1 to 100, write reviews and see what other readers think.

Shelfmate is a **Go REST API** with a **static HTML + vanilla JS frontend** served by the same binary, backed by **PostgreSQL**. The catalog ships with **1,022 books** loaded by a migration.

> The original specification (in Russian) lives in [Техническое задание.md](Техническое%20задание.md).

## Features

- **Accounts and sessions.** Sign up, log in, log out. Server-side sessions in an `HttpOnly` cookie, passwords stored as bcrypt hashes.
- **Catalog.** Browse and search books by title or author, open a book page with its reviews.
- **Personal shelf.** Add a book ("want to read"), mark it as read, give it a rating from 1 to 100, write a review of up to 5,000 characters, or remove it.
- **Reader profiles.** Shelf, recent activity and personal statistics for every user.
- **Statistics.** System-wide, per-user and per-book numbers, top books by reads and by score.
- **Website.** Home, catalog, book, readers, profile, statistics, sign up and log in pages. The browser talks to the same JSON API.
- **Observability.** Prometheus metrics, a ready-made Grafana dashboard, structured logs.

## Tech stack

| Area | Technology |
|---|---|
| Language | Go 1.26 |
| HTTP | `net/http` standard library router (no framework) |
| Database | PostgreSQL 18 via [pgx](https://github.com/jackc/pgx) |
| Migrations | [golang-migrate](https://github.com/golang-migrate/migrate) (schema + seed data) |
| Configuration | Environment variables ([envconfig](https://github.com/kelseyhightower/envconfig)) |
| Validation | [validator](https://github.com/go-playground/validator) |
| Logging | [zap](https://github.com/uber-go/zap) |
| API docs | Swagger via [swag](https://github.com/swaggo/swag) |
| Metrics | Prometheus client, Prometheus 3, Grafana 13 |
| Frontend | Static HTML, CSS and vanilla JS, embedded into the binary with `embed` |
| Packaging | Docker, Docker Compose, Makefile |

## Architecture

Every feature is a vertical slice with three layers. Each layer depends only on the one below it, through an interface:

```
transport (HTTP handlers, DTOs)  ->  service (business rules)  ->  repository (PostgreSQL)
```

```
cmd/bookshelfapp/        entry point: wiring, middleware, servers
internal/
  core/                  shared building blocks
    auth/                session cookie, identity in the request context
    config/              application config
    domain/              entities and validation rules (User, Book, ShelfBook, ...)
    errors/              sentinel errors mapped to HTTP statuses
    logger/              zap logger
    metrics/             Prometheus registry and metrics
    repository/postgres/ pgx connection pool and error mapping
    transport/http/      server, router, middleware, request and response helpers
  features/
    auth/                register, login, logout, me
    users/               profiles
    books/               catalog, reviews
    bookshelf/           a user's shelf and activity
    statistics/          system, user and book statistics
migrations/              SQL schema and seed data
web/                     frontend (HTML, CSS, JS), embedded into the binary
deploy/                  Prometheus config and Grafana provisioning
docs/                    generated Swagger documentation
```

Services are tested against in-memory fake repositories, so the test suite needs no database.

## Getting started

### Prerequisites

- Go 1.26 or newer
- Docker with Docker Compose
- `make`

### Run locally

```bash
# 1. Configure the environment
cp .env.example .env        # then edit .env, see the table below

# 2. Start PostgreSQL in Docker and apply migrations
make env-up
make migrate-up-all

# 3. Expose the database port to your machine (127.0.0.1:5432)
make env-port-forward

# 4. Run the application
make app-run
```

Open <http://localhost:8080>. Swagger UI is at <http://localhost:8080/swagger/>.

`make app-run` runs the application on your machine and connects to PostgreSQL through the forwarded port from step 3.

### Configuration

All settings are environment variables, usually kept in `.env` (git-ignored). `.env.example` lists every variable.

| Variable | Default | Description |
|---|---|---|
| `HTTP_ADDR` | required | Address of the public HTTP server, for example `:8080` |
| `HTTP_SHUTDOWN_TIMEOUT` | `30s` | Graceful shutdown timeout |
| `POSTGRES_USER` / `POSTGRES_PASSWORD` / `POSTGRES_DB` | required | Database credentials |
| `POSTGRES_TIMEOUT` | | Database operation timeout |
| `LOGGER_LVL` | `DEBUG` | Log level |
| `TIME_ZONE` | | Application time zone, for example `Europe/Moscow` |
| `STATS_MIN_RATINGS` | `30` | Minimum number of ratings for a book to appear in the top by score |
| `AUTH_SESSION_TTL` | `720h` | Session lifetime |
| `AUTH_BCRYPT_COST` | `12` | bcrypt cost, from 4 to 31 (use a low value in tests only) |
| `AUTH_COOKIE_SECURE` | `false` | Set the `Secure` flag on the session cookie. **Must be `true` behind HTTPS** |
| `METRICS_ADDR` | `:2112` | Address of the separate metrics server |
| `GRAFANA_ADMIN_PASSWORD` | | Password of the Grafana `admin` user |

### Make targets

| Target | What it does |
|---|---|
| `make env-up` / `env-down` | Start / stop the PostgreSQL container |
| `make env-port-forward` / `env-port-close` | Open / close `127.0.0.1:5432` to the database |
| `make env-cleanup` | Delete the database files (asks for confirmation) |
| `make migrate-up-all` / `migrate-up n=1` | Apply all / the next `n` migrations |
| `make migrate-down n=1` / `migrate-down-all` | Roll back `n` / all migrations |
| `make migrate-create seq=name` | Create a new migration pair |
| `make migrate-force version=N` | Force the migration version after a failed run |
| `make app-run` | Run the application locally |
| `make app-deploy` / `app-undeploy` | Build and start / stop the application in Docker |
| `make swagger-gen` | Regenerate the Swagger documentation |
| `make monitoring-up` / `monitoring-down` | Start / stop Prometheus and Grafana |
| `make logs-cleanup` | Delete log files (asks for confirmation) |
| `make ps` | Show the Compose services |

## API

Base path: `/api/v1`. Interactive documentation: `/swagger/`.

| Method | Path | Access | Description |
|---|---|---|---|
| `POST` | `/auth/register` | public | Create an account and log in |
| `POST` | `/auth/login` | public | Log in |
| `POST` | `/auth/logout` | public | Log out (always clears the cookie, answers `204`) |
| `GET` | `/auth/me` | logged in | The current user |
| `GET` | `/users` | public | List readers |
| `GET` | `/users/{id}` | public | One reader |
| `PATCH` | `/users/{id}` | owner | Edit a profile |
| `DELETE` | `/users/{id}` | owner | Delete an account |
| `GET` | `/books` | public | List the catalog (`title`, `author` search) |
| `GET` | `/books/{id}` | public | One book |
| `GET` | `/books/{id}/reviews` | public | Ratings and reviews of a book |
| `GET` | `/reviews` | public | Latest reviews of all readers |
| `GET` | `/users/{user_id}/bookshelf` | public | A reader's shelf (`read=true\|false` filter) |
| `GET` | `/users/{user_id}/bookshelf/{book_id}` | public | One shelf entry |
| `POST` | `/users/{user_id}/bookshelf` | owner | Add a book to the shelf |
| `PATCH` | `/users/{user_id}/bookshelf/{book_id}` | owner | Mark as read, rate, review |
| `DELETE` | `/users/{user_id}/bookshelf/{book_id}` | owner | Remove from the shelf |
| `GET` | `/users/{user_id}/activity` | public | A reader's recent activity |
| `GET` | `/stats` | public | Statistics (`user_id` or `book_id`, `from`, `to`, `top`) |

"Owner" means the logged-in user must be the user from the path; otherwise the API answers `403`.

Lists accept `limit` (1 to 100, default 20) and `offset` (default 0). Errors share one JSON shape:

```json
{ "message": "failed to get book", "error": "book with id 999: not found" }
```

### Example

State-changing requests must be sent as `application/json`. The session cookie is kept in a cookie jar:

```bash
# Create an account (this also logs you in)
curl -c jar.txt -X POST http://localhost:8080/api/v1/auth/register \
  -H 'Content-Type: application/json' \
  -d '{"username":"reader","email":"reader@example.com","password":"Passw0rd!"}'

# Add book 1 to your shelf (use your own user id from the response above)
curl -b jar.txt -X POST http://localhost:8080/api/v1/users/1/bookshelf \
  -H 'Content-Type: application/json' \
  -d '{"book_id":1}'

# Mark it as read, rate it and review it
curl -b jar.txt -X PATCH http://localhost:8080/api/v1/users/1/bookshelf/1 \
  -H 'Content-Type: application/json' \
  -d '{"read":true,"rating":90,"review":"A masterpiece."}'
```

Passwords are 8 to 72 characters of Latin letters, digits and symbols.

## Security

- **Sessions.** A random 32-byte token is sent in the `shelfmate_session` cookie (`HttpOnly`, `SameSite=Lax`, `Secure` when `AUTH_COOKIE_SECURE=true`). Only the SHA-256 hash of the token is stored in the database, together with its expiry.
- **Passwords.** bcrypt hashes with a configurable cost. Login takes the same time for unknown users as for wrong passwords.
- **CSRF.** For unsafe methods the server checks `Sec-Fetch-Site` and `Origin` against its own host and requires `Content-Type: application/json`.
- **Access control.** Changing a profile or a shelf requires being that user. Emails are never returned by the public API.
- **Limits.** Request bodies are capped at 1 MiB.

> Without HTTPS the session cookie travels in plain text. Do not expose the application to the internet until it sits behind TLS, and set `AUTH_COOKIE_SECURE=true` once it does.

## Monitoring

```bash
make monitoring-up
```

| Service | Address | Notes |
|---|---|---|
| Application metrics | `:2112/metrics` | Separate port, not published outside Docker |
| Prometheus | <http://127.0.0.1:9090> | Scrapes the application every 15 s |
| Grafana | <http://127.0.0.1:3000> | Login `admin`, password from `GRAFANA_ADMIN_PASSWORD`; the **Shelfmate** dashboard is provisioned automatically |

The application exposes:

- `shelfmate_http_requests_total{method,route,status}` and `shelfmate_http_request_duration_seconds{method,route}`. The `route` label is the route **pattern** (`/api/v1/books/{id}`), never the raw URL, so the number of time series stays bounded.
- Business counters: registrations, successful and failed logins, books added to shelves, books marked as read, published reviews, books removed.
- PostgreSQL pool gauges (`shelfmate_db_pool_*`) plus the standard Go runtime and process metrics.

Prometheus and Grafana listen on `127.0.0.1` only. On a server, reach Grafana through an SSH tunnel:

```bash
ssh -L 3001:127.0.0.1:3000 user@your-server    # then open http://localhost:3001
```

## Testing

```bash
go test -race ./...
go vet ./...
```

The tests cover domain rules, middleware, the HTTP layer, authentication, statistics and metrics. Repositories are replaced by in-memory fakes, so no database is required.

## Deployment

`make app-deploy` builds the application image (multi-stage [Dockerfile](cmd/bookshelfapp/dockerfile)) and starts it with PostgreSQL through [docker-compose.yaml](docker-compose.yaml). Apply migrations with `make migrate-up-all`. The frontend is embedded in the binary, so nothing else has to be copied into the image.

## Roadmap

Not implemented yet: HTTPS termination, email confirmation, password change.
