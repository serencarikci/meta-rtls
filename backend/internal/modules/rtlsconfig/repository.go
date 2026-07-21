package rtlsconfig

import (
	"context"
	"database/sql"
	"fmt"
	"strings"
)

type Repository struct {
	db *sql.DB
}

func NewRepository(db *sql.DB) *Repository {
	return &Repository{db: db}
}

func (r *Repository) ListSites(ctx context.Context, tenantID string) ([]Site, error) {
	const q = `
SELECT RAWTOHEX(id), RAWTOHEX(tenant_id), code, name, timezone, status, created_at
FROM sites WHERE tenant_id = HEXTORAW(:1) ORDER BY name`

	rows, err := r.db.QueryContext(ctx, q, strings.ToUpper(tenantID))
	if err != nil {
		return nil, fmt.Errorf("list sites: %w", err)
	}
	defer rows.Close()

	var out []Site
	for rows.Next() {
		var s Site
		if err := rows.Scan(&s.ID, &s.TenantID, &s.Code, &s.Name, &s.Timezone, &s.Status, &s.CreatedAt); err != nil {
			return nil, err
		}
		s.ID = strings.ToLower(s.ID)
		s.TenantID = strings.ToLower(s.TenantID)
		out = append(out, s)
	}
	return out, rows.Err()
}

func (r *Repository) CreateSite(ctx context.Context, tenantID string, req CreateSiteRequest) (*Site, error) {
	tz := req.Timezone
	if tz == "" {
		tz = "UTC"
	}
	const insertQ = `
INSERT INTO sites (tenant_id, code, name, timezone, status)
VALUES (HEXTORAW(:1), :2, :3, :4, 'ACTIVE')`
	if _, err := r.db.ExecContext(ctx, insertQ, strings.ToUpper(tenantID), req.Code, req.Name, tz); err != nil {
		return nil, fmt.Errorf("create site: %w", err)
	}

	const selectQ = `
SELECT RAWTOHEX(id), RAWTOHEX(tenant_id), code, name, timezone, status, created_at
FROM sites WHERE tenant_id = HEXTORAW(:1) AND UPPER(code) = UPPER(:2)`
	var s Site
	if err := r.db.QueryRowContext(ctx, selectQ, strings.ToUpper(tenantID), req.Code).Scan(
		&s.ID, &s.TenantID, &s.Code, &s.Name, &s.Timezone, &s.Status, &s.CreatedAt,
	); err != nil {
		return nil, err
	}
	s.ID = strings.ToLower(s.ID)
	s.TenantID = strings.ToLower(s.TenantID)
	return &s, nil
}

func (r *Repository) ListFloors(ctx context.Context, tenantID string) ([]Floor, error) {
	const q = `
SELECT RAWTOHEX(id), RAWTOHEX(tenant_id), RAWTOHEX(building_id), code, name, level_index, width_m, height_m
FROM floors WHERE tenant_id = HEXTORAW(:1) ORDER BY level_index`

	rows, err := r.db.QueryContext(ctx, q, strings.ToUpper(tenantID))
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var out []Floor
	for rows.Next() {
		var f Floor
		if err := rows.Scan(&f.ID, &f.TenantID, &f.BuildingID, &f.Code, &f.Name, &f.LevelIndex, &f.WidthM, &f.HeightM); err != nil {
			return nil, err
		}
		f.ID = strings.ToLower(f.ID)
		f.TenantID = strings.ToLower(f.TenantID)
		f.BuildingID = strings.ToLower(f.BuildingID)
		out = append(out, f)
	}
	return out, rows.Err()
}

func (r *Repository) ListZones(ctx context.Context, tenantID, floorID string) ([]Zone, error) {
	const q = `
SELECT RAWTOHEX(id), RAWTOHEX(tenant_id), RAWTOHEX(floor_id), code, name, zone_type, min_x, min_y, max_x, max_y
FROM zones
WHERE tenant_id = HEXTORAW(:1) AND floor_id = HEXTORAW(:2)
ORDER BY name`

	rows, err := r.db.QueryContext(ctx, q, strings.ToUpper(tenantID), strings.ToUpper(floorID))
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var out []Zone
	for rows.Next() {
		var z Zone
		if err := rows.Scan(&z.ID, &z.TenantID, &z.FloorID, &z.Code, &z.Name, &z.ZoneType, &z.MinX, &z.MinY, &z.MaxX, &z.MaxY); err != nil {
			return nil, err
		}
		z.ID = strings.ToLower(z.ID)
		z.TenantID = strings.ToLower(z.TenantID)
		z.FloorID = strings.ToLower(z.FloorID)
		out = append(out, z)
	}
	return out, rows.Err()
}

func (r *Repository) ListBuildings(ctx context.Context, tenantID string) ([]Building, error) {
	const q = `
SELECT RAWTOHEX(id), RAWTOHEX(tenant_id), RAWTOHEX(site_id), code, name
FROM buildings WHERE tenant_id = HEXTORAW(:1) ORDER BY name`

	rows, err := r.db.QueryContext(ctx, q, strings.ToUpper(tenantID))
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var out []Building
	for rows.Next() {
		var b Building
		if err := rows.Scan(&b.ID, &b.TenantID, &b.SiteID, &b.Code, &b.Name); err != nil {
			return nil, err
		}
		b.ID = strings.ToLower(b.ID)
		b.TenantID = strings.ToLower(b.TenantID)
		b.SiteID = strings.ToLower(b.SiteID)
		out = append(out, b)
	}
	return out, rows.Err()
}
