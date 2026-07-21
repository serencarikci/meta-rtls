
CREATE TABLE tenants (
    id              RAW(16) DEFAULT SYS_GUID() PRIMARY KEY,
    code            VARCHAR2(64)  NOT NULL,
    name            VARCHAR2(255) NOT NULL,
    profile_scale   VARCHAR2(16)  DEFAULT 'SMALL' NOT NULL,
    status          VARCHAR2(32)  DEFAULT 'ACTIVE' NOT NULL,
    created_at      TIMESTAMP WITH TIME ZONE DEFAULT SYSTIMESTAMP NOT NULL,
    updated_at      TIMESTAMP WITH TIME ZONE DEFAULT SYSTIMESTAMP NOT NULL,
    CONSTRAINT uq_tenants_code UNIQUE (code),
    CONSTRAINT ck_tenants_scale CHECK (profile_scale IN ('SMALL', 'MEDIUM', 'LARGE')),
    CONSTRAINT ck_tenants_status CHECK (status IN ('ACTIVE', 'SUSPENDED', 'ARCHIVED'))
);

CREATE TABLE users (
    id              RAW(16) DEFAULT SYS_GUID() PRIMARY KEY,
    tenant_id       RAW(16) NOT NULL,
    email           VARCHAR2(320) NOT NULL,
    password_hash   VARCHAR2(255) NOT NULL,
    display_name    VARCHAR2(255) NOT NULL,
    role            VARCHAR2(32)  DEFAULT 'TENANT_ADMIN' NOT NULL,
    status          VARCHAR2(32)  DEFAULT 'ACTIVE' NOT NULL,
    created_at      TIMESTAMP WITH TIME ZONE DEFAULT SYSTIMESTAMP NOT NULL,
    updated_at      TIMESTAMP WITH TIME ZONE DEFAULT SYSTIMESTAMP NOT NULL,
    CONSTRAINT fk_users_tenant FOREIGN KEY (tenant_id) REFERENCES tenants(id),
    CONSTRAINT uq_users_tenant_email UNIQUE (tenant_id, email),
    CONSTRAINT ck_users_role CHECK (role IN ('PLATFORM_ADMIN', 'TENANT_ADMIN', 'OPERATOR', 'VIEWER')),
    CONSTRAINT ck_users_status CHECK (status IN ('ACTIVE', 'DISABLED'))
);

CREATE INDEX ix_users_tenant ON users (tenant_id);

CREATE TABLE sites (
    id              RAW(16) DEFAULT SYS_GUID() PRIMARY KEY,
    tenant_id       RAW(16) NOT NULL,
    code            VARCHAR2(64)  NOT NULL,
    name            VARCHAR2(255) NOT NULL,
    timezone        VARCHAR2(64)  DEFAULT 'UTC' NOT NULL,
    status          VARCHAR2(32)  DEFAULT 'ACTIVE' NOT NULL,
    created_at      TIMESTAMP WITH TIME ZONE DEFAULT SYSTIMESTAMP NOT NULL,
    updated_at      TIMESTAMP WITH TIME ZONE DEFAULT SYSTIMESTAMP NOT NULL,
    CONSTRAINT fk_sites_tenant FOREIGN KEY (tenant_id) REFERENCES tenants(id),
    CONSTRAINT uq_sites_tenant_code UNIQUE (tenant_id, code)
);

CREATE TABLE buildings (
    id              RAW(16) DEFAULT SYS_GUID() PRIMARY KEY,
    tenant_id       RAW(16) NOT NULL,
    site_id         RAW(16) NOT NULL,
    code            VARCHAR2(64)  NOT NULL,
    name            VARCHAR2(255) NOT NULL,
    created_at      TIMESTAMP WITH TIME ZONE DEFAULT SYSTIMESTAMP NOT NULL,
    updated_at      TIMESTAMP WITH TIME ZONE DEFAULT SYSTIMESTAMP NOT NULL,
    CONSTRAINT fk_buildings_tenant FOREIGN KEY (tenant_id) REFERENCES tenants(id),
    CONSTRAINT fk_buildings_site FOREIGN KEY (site_id) REFERENCES sites(id),
    CONSTRAINT uq_buildings_site_code UNIQUE (site_id, code)
);

