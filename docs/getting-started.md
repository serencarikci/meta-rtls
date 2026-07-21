# Getting Started

## Requirements

- Docker Desktop
- Go 1.24+
- Node.js 22+

## Setup

```bash
cp .env.example .env
make up
```

Oracle may need 1–3 minutes on first start. The schema is created under `APP_USER` (`metartls`).

```bash
make deps
make backend-run   # http://localhost:8090/health
make frontend-run  # http://localhost:5173
```

## Demo login

| Tenant | Email | Password |
|--------|-------|----------|
| warehouse-s | admin@warehouse-s.demo | MetaRTLS!2026 |
| hospital-m | admin@hospital-m.demo | MetaRTLS!2026 |
| factory-l | admin@factory-l.demo | MetaRTLS!2026 |

On first start, the API creates these users with bcrypt.

## API surface (Phase 1)

- `POST /api/v1/auth/login`
- `GET /api/v1/auth/me`
- `GET/POST /api/v1/tenants`
- `GET/POST /api/v1/sites`
- `GET /api/v1/buildings`
- `GET /api/v1/floors`
- `GET /api/v1/floors/:floorId/zones`
- `GET/POST /api/v1/metadata/definitions`
- `GET /api/v1/metadata/definitions/:id`
- `GET/POST /api/v1/metadata/definitions/:id/versions`
- `GET/POST /api/v1/metadata/versions/:versionId/fields`
- `POST /api/v1/metadata/validate`
- `GET /api/v1/metadata/features`

## Next step

Phase 3: MQTT simulator + location ingestion + WebSocket live map.
