# Bookmark Management

A small URL-shortening HTTP service written in Go, built with Gin and Redis. It exposes a
health-check endpoint, a cryptographically secure password generator, and a short-link
endpoint — with Swagger docs, unit tests, and integration tests.

## Features

- **Shorten URL** — `POST /v1/links/shorten` generates a 7-character random code and stores
  the original URL in Redis with a TTL
- **Health check** — `GET /health-check` pings Redis and reports status, service name and instance id
- **Password generator** — `GET /genpass` returns a 12-character random password from
  `crypto/rand`
- **Swagger UI** — served at `/swagger/index.html`
- **Layered architecture** — handler → service → repository, wired with interfaces and
  generated mocks so each layer can be tested in isolation

## Tech stack

| Concern | Choice |
|---|---|
| Language | Go 1.26 |
| HTTP | [Gin](https://github.com/gin-gonic/gin) |
| Storage | Redis ([go-redis/v9](https://github.com/redis/go-redis/v9)) |
| Config | [envconfig](https://github.com/kelseyhightower/envconfig) |
| Docs | [swaggo/swag](https://github.com/swaggo/swag) |
| Tests | `testing`, [testify](https://github.com/stretchr/testify), [miniredis](https://github.com/alicebob/miniredis) |
| Mocks | [mockery](https://github.com/vektra/mockery) v2 |

## Project structure

```
.
├── cmd/
│   ├── api/main.go            # entrypoint: load config, build engine, start server
│   └── test/main.go           # scratch script for manual checks against real Redis
├── internal/
│   ├── api/                   # HTTP engine + route wiring
│   ├── config/                # environment configuration
│   ├── handler/               # Gin handlers (bind input, call service, write response)
│   │   └── mocks/             # (generated) service mocks used by handler tests
│   ├── service/               # business logic
│   │   └── mocks/             # (generated) service mocks used outside this package
│   ├── repository/            # data access (Redis)
│   │   └── mocks/             # (generated) repository mocks
│   └── integration_test/      # end-to-end tests: real engine + in-memory Redis
├── pkg/redis/                 # Redis client factory + config
├── docs/                      # generated Swagger spec
└── Makefile
```

Request flow for `POST /v1/links/shorten`:

```
client → handler.ShortenLink   (bind JSON, map errors → HTTP status)
       → service.ShortenUrl    (generate code, ensure uniqueness)
       → repository.URLStorage (Redis)
```

## Requirements

- Go 1.26 or newer
- Redis (only needed to run the server; tests use an in-memory Redis)

## Quick start

```bash
# 1. install dependencies
go mod download

# 2. make sure Redis is reachable
redis-server                       # or use REDIS_ADDRESS to point at an existing one

# 3. run the API
make run                           # == go run cmd/api/main.go

# 4. try it out
curl http://localhost:8080/health-check
```

Then open the Swagger UI at <http://localhost:8080/swagger/index.html>.

## Configuration

All settings come from environment variables.

| Variable | Default | Used by | Description |
|---|---|---|---|
| `APP_PORT` | `8080` | `internal/config` | Port the HTTP server listens on |
| `SERVICE_NAME` | `bookmark-management` | `internal/config` | Name reported by `/health-check` |
| `INSTANCE_ID` | *(empty)* | `internal/config` | Instance identifier reported by `/health-check` |
| `REDIS_ADDRESS` | `localhost:6379` | `pkg/redis` | Redis host:port |
| `REDIS_PASSWORD` | *(empty)* | `pkg/redis` | Redis password |
| `REDIS_DB` | `0` | `pkg/redis` | Redis database index |

```bash
export APP_PORT=8080
export REDIS_ADDRESS=localhost:6379
export REDIS_DB=0
make run
```

## API endpoints

### `GET /health-check`

```bash
curl http://localhost:8080/health-check
```

```json
{ "message": "OK", "service_name": "bookmark-management", "instance_id": "..." }
```

Returns `500` with `{"error":"Internal Server Error"}` when Redis is unreachable.

### `GET /genpass`

```bash
curl http://localhost:8080/genpass
```

```json
{ "password": "k3Yq8Zt1aBc9" }
```

Always 12 characters, from `abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789`.

### `POST /v1/links/shorten`

| Field | Type | Required | Description |
|---|---|---|---|
| `url` | string | yes | Original URL to shorten |
| `exp` | int | yes | Time to live **in seconds** (`0` means no expiry) |

```bash
curl -X POST http://localhost:8080/v1/links/shorten \
  -H 'Content-Type: application/json' \
  -d '{"url":"https://google.com","exp":3600}'
```

```json
{ "code": "aG7xK2p", "message": "Shorten URL generated successfully!" }
```

- `400` — body is not valid JSON: `{"error":"Invalid input"}`
- `500` — storage failure: `{"error":"Internal Server Error"}`

### `GET /swagger/*any`

Interactive API documentation generated from the handler annotations.

## Testing

```bash
make test          # unit + integration tests, coverage report, enforces 80% threshold
go test ./...      # tests only
go test -race ./...
go test -shuffle=on ./...   # catch tests that depend on execution order
```

`make test` writes `coverage.html` (open it in a browser) and fails the build when total
coverage drops below `COVERAGE_THRESHOLD` (80%, currently ~97%).

Test conventions used throughout the repo:

| Convention | Why |
|---|---|
| Black-box tests live in `package X_test` | they exercise the public API only, so refactoring internals cannot silently break them |
| `internal/app/handler` keeps `package handler` | its tests need the unexported `passwordLength` constant |
| Table-driven tests with `name / setUpX / wantX` | one place to add a scenario, subtests run in parallel |
| `require` stops the test, `assert` keeps going | use `require` for preconditions, `assert` for the checks that belong to the case |
| `t.Context()` everywhere | context is created and cancelled by the test framework |
| `verifyFunc` for side-effect assertions | e.g. "storage lookup failed ⇒ StoreURL must never be called" |
| `mock.Anything` + `.Once()` for repeated calls | documents the exact number of interactions with the dependency |
| Errors are compared with `require.ErrorIs` | works with sentinel errors and wrapped errors alike |

What the suite covers:

| Layer | Approach |
|---|---|
| `internal/app/service`, `internal/app/handler` | table-driven unit tests with mockery mocks; every branch (success, dependency error, invalid input) has a case |
| `internal/app/repository` | real Redis behaviour through `miniredis` |
| `internal/api` | route registration through `ServeHTTP` |
| `internal/config`, `pkg/redis` | environment parsing, using `miniredis` as the server |
| `internal/integration_test` | real engine + in-memory Redis, driving HTTP endpoints end to end, including what ends up in Redis and with which TTL |

Two functions are deliberately not covered: `api.Start` (it blocks forever and the engine has
no `Shutdown` method) and the `rand.Int` error branch in `GeneratePassword`.

### Regenerating mocks

Mocks are generated by [mockery](https://github.com/vektra/mockery) and live next to the
interface they mock. After changing an interface, regenerate it:

```bash
# interface carries the directive, so from the package directory:
go generate ./internal/service/...
go generate ./internal/repository/...
```

or explicitly:

```bash
mockery --name=ShortenUrl --filename=shortenurl.go --dir=internal/service --output=internal/service/mocks
```

## Design notes

A few decisions worth knowing before changing the code:

- **Dependencies are injected as interfaces.** `NewShortenUrl(repo repository.URLStorage,
  codeGen GenPass)` receives everything it needs, so tests can pass mocks and no globals are
  used. `URLStorage` is declared in `repository`, while `GenPass`, `ShortenUrl` and
  `HealthCheck` are declared in `service` next to their consumers.
- **Each layer only knows the layer below.** The service must not import
  `go-redis`. Storage errors are translated at the repository boundary into
  `repository.ErrNotFound` (a sentinel error), so the service reads
  `errors.Is(err, repository.ErrNotFound)` and stays storage-agnostic.
- **Context reaches the storage layer.** Handlers own the context they pass down — today
  they create `context.Background()`, which is one of the limitations listed below; the
  target is `c.Request.Context()` so cancellations and deadlines propagate.
- **Tests are black-box by default.** Every package is tested from `package X_test` so that
  only the exported API is relied upon; `internal/app/handler` is the exception because its tests
  need the unexported `passwordLength`. This also avoids an import cycle: the generated
  `service/mocks` package imports `service`.

## Known limitations

- No graceful shutdown: `engine.Start()` blocks and the process is killed on deploy, so
  in-flight requests can be cut.
- Short-code generation is a non-atomic check-then-set (`GetURL` then `StoreURL`): two
  concurrent requests can be served the same code. `SET key value NX EX ttl` would make it
  atomic.
- Code generation retries in an unbounded `for` loop — a busy keyspace can spin forever
  instead of returning an error.
- No structured logging or request ids; every failure is reported as a generic 500.
- No CI pipeline yet — `make test` is only enforced locally.
- `cmd/test/main.go` is a scratch script and is not part of the build.

## Contributing

1. Branch off `main`.
2. Keep layers separated — do not import Redis (or any other client) into `internal/app/service`.
3. Add or update table-driven tests for every behaviour change.
4. Run `make test` before opening a PR.

## License

Not specified — add a LICENSE file before publishing this code.
