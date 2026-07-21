# MetaRTLS Backend

Go API for MetaRTLS. It is a modular monolith (one service, clear folders).

## Stack

- Go + Gin
- Oracle (`database/sql` + go-ora)
- JWT auth
- MQTT (Eclipse Paho)
- WebSocket (Gorilla)

## Run

From the repo root:

```bash
cp config/config-temp.env config/config.env
make up
cd backend
go mod tidy
go run ./cmd/api
```

API base: `http://localhost:8090`

- Health (process up): `GET /health`
- Ready (Oracle ping): `GET /ready`

## Main folders

```text
cmd/api/                 app entry
internal/app/            wire modules together
internal/config/         env config + production checks
internal/modules/
  identity/              login and users
  tenant/                tenants
  rtlsconfig/            sites, floors, zones
  metadata/              definitions, fields, validate
  location/              MQTT, simulator, live positions
  analysis/              requirements, compare, impact
internal/platform/       db, auth, response helpers
migrations/oracle/       SQL schema and seed
```

## Useful API groups

- `/api/v1/auth/*`
- `/api/v1/tenants`, `/api/v1/sites`, `/api/v1/floors`, `/api/v1/floors/:id/zones`
- `/api/v1/metadata/*`
- `/api/v1/locations/latest`, `/api/v1/ws/locations?token=...`
- `/api/v1/simulator/*`
- `/api/v1/analysis/*`

## Quality checks

```bash
go test ./...
go vet ./...
gofmt -l .
```

CI runs the same checks on push and pull requests.

## Security basics

- JWT required for protected routes
- Production rejects short or default `jwtSecret`
- CORS limited to `corsOrigins` in JSON config
- Basic security headers on every response
- Request logs include method, path, status, latency

## Notes

- Config is JSON: `config/config.env` (local) and `config/config-temp.env` (template).
- Optional override: `CONFIG_PATH=/path/to/file`
- Default API port is `8090`.
- Demo users and sample metadata load on startup when the DB is ready.
- If MQTT is down, the simulator still updates positions locally.
