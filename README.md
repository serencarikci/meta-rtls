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

## What it does

- Tenant login and isolation
- Site → building → floor → zone setup
- Metadata definitions, fields and validation
- MQTT tag simulator and location history in Oracle
- Live map over WebSocket
- Compare small / medium / large customer needs
- Simple change impact score for metadata updates

## Quick start

```bash
cp .env.example .env
make up
make deps
make backend-run
make frontend-run
```

- API: http://localhost:8090/health
- UI: http://localhost:5173

Oracle may need a few minutes on first start.

### Demo login

| Tenant | Email | Password |
|--------|-------|----------|
| warehouse-s | admin@warehouse-s.demo | MetaRTLS!2026 |
| hospital-m | admin@hospital-m.demo | MetaRTLS!2026 |
| factory-l | admin@factory-l.demo | MetaRTLS!2026 |

## Folders

- `backend/` — Go API (see `backend/README.md`)
- `frontend/` — React app (see `frontend/README.md`)
- `deploy/` — Mosquitto config
- `docs/images/` — sample images

## License

Private / educational portfolio project.
