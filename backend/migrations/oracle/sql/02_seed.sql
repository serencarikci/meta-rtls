INSERT INTO tenants (id, code, name, profile_scale, status)
VALUES (HEXTORAW('11111111111111111111111111111101'), 'warehouse-s', 'Demo Warehouse (Small)', 'SMALL', 'ACTIVE');

INSERT INTO tenants (id, code, name, profile_scale, status)
VALUES (HEXTORAW('11111111111111111111111111111102'), 'hospital-m', 'Demo Hospital (Medium)', 'MEDIUM', 'ACTIVE');

INSERT INTO tenants (id, code, name, profile_scale, status)
VALUES (HEXTORAW('11111111111111111111111111111103'), 'factory-l', 'Demo Factory (Large)', 'LARGE', 'ACTIVE');

INSERT INTO feature_definitions (id, code, name, description, category) VALUES
    (HEXTORAW('33333333333333333333333333333301'), 'LIVE_MAP', 'Live Map', 'WebSocket live asset tracking', 'CORE');
INSERT INTO feature_definitions (id, code, name, description, category) VALUES
    (HEXTORAW('33333333333333333333333333333302'), 'ZONE_EVENTS', 'Zone Events', 'Enter/exit zone detection', 'CORE');
INSERT INTO feature_definitions (id, code, name, description, category) VALUES
    (HEXTORAW('33333333333333333333333333333303'), 'METADATA_ENGINE', 'Metadata Engine', 'Dynamic field definitions and validation', 'CORE');
INSERT INTO feature_definitions (id, code, name, description, category) VALUES
    (HEXTORAW('33333333333333333333333333333304'), 'IMPACT_ANALYSIS', 'Change Impact Analysis', 'Estimate metadata change cost/risk', 'ANALYTICS');
INSERT INTO feature_definitions (id, code, name, description, category) VALUES
    (HEXTORAW('33333333333333333333333333333305'), 'SIMULATOR', 'Location Simulator', 'MQTT-based virtual tag simulator', 'RTLS');

INSERT INTO tenant_features (tenant_id, feature_id, enabled)
SELECT t.id, f.id, 'Y'
FROM tenants t
CROSS JOIN feature_definitions f
WHERE f.code IN ('LIVE_MAP', 'ZONE_EVENTS', 'METADATA_ENGINE', 'SIMULATOR');

INSERT INTO tenant_features (tenant_id, feature_id, enabled)
SELECT t.id, f.id, 'Y'
FROM tenants t
CROSS JOIN feature_definitions f
WHERE t.profile_scale = 'LARGE' AND f.code = 'IMPACT_ANALYSIS';

INSERT INTO sites (id, tenant_id, code, name, timezone)
VALUES (
    HEXTORAW('44444444444444444444444444444401'),
    HEXTORAW('11111111111111111111111111111101'),
    'MAIN', 'Main Campus', 'Europe/Istanbul'
);

INSERT INTO buildings (id, tenant_id, site_id, code, name)
VALUES (
    HEXTORAW('55555555555555555555555555555501'),
    HEXTORAW('11111111111111111111111111111101'),
    HEXTORAW('44444444444444444444444444444401'),
    'WH-A', 'Warehouse A'
);

INSERT INTO floors (id, tenant_id, building_id, code, name, level_index, width_m, height_m)
VALUES (
    HEXTORAW('66666666666666666666666666666601'),
    HEXTORAW('11111111111111111111111111111101'),
    HEXTORAW('55555555555555555555555555555501'),
    'L0', 'Ground Floor', 0, 80, 40
);

INSERT INTO zones (id, tenant_id, floor_id, code, name, zone_type, min_x, min_y, max_x, max_y)
VALUES (
    HEXTORAW('77777777777777777777777777777701'),
    HEXTORAW('11111111111111111111111111111101'),
    HEXTORAW('66666666666666666666666666666601'),
    'RECV', 'Receiving', 'PROCESS', 0, 0, 20, 15
);

INSERT INTO zones (id, tenant_id, floor_id, code, name, zone_type, min_x, min_y, max_x, max_y)
VALUES (
    HEXTORAW('77777777777777777777777777777702'),
    HEXTORAW('11111111111111111111111111111101'),
    HEXTORAW('66666666666666666666666666666601'),
    'STOR', 'Storage', 'STORAGE', 20, 0, 60, 40
);

COMMIT;