CREATE TABLE floors (
    id              RAW(16) DEFAULT SYS_GUID() PRIMARY KEY,
    tenant_id       RAW(16) NOT NULL,
    building_id     RAW(16) NOT NULL,
    code            VARCHAR2(64)  NOT NULL,
    name            VARCHAR2(255) NOT NULL,
    level_index     NUMBER(5) DEFAULT 0 NOT NULL,
    width_m         NUMBER(12,3) DEFAULT 100 NOT NULL,
    height_m        NUMBER(12,3) DEFAULT 60 NOT NULL,
    created_at      TIMESTAMP WITH TIME ZONE DEFAULT SYSTIMESTAMP NOT NULL,
    updated_at      TIMESTAMP WITH TIME ZONE DEFAULT SYSTIMESTAMP NOT NULL,
    CONSTRAINT fk_floors_tenant FOREIGN KEY (tenant_id) REFERENCES tenants(id),
    CONSTRAINT fk_floors_building FOREIGN KEY (building_id) REFERENCES buildings(id),
    CONSTRAINT uq_floors_building_code UNIQUE (building_id, code)
);

CREATE TABLE zones (
    id              RAW(16) DEFAULT SYS_GUID() PRIMARY KEY,
    tenant_id       RAW(16) NOT NULL,
    floor_id        RAW(16) NOT NULL,
    code            VARCHAR2(64)  NOT NULL,
    name            VARCHAR2(255) NOT NULL,
    zone_type       VARCHAR2(64)  DEFAULT 'GENERIC' NOT NULL,
    min_x           NUMBER(12,3) NOT NULL,
    min_y           NUMBER(12,3) NOT NULL,
    max_x           NUMBER(12,3) NOT NULL,
    max_y           NUMBER(12,3) NOT NULL,
    created_at      TIMESTAMP WITH TIME ZONE DEFAULT SYSTIMESTAMP NOT NULL,
    updated_at      TIMESTAMP WITH TIME ZONE DEFAULT SYSTIMESTAMP NOT NULL,
    CONSTRAINT fk_zones_tenant FOREIGN KEY (tenant_id) REFERENCES tenants(id),
    CONSTRAINT fk_zones_floor FOREIGN KEY (floor_id) REFERENCES floors(id),
    CONSTRAINT uq_zones_floor_code UNIQUE (floor_id, code),
    CONSTRAINT ck_zones_bbox CHECK (min_x < max_x AND min_y < max_y)
);

CREATE INDEX ix_sites_tenant ON sites (tenant_id);
CREATE INDEX ix_buildings_tenant ON buildings (tenant_id);
CREATE INDEX ix_floors_tenant ON floors (tenant_id);
CREATE INDEX ix_zones_tenant ON zones (tenant_id);

CREATE TABLE device_types (
    id              RAW(16) DEFAULT SYS_GUID() PRIMARY KEY,
    tenant_id       RAW(16),
    code            VARCHAR2(64)  NOT NULL,
    name            VARCHAR2(255) NOT NULL,
    category        VARCHAR2(64)  NOT NULL,
    created_at      TIMESTAMP WITH TIME ZONE DEFAULT SYSTIMESTAMP NOT NULL,
    CONSTRAINT ck_device_types_category CHECK (category IN ('ANCHOR', 'GATEWAY', 'TAG', 'OTHER')),
    CONSTRAINT uq_device_types_scope_code UNIQUE (tenant_id, code)
);

