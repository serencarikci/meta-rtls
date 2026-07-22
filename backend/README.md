# MetaRTLS Backend

This is the Go API for MetaRTLS.
One service. Clear folders. No microservices.

## What you need

\- Docker Desktop (for Oracle and MQTT)
\- Go 1.24+
\- Config file in the repo root: `config/config.env`

## How to start the backend

Do these steps from the **repo root**.

### 1) Config

```bash
cp config/config-temp.env config/config.env
```

### 2) Start Oracle and MQTT

```bash
make up
```

Wait until Oracle is ready (first start can take a few minutes).

### 3) Install Go modules

```bash
cd backend
go mod tidy
```

### 4) Run the API

```bash
go run ./cmd/api
```

Or from the repo root:

```bash
make backend-run
```

### 5) Check that it works

\- http://localhost:8090/health — process is up
\- http://localhost:8090/ready — Oracle is up
\- http://localhost:8090?func=getversion — API version
\- http://localhost:8090?func=getconfig — public config (no secrets)

If `/ready` fails, wait for Oracle and try again.

Current backend version: `0.1.0` (`internal/version`).

## Stack

\- Go + Gin
\- Oracle (`database/sql` + go-ora)
\- JWT auth
\- MQTT (Eclipse Paho)
\- WebSocket (Gorilla)

## Main folders

```text
cmd/api/                 app entry
internal/app/            connect modules
internal/config/         load JSON config
internal/modules/
  identity/              login and users
  tenant/                tenants
  rtlsconfig/            sites, floors, zones
  metadata/              definitions, fields, validate
  location/              MQTT, simulator, live positions
  analysis/              requirements, compare, impact
internal/platform/       db, auth, response, logging
migrations/oracle/       SQL schema and seed
```

## Main API groups

\- `/api/v1/auth/*`
\- `/api/v1/tenants`, `/api/v1/sites`, `/api/v1/floors`, `/api/v1/floors/:id/zones`
\- `/api/v1/metadata/*`
\- `/api/v1/locations/latest`, `/api/v1/ws/locations?token=...`
\- `/api/v1/simulator/*`
\- `/api/v1/analysis/*`

## Tests

```bash
cd backend
go test ./...
go vet ./...
```

## Notes

\- Default port: `8090`
\- Config files: `config/config.env` (local) and `config/config-temp.env` (template)
\- Optional: `CONFIG_PATH=/path/to/file`
\- Logs: `logs/debug.log`, `logs/info.log`, `logs/error.log` (+ console)
\- Optional: `LOG_DIR=/path/to/logs`
\- Demo users load on startup when the DB is ready
\- If MQTT is down, the simulator can still run locally
