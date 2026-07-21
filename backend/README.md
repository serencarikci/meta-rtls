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
cp .env.example .env
make up
cd backend
go mod tidy
go run ./cmd/api
```

API base: `http://localhost:8090`

Health check: `GET /health`

## Main folders

```text
cmd/api/                 app entry
internal/app/            wire modules together
internal/config/         env config
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

## Tests and format

```bash
go test ./...
go fmt ./...
```

## Notes

- Default API port is `8090` (see `.env` / `.env.example`).
- Demo users and sample metadata load on startup when the DB is ready.
- If MQTT is down, the simulator still updates positions locally.
