# MetaRTLS

Metadata-driven, multi-tenant RTLS data architecture and location tracking platform built with Go, React and Oracle DB.

This project is an **RTLS platform simulator** and **metadata management system**. It works without real devices. It turns a thesis idea about *metadata-oriented data architecture* into a modular monolith you can run and show.

## Tech stack

| Layer | Technology |
|--------|-----------|
| Backend | Go, Gin, database/sql + go-ora |
| Frontend | React, TypeScript, Vite, Material UI |
| Database | Oracle Database 23ai Free |
| Cache | Redis |
| Messaging | MQTT (Eclipse Mosquitto) |
| Auth | JWT |
| Realtime | WebSocket |
| API docs | Swagger / OpenAPI |
| DevOps | Docker Compose, Makefile |

## Architecture overview

The first version is a **modular monolith**. There are no microservices.

```
React Web Application
        |
        | REST API / WebSocket
        v
Go Modular Monolith
        |
        +-- Identity & Access
        +-- Tenant Management
        +-- Metadata Management
        +-- RTLS Configuration
        +-- Location Ingestion
        +-- Location Query
        +-- Rule & Zone Management
        +-- Schema Generation
        +-- Audit Management
        |
        +-- Oracle 23ai Free
        +-- Redis
        +-- MQTT
```

## Link to the thesis topic

The main academic value is not the programming language. It is a **metadata-driven multi-tenant data architecture**:

- Asset and field definitions per tenant (customer)
- Metadata validation and schema versioning
- Compare needs of small / medium / large customers
- Impact and complexity analysis for metadata changes
- RTLS location simulation and live tracking

## Quick start

```bash
cp .env.example .env
make up          # Oracle + Redis + Mosquitto
make deps
make backend-run # :8080
make frontend-run # :5173
```

Oracle may need a few minutes on first start. The init script `backend/migrations/oracle/00_run_as_app_user.sh` creates the schema for the `metartls` user.

### Demo login

| Tenant | Email | Password |
|--------|-------|----------|
| `warehouse-s` | `admin@warehouse-s.demo` | `MetaRTLS!2026` |
| `hospital-m` | `admin@hospital-m.demo` | `MetaRTLS!2026` |
| `factory-l` | `admin@factory-l.demo` | `MetaRTLS!2026` |

## Delivery phases

1. **Core platform** — Auth, tenant, site/building/floor/zone
2. **Metadata engine** — Definition, validation, schema versioning
3. **RTLS flow** — MQTT simulator, ingestion, WebSocket, live map
4. **Architecture analysis** — Requirement matrix, impact analysis
5. **Production quality** — Tests, CI/CD, monitoring, docs

## Oracle enterprise features (planned)

- `LOCATION_EVENTS` daily partition
- `LOCATION_HISTORY` materialized view
- `DEVICE_STATUS` scheduler job
- `AUDIT_LOG` trigger-based audit

## License

Private / educational portfolio project.