CREATE TABLE anchors (
    id              RAW(16) DEFAULT SYS_GUID() PRIMARY KEY,
    tenant_id       RAW(16) NOT NULL,
    floor_id        RAW(16) NOT NULL,
    device_type_id  RAW(16),
    code            VARCHAR2(64)  NOT NULL,
    name            VARCHAR2(255) NOT NULL,
    x               NUMBER(12,3) NOT NULL,
    y               NUMBER(12,3) NOT NULL,
    status          VARCHAR2(32)  DEFAULT 'ACTIVE' NOT NULL,
    created_at      TIMESTAMP WITH TIME ZONE DEFAULT SYSTIMESTAMP NOT NULL,
    CONSTRAINT fk_anchors_tenant FOREIGN KEY (tenant_id) REFERENCES tenants(id),
    CONSTRAINT fk_anchors_floor FOREIGN KEY (floor_id) REFERENCES floors(id),
    CONSTRAINT uq_anchors_tenant_code UNIQUE (tenant_id, code)
);

CREATE TABLE gateways (
    id              RAW(16) DEFAULT SYS_GUID() PRIMARY KEY,
    tenant_id       RAW(16) NOT NULL,
    site_id         RAW(16) NOT NULL,
    code            VARCHAR2(64)  NOT NULL,
    name            VARCHAR2(255) NOT NULL,
    endpoint_url    VARCHAR2(512),
    status          VARCHAR2(32)  DEFAULT 'ACTIVE' NOT NULL,
    created_at      TIMESTAMP WITH TIME ZONE DEFAULT SYSTIMESTAMP NOT NULL,
    CONSTRAINT fk_gateways_tenant FOREIGN KEY (tenant_id) REFERENCES tenants(id),
    CONSTRAINT fk_gateways_site FOREIGN KEY (site_id) REFERENCES sites(id),
    CONSTRAINT uq_gateways_tenant_code UNIQUE (tenant_id, code)
);

CREATE TABLE tags (
    id              RAW(16) DEFAULT SYS_GUID() PRIMARY KEY,
    tenant_id       RAW(16) NOT NULL,
    code            VARCHAR2(64)  NOT NULL,
    name            VARCHAR2(255) NOT NULL,
    battery_pct     NUMBER(5,2),
    status          VARCHAR2(32)  DEFAULT 'ACTIVE' NOT NULL,
    created_at      TIMESTAMP WITH TIME ZONE DEFAULT SYSTIMESTAMP NOT NULL,
    updated_at      TIMESTAMP WITH TIME ZONE DEFAULT SYSTIMESTAMP NOT NULL,
    CONSTRAINT fk_tags_tenant FOREIGN KEY (tenant_id) REFERENCES tenants(id),
    CONSTRAINT uq_tags_tenant_code UNIQUE (tenant_id, code)
);

CREATE TABLE assets (
    id              RAW(16) DEFAULT SYS_GUID() PRIMARY KEY,
    tenant_id       RAW(16) NOT NULL,
    tag_id          RAW(16),
    code            VARCHAR2(64)  NOT NULL,
    name            VARCHAR2(255) NOT NULL,
    asset_kind      VARCHAR2(64)  DEFAULT 'GENERIC' NOT NULL,
    metadata_json   CLOB CHECK (metadata_json IS JSON),
    schema_version  NUMBER(10) DEFAULT 1 NOT NULL,
    created_at      TIMESTAMP WITH TIME ZONE DEFAULT SYSTIMESTAMP NOT NULL,
    updated_at      TIMESTAMP WITH TIME ZONE DEFAULT SYSTIMESTAMP NOT NULL,
    CONSTRAINT fk_assets_tenant FOREIGN KEY (tenant_id) REFERENCES tenants(id),
    CONSTRAINT fk_assets_tag FOREIGN KEY (tag_id) REFERENCES tags(id),
    CONSTRAINT uq_assets_tenant_code UNIQUE (tenant_id, code)
);

CREATE TABLE persons (
    id              RAW(16) DEFAULT SYS_GUID() PRIMARY KEY,
    tenant_id       RAW(16) NOT NULL,
    tag_id          RAW(16),
    code            VARCHAR2(64)  NOT NULL,
    full_name       VARCHAR2(255) NOT NULL,
    department      VARCHAR2(128),
    metadata_json   CLOB CHECK (metadata_json IS JSON),
    schema_version  NUMBER(10) DEFAULT 1 NOT NULL,
    created_at      TIMESTAMP WITH TIME ZONE DEFAULT SYSTIMESTAMP NOT NULL,
    updated_at      TIMESTAMP WITH TIME ZONE DEFAULT SYSTIMESTAMP NOT NULL,
    CONSTRAINT fk_persons_tenant FOREIGN KEY (tenant_id) REFERENCES tenants(id),
    CONSTRAINT fk_persons_tag FOREIGN KEY (tag_id) REFERENCES tags(id),
    CONSTRAINT uq_persons_tenant_code UNIQUE (tenant_id, code)
);

