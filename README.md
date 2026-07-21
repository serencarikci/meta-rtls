# MetaRTLS

Metadata-driven, multi-tenant RTLS platform. Built with Go, React and Oracle DB.

This app simulates indoor location tracking without real devices. Each customer (tenant) has its own sites, zones, metadata fields and live tags.

## Stack

| Part | Tech |
|------|------|
| Backend | Go, Gin, go-ora |
| Frontend | React, TypeScript, Vite, MUI |
| Database | Oracle 23ai Free |
| Cache | Redis |
| Messaging | MQTT (Mosquitto) |
| Auth | JWT |
| Live updates | WebSocket |
| Local run | Docker Compose, Makefile |
| CI | GitHub Actions |

## What it does

- Tenant login and isolation
- Site → building → floor → zone setup
- Metadata definitions, fields and validation
- MQTT tag simulator and location history in Oracle
- Live map over WebSocket
- Compare small / medium / large customer needs
- Simple change impact score for metadata updates
- Health / ready checks, basic security headers, CI tests

## Quick start

```bash
cp config/config-temp.env config/config.env
make up
make deps
make backend-run
make frontend-run
```

- API health: http://localhost:8090/health
- API ready (Oracle ping): http://localhost:8090/ready
- UI: http://localhost:5173

Oracle may need a few minutes on first start.

```bash
make ready   # when the API is running
make test    # backend unit tests
```

### Demo login

| Tenant | Email | Password |
|--------|-------|----------|
| warehouse-s | admin@warehouse-s.demo | MetaRTLS!2026 |
| hospital-m | admin@hospital-m.demo | MetaRTLS!2026 |
| factory-l | admin@factory-l.demo | MetaRTLS!2026 |

## Production notes

- Edit `config/config.env` (JSON)
- Set `appEnv` to `production`
- Use a long random `jwtSecret` (32+ characters, not the template value)
- Set `corsOrigins` to your real UI URL
- Keep secrets in `config/config.env` (gitignored; use `config-temp.env` as template)

## Folders

- `backend/` — Go API (see `backend/README.md`)
- `frontend/` — React app (see `frontend/README.md`)
- `config/` — JSON app config (`config.env`, `config-temp.env`)
- `deploy/` — Mosquitto config
- `docs/images/` — sample images
- `.github/workflows/` — CI (gofmt, vet, test, build)

## License

Private / educational portfolio project.
