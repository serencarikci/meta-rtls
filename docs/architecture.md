# MetaRTLS Architecture

## Decision summary

| Decision | Choice | Why |
|----------|--------|-----|
| Architecture style | Modular monolith | One person can build it; domain borders stay clear |
| Backend | Go + Gin | High concurrent location data, low memory, good for IoT |
| Frontend | React + TypeScript | Tenant/metadata UI + live map |
| Database | Oracle 23ai Free | Partition, MV, scheduler, trigger; fits Oracle experience |
| Messaging | MQTT (Mosquitto) | Path for RTLS device / simulator data |
| Cache | Redis | Live location snapshot / rate limit (later phases) |
| Auth | JWT | Stateless API auth with tenant claim |

The thesis stack was Python/Django + SQL Server. This product version uses **Go + Oracle**. The scientific value is still the metadata-driven multi-tenant model, not the language.

## Module borders

```
internal/modules/
  identity/      JWT login, user bootstrap
  tenant/        Tenant CRUD and profile scale
  rtlsconfig/    Site → Building → Floor → Zone
  metadata/      Definition/field/version (Phase 2)
  location/      Ingestion + query (Phase 3)
  analysis/      Requirement & impact (Phase 4)
```

Business rules stay inside each module. Shared infrastructure is under `internal/platform/` (`db`, `auth`, `response`, `logging`).

## Tenant isolation

- Almost every table has a required `tenant_id`.
- API middleware reads `tenant_id` from the JWT. Repository queries filter by it.
- Later we can add Oracle VPD / Row-Level Security.

## Metadata model

Metadata is not one big free JSON field:

- `metadata_definitions` — entity type + definition
- `schema_versions` — version history
- `metadata_fields` — type, required, min/max, regex, enum
- `metadata_values` — validated JSON values (CLOB IS JSON)

In the first version, values stay in controlled JSON. The system does not create physical SQL tables dynamically.

## Oracle enterprise goals

1. `LOCATION_EVENTS` — daily interval partition (already in schema)
2. `LOCATION_HISTORY` — materialized view (Phase 3+)
3. `DEVICE_STATUS` — DBMS_SCHEDULER job (Phase 3+)
4. `AUDIT_LOG` — write with triggers (Phase 4+)

## Phase plan

1. Core platform (current skeleton)
2. Metadata engine
3. RTLS MQTT + simulator + WebSocket live map
4. Requirement compare + impact analysis
5. Tests, CI/CD, observability

## Out of scope (on purpose)

- Kubernetes / microservices
- Real UWB engine
- Machine learning
- Full event sourcing
- Dynamic SQL table generation