CREATE TABLE metadata_definitions (
    id              RAW(16) DEFAULT SYS_GUID() PRIMARY KEY,
    tenant_id       RAW(16) NOT NULL,
    entity_type     VARCHAR2(64)  NOT NULL,
    code            VARCHAR2(64)  NOT NULL,
    name            VARCHAR2(255) NOT NULL,
    description     VARCHAR2(1000),
    current_version NUMBER(10) DEFAULT 1 NOT NULL,
    status          VARCHAR2(32)  DEFAULT 'ACTIVE' NOT NULL,
    created_at      TIMESTAMP WITH TIME ZONE DEFAULT SYSTIMESTAMP NOT NULL,
    updated_at      TIMESTAMP WITH TIME ZONE DEFAULT SYSTIMESTAMP NOT NULL,
    CONSTRAINT fk_md_def_tenant FOREIGN KEY (tenant_id) REFERENCES tenants(id),
    CONSTRAINT uq_md_def UNIQUE (tenant_id, entity_type, code),
    CONSTRAINT ck_md_def_entity CHECK (entity_type IN ('ASSET', 'PERSON', 'TAG', 'ZONE', 'DEVICE'))
);

CREATE TABLE schema_versions (
    id              RAW(16) DEFAULT SYS_GUID() PRIMARY KEY,
    tenant_id       RAW(16) NOT NULL,
    definition_id   RAW(16) NOT NULL,
    version_no      NUMBER(10) NOT NULL,
    changelog       VARCHAR2(2000),
    is_current      CHAR(1) DEFAULT 'N' NOT NULL,
    created_at      TIMESTAMP WITH TIME ZONE DEFAULT SYSTIMESTAMP NOT NULL,
    CONSTRAINT fk_sv_tenant FOREIGN KEY (tenant_id) REFERENCES tenants(id),
    CONSTRAINT fk_sv_def FOREIGN KEY (definition_id) REFERENCES metadata_definitions(id),
    CONSTRAINT uq_sv UNIQUE (definition_id, version_no),
    CONSTRAINT ck_sv_current CHECK (is_current IN ('Y', 'N'))
);

CREATE TABLE metadata_fields (
    id              RAW(16) DEFAULT SYS_GUID() PRIMARY KEY,
    tenant_id       RAW(16) NOT NULL,
    definition_id   RAW(16) NOT NULL,
    schema_version_id RAW(16) NOT NULL,
    field_key       VARCHAR2(128) NOT NULL,
    label           VARCHAR2(255) NOT NULL,
    data_type       VARCHAR2(32)  NOT NULL,
    is_required     CHAR(1) DEFAULT 'N' NOT NULL,
    min_value       NUMBER,
    max_value       NUMBER,
    regex_pattern   VARCHAR2(512),
    enum_values     CLOB CHECK (enum_values IS JSON),
    default_value   VARCHAR2(1000),
    sort_order      NUMBER(5) DEFAULT 0 NOT NULL,
    created_at      TIMESTAMP WITH TIME ZONE DEFAULT SYSTIMESTAMP NOT NULL,
    CONSTRAINT fk_mf_tenant FOREIGN KEY (tenant_id) REFERENCES tenants(id),
    CONSTRAINT fk_mf_def FOREIGN KEY (definition_id) REFERENCES metadata_definitions(id),
    CONSTRAINT fk_mf_sv FOREIGN KEY (schema_version_id) REFERENCES schema_versions(id),
    CONSTRAINT uq_mf_key UNIQUE (schema_version_id, field_key),
    CONSTRAINT ck_mf_type CHECK (data_type IN ('STRING', 'NUMBER', 'BOOLEAN', 'ENUM', 'DATE', 'JSON')),
    CONSTRAINT ck_mf_required CHECK (is_required IN ('Y', 'N'))
);

