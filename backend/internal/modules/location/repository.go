package location

import (
	"context"
	"database/sql"
	"fmt"
	"strings"
	"time"
)

type Repository struct {
	db *sql.DB
}

func NewRepository(db *sql.DB) *Repository {
	return &Repository{db: db}
}

func (r *Repository) InsertEvent(ctx context.Context, event PositionEvent) error {
	const q = `
INSERT INTO location_events (tenant_id, tag_id, floor_id, x, y, quality, event_ts)
VALUES (HEXTORAW(:1), HEXTORAW(:2), HEXTORAW(:3), :4, :5, 1, :6)`

	ts := event.Timestamp
	if ts.IsZero() {
		ts = time.Now().UTC()
	}
	_, err := r.db.ExecContext(ctx, q,
		strings.ToUpper(event.TenantID),
		strings.ToUpper(event.TagID),
		strings.ToUpper(event.FloorID),
		event.X,
		event.Y,
		ts,
	)
	return err
}

func (r *Repository) InsertZoneEvent(ctx context.Context, tenantID, tagID, zoneID, eventType string, at time.Time) error {
	const q = `
INSERT INTO zone_events (tenant_id, tag_id, zone_id, event_type, event_ts)
VALUES (HEXTORAW(:1), HEXTORAW(:2), HEXTORAW(:3), :4, :5)`
	_, err := r.db.ExecContext(ctx, q,
		strings.ToUpper(tenantID),
		strings.ToUpper(tagID),
		strings.ToUpper(zoneID),
		eventType,
		at,
	)
	return err
}

func (r *Repository) ListZonesForFloor(ctx context.Context, tenantID, floorID string) ([]ZoneBox, error) {
	const q = `
SELECT RAWTOHEX(id), code, min_x, min_y, max_x, max_y
FROM zones
WHERE tenant_id = HEXTORAW(:1) AND floor_id = HEXTORAW(:2)`

	rows, err := r.db.QueryContext(ctx, q, strings.ToUpper(tenantID), strings.ToUpper(floorID))
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var out []ZoneBox
	for rows.Next() {
		var z ZoneBox
		if err := rows.Scan(&z.ID, &z.Code, &z.MinX, &z.MinY, &z.MaxX, &z.MaxY); err != nil {
			return nil, err
		}
		z.ID = strings.ToLower(z.ID)
		out = append(out, z)
	}
	return out, rows.Err()
}

func (r *Repository) EnsureDemoTags(ctx context.Context, tenantCode string) (tenantID, floorID string, tags []SimTag, err error) {
	err = r.db.QueryRowContext(ctx, `
SELECT RAWTOHEX(id) FROM tenants WHERE UPPER(code) = UPPER(:1)`, tenantCode).Scan(&tenantID)
	if err != nil {
		return "", "", nil, fmt.Errorf("tenant: %w", err)
	}
	tenantID = strings.ToLower(tenantID)

	err = r.db.QueryRowContext(ctx, `
SELECT RAWTOHEX(id) FROM floors
WHERE tenant_id = HEXTORAW(:1)
FETCH FIRST 1 ROWS ONLY`, strings.ToUpper(tenantID)).Scan(&floorID)
	if err != nil {
		return "", "", nil, fmt.Errorf("floor: %w", err)
	}
	floorID = strings.ToLower(floorID)

	demoCodes := []string{"TAG-01", "TAG-02"}
	for i, code := range demoCodes {
		var tagID string
		err = r.db.QueryRowContext(ctx, `
SELECT RAWTOHEX(id) FROM tags
WHERE tenant_id = HEXTORAW(:1) AND UPPER(code) = UPPER(:2)`,
			strings.ToUpper(tenantID), code).Scan(&tagID)
		if err == sql.ErrNoRows {
			_, err = r.db.ExecContext(ctx, `
INSERT INTO tags (tenant_id, code, name, battery_pct, status)
VALUES (HEXTORAW(:1), :2, :3, 90, 'ACTIVE')`,
				strings.ToUpper(tenantID), code, "Demo tag "+code)
			if err != nil {
				return "", "", nil, err
			}
			err = r.db.QueryRowContext(ctx, `
SELECT RAWTOHEX(id) FROM tags
WHERE tenant_id = HEXTORAW(:1) AND UPPER(code) = UPPER(:2)`,
				strings.ToUpper(tenantID), code).Scan(&tagID)
		}
		if err != nil {
			return "", "", nil, err
		}
		tags = append(tags, SimTag{
			TagID:   strings.ToLower(tagID),
			TagCode: code,
			FloorID: floorID,
			X:       10 + float64(i*15),
			Y:       8 + float64(i*5),
			DirX:    0.8,
			DirY:    0.5,
		})
	}
	return tenantID, floorID, tags, nil
}