CREATE TABLE metadata_values (
    id              RAW(16) DEFAULT SYS_GUID() PRIMARY KEY,
    tenant_id       RAW(16) NOT NULL,
    entity_type     VARCHAR2(64)  NOT NULL,
    entity_id       RAW(16) NOT NULL,
    definition_id   RAW(16) NOT NULL,
    schema_version_id RAW(16) NOT NULL,
    values_json     CLOB NOT NULL CHECK (values_json IS JSON),
    created_at      TIMESTAMP WITH TIME ZONE DEFAULT SYSTIMESTAMP NOT NULL,
    updated_at      TIMESTAMP WITH TIME ZONE DEFAULT SYSTIMESTAMP NOT NULL,
    CONSTRAINT fk_mv_tenant FOREIGN KEY (tenant_id) REFERENCES tenants(id),
    CONSTRAINT fk_mv_def FOREIGN KEY (definition_id) REFERENCES metadata_definitions(id),
    CONSTRAINT fk_mv_sv FOREIGN KEY (schema_version_id) REFERENCES schema_versions(id),
    CONSTRAINT uq_mv_entity UNIQUE (tenant_id, entity_type, entity_id, definition_id)
);

CREATE TABLE feature_definitions (
    id              RAW(16) DEFAULT SYS_GUID() PRIMARY KEY,
    code            VARCHAR2(64)  NOT NULL,
    name            VARCHAR2(255) NOT NULL,
    description     VARCHAR2(1000),
    category        VARCHAR2(64)  DEFAULT 'CORE' NOT NULL,
    CONSTRAINT uq_feature_code UNIQUE (code)
);

CREATE TABLE tenant_features (
    id              RAW(16) DEFAULT SYS_GUID() PRIMARY KEY,
    tenant_id       RAW(16) NOT NULL,
    feature_id      RAW(16) NOT NULL,
    enabled         CHAR(1) DEFAULT 'Y' NOT NULL,
    config_json     CLOB CHECK (config_json IS JSON),
    created_at      TIMESTAMP WITH TIME ZONE DEFAULT SYSTIMESTAMP NOT NULL,
    CONSTRAINT fk_tf_tenant FOREIGN KEY (tenant_id) REFERENCES tenants(id),
    CONSTRAINT fk_tf_feature FOREIGN KEY (feature_id) REFERENCES feature_definitions(id),
    CONSTRAINT uq_tf UNIQUE (tenant_id, feature_id),
    CONSTRAINT ck_tf_enabled CHECK (enabled IN ('Y', 'N'))
);

CREATE TABLE customer_requirements (
    id              RAW(16) DEFAULT SYS_GUID() PRIMARY KEY,
    tenant_id       RAW(16) NOT NULL,
    code            VARCHAR2(64)  NOT NULL,
    title           VARCHAR2(255) NOT NULL,
    description     VARCHAR2(2000),
    priority        VARCHAR2(16)  DEFAULT 'MEDIUM' NOT NULL,
    expected_tags   NUMBER(10),
    expected_eps    NUMBER(12,2),
    retention_days  NUMBER(10),
    created_at      TIMESTAMP WITH TIME ZONE DEFAULT SYSTIMESTAMP NOT NULL,
    CONSTRAINT fk_cr_tenant FOREIGN KEY (tenant_id) REFERENCES tenants(id),
    CONSTRAINT uq_cr UNIQUE (tenant_id, code),
    CONSTRAINT ck_cr_priority CHECK (priority IN ('LOW', 'MEDIUM', 'HIGH', 'CRITICAL'))
);

CREATE TABLE change_requests (
    id              RAW(16) DEFAULT SYS_GUID() PRIMARY KEY,
    tenant_id       RAW(16) NOT NULL,
    request_type    VARCHAR2(64)  NOT NULL,
    title           VARCHAR2(255) NOT NULL,
    payload_json    CLOB CHECK (payload_json IS JSON),
    affected_tenants NUMBER(10) DEFAULT 0,
    affected_entities NUMBER(10) DEFAULT 0,
    migration_required CHAR(1) DEFAULT 'N' NOT NULL,
    risk_level      VARCHAR2(16)  DEFAULT 'LOW' NOT NULL,
    complexity_score NUMBER(10,2) DEFAULT 0,
    backward_compatible CHAR(1) DEFAULT 'Y' NOT NULL,
    status          VARCHAR2(32)  DEFAULT 'DRAFT' NOT NULL,
    created_at      TIMESTAMP WITH TIME ZONE DEFAULT SYSTIMESTAMP NOT NULL,
    CONSTRAINT fk_chreq_tenant FOREIGN KEY (tenant_id) REFERENCES tenants(id),
    CONSTRAINT ck_chreq_mig CHECK (migration_required IN ('Y', 'N')),
    CONSTRAINT ck_chreq_bc CHECK (backward_compatible IN ('Y', 'N')),
    CONSTRAINT ck_chreq_risk CHECK (risk_level IN ('LOW', 'MEDIUM', 'HIGH', 'CRITICAL'))
);

CREATE TABLE location_events (
    id              RAW(16) DEFAULT SYS_GUID() NOT NULL,
    tenant_id       RAW(16) NOT NULL,
    tag_id          RAW(16) NOT NULL,
    floor_id        RAW(16) NOT NULL,
    x               NUMBER(12,3) NOT NULL,
    y               NUMBER(12,3) NOT NULL,
    quality         NUMBER(5,2),
    event_ts        TIMESTAMP WITH TIME ZONE NOT NULL,
    ingested_at     TIMESTAMP WITH TIME ZONE DEFAULT SYSTIMESTAMP NOT NULL,
    CONSTRAINT pk_location_events PRIMARY KEY (id, event_ts)
)
PARTITION BY RANGE (event_ts)
INTERVAL (NUMTODSINTERVAL(1, 'DAY'))
(
    PARTITION p_loc_bootstrap VALUES LESS THAN (TIMESTAMP '2026-01-01 00:00:00 UTC')
);

CREATE INDEX ix_loc_tenant_ts ON location_events (tenant_id, event_ts);
CREATE INDEX ix_loc_tag_ts ON location_events (tag_id, event_ts);

CREATE TABLE zone_events (
    id              RAW(16) DEFAULT SYS_GUID() PRIMARY KEY,
    tenant_id       RAW(16) NOT NULL,
    tag_id          RAW(16) NOT NULL,
    zone_id         RAW(16) NOT NULL,
    event_type      VARCHAR2(32) NOT NULL,
    event_ts        TIMESTAMP WITH TIME ZONE NOT NULL,
    created_at      TIMESTAMP WITH TIME ZONE DEFAULT SYSTIMESTAMP NOT NULL,
    CONSTRAINT fk_ze_tenant FOREIGN KEY (tenant_id) REFERENCES tenants(id),
    CONSTRAINT fk_ze_tag FOREIGN KEY (tag_id) REFERENCES tags(id),
    CONSTRAINT fk_ze_zone FOREIGN KEY (zone_id) REFERENCES zones(id),
    CONSTRAINT ck_ze_type CHECK (event_type IN ('ZONE_ENTERED', 'ZONE_EXITED'))
);

CREATE INDEX ix_zone_events_tenant_ts ON zone_events (tenant_id, event_ts);

CREATE TABLE audit_logs (
    id              RAW(16) DEFAULT SYS_GUID() PRIMARY KEY,
    tenant_id       RAW(16),
    actor_user_id   RAW(16),
    action          VARCHAR2(64)  NOT NULL,
    entity_type     VARCHAR2(64),
    entity_id       RAW(16),
    details_json    CLOB CHECK (details_json IS JSON),
    created_at      TIMESTAMP WITH TIME ZONE DEFAULT SYSTIMESTAMP NOT NULL
);

CREATE INDEX ix_audit_tenant_ts ON audit_logs (tenant_id, created_at);
