package analysis

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"strings"
)

type Repository struct {
	db *sql.DB
}

func NewRepository(db *sql.DB) *Repository {
	return &Repository{db: db}
}

func (r *Repository) ListRequirements(ctx context.Context) ([]Requirement, error) {
	const q = `
SELECT RAWTOHEX(r.id), RAWTOHEX(r.tenant_id), t.code, t.profile_scale, r.code, r.title,
       NVL(r.description, ''), r.priority, NVL(r.expected_tags, 0), NVL(r.expected_eps, 0),
       NVL(r.retention_days, 0), r.created_at
FROM customer_requirements r
JOIN tenants t ON t.id = r.tenant_id
ORDER BY t.profile_scale, t.code`

	rows, err := r.db.QueryContext(ctx, q)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var out []Requirement
	for rows.Next() {
		var item Requirement
		if err := rows.Scan(
			&item.ID, &item.TenantID, &item.TenantCode, &item.ProfileScale, &item.Code, &item.Title,
			&item.Description, &item.Priority, &item.ExpectedTags, &item.ExpectedEPS,
			&item.RetentionDays, &item.CreatedAt,
		); err != nil {
			return nil, err
		}
		item.ID = strings.ToLower(item.ID)
		item.TenantID = strings.ToLower(item.TenantID)
		out = append(out, item)
	}
	return out, rows.Err()
}

func (r *Repository) ListProfiles(ctx context.Context) ([]TenantProfile, error) {
	const q = `
SELECT RAWTOHEX(t.id), t.code, t.name, t.profile_scale,
       NVL(r.title, ''), NVL(r.expected_tags, 0), NVL(r.expected_eps, 0), NVL(r.retention_days, 0)
FROM tenants t
LEFT JOIN customer_requirements r ON r.tenant_id = t.id
ORDER BY CASE t.profile_scale WHEN 'SMALL' THEN 1 WHEN 'MEDIUM' THEN 2 ELSE 3 END, t.code`

	rows, err := r.db.QueryContext(ctx, q)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var out []TenantProfile
	for rows.Next() {
		var item TenantProfile
		if err := rows.Scan(
			&item.TenantID, &item.Code, &item.Name, &item.ProfileScale,
			&item.RequirementTitle, &item.ExpectedTags, &item.ExpectedEPS, &item.RetentionDays,
		); err != nil {
			return nil, err
		}
		item.TenantID = strings.ToLower(item.TenantID)
		item.FeatureCodes = []string{}
		out = append(out, item)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}

	for i := range out {
		features, err := r.listFeatureCodes(ctx, out[i].TenantID)
		if err != nil {
			return nil, err
		}
		out[i].FeatureCodes = features
		fields, err := r.countMetadataFields(ctx, out[i].TenantID)
		if err != nil {
			return nil, err
		}
		out[i].MetadataFields = fields
	}
	return out, nil
}

func (r *Repository) listFeatureCodes(ctx context.Context, tenantID string) ([]string, error) {
	const q = `
SELECT f.code
FROM tenant_features tf
JOIN feature_definitions f ON f.id = tf.feature_id
WHERE tf.tenant_id = HEXTORAW(:1) AND tf.enabled = 'Y'
ORDER BY f.code`
	rows, err := r.db.QueryContext(ctx, q, strings.ToUpper(tenantID))
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var out []string
	for rows.Next() {
		var code string
		if err := rows.Scan(&code); err != nil {
			return nil, err
		}
		out = append(out, code)
	}
	if out == nil {
		out = []string{}
	}
	return out, rows.Err()
}

func (r *Repository) countMetadataFields(ctx context.Context, tenantID string) (int, error) {
	var count int
	err := r.db.QueryRowContext(ctx, `
SELECT COUNT(*) FROM metadata_fields WHERE tenant_id = HEXTORAW(:1)`,
		strings.ToUpper(tenantID)).Scan(&count)
	return count, err
}

func (r *Repository) CountTenants(ctx context.Context) (int, error) {
	var count int
	err := r.db.QueryRowContext(ctx, `SELECT COUNT(*) FROM tenants WHERE status = 'ACTIVE'`).Scan(&count)
	return count, err
}

func (r *Repository) CountTagsAndAssets(ctx context.Context) (int, error) {
	var tags, assets int
	if err := r.db.QueryRowContext(ctx, `SELECT COUNT(*) FROM tags`).Scan(&tags); err != nil {
		return 0, err
	}
	if err := r.db.QueryRowContext(ctx, `SELECT COUNT(*) FROM assets`).Scan(&assets); err != nil {
		return 0, err
	}
	return tags + assets, nil
}

func (r *Repository) SumExpectedTags(ctx context.Context) (int, error) {
	var sum sql.NullInt64
	err := r.db.QueryRowContext(ctx, `SELECT SUM(expected_tags) FROM customer_requirements`).Scan(&sum)
	if err != nil {
		return 0, err
	}
	if !sum.Valid {
		return 0, nil
	}
	return int(sum.Int64), nil
}

func (r *Repository) SaveChangeRequest(ctx context.Context, tenantID string, result ImpactResult, payload any) (string, error) {
	body, err := json.Marshal(payload)
	if err != nil {
		return "", err
	}
	mig := "N"
	if result.MigrationRequired {
		mig = "Y"
	}
	bc := "N"
	if result.BackwardCompatible {
		bc = "Y"
	}

	_, err = r.db.ExecContext(ctx, `
INSERT INTO change_requests (
  tenant_id, request_type, title, payload_json, affected_tenants, affected_entities,
  migration_required, risk_level, complexity_score, backward_compatible, status
) VALUES (
  HEXTORAW(:1), :2, :3, :4, :5, :6, :7, :8, :9, :10, 'DRAFT'
)`,
		strings.ToUpper(tenantID), result.RequestType, result.Title, string(body),
		result.AffectedTenants, result.AffectedEntities, mig, result.RiskLevel,
		result.ComplexityScore, bc,
	)
	if err != nil {
		return "", err
	}

	var id string
	err = r.db.QueryRowContext(ctx, `
SELECT RAWTOHEX(id) FROM change_requests
WHERE tenant_id = HEXTORAW(:1) AND title = :2
ORDER BY created_at DESC FETCH FIRST 1 ROWS ONLY`,
		strings.ToUpper(tenantID), result.Title).Scan(&id)
	if err != nil {
		return "", err
	}
	return strings.ToLower(id), nil
}

func (r *Repository) ListChangeRequests(ctx context.Context, tenantID string) ([]ChangeRequest, error) {
	const q = `
SELECT RAWTOHEX(id), RAWTOHEX(tenant_id), request_type, title, affected_tenants, affected_entities,
       migration_required, risk_level, complexity_score, backward_compatible, status, created_at
FROM change_requests
WHERE tenant_id = HEXTORAW(:1)
ORDER BY created_at DESC`

	rows, err := r.db.QueryContext(ctx, q, strings.ToUpper(tenantID))
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var out []ChangeRequest
	for rows.Next() {
		var item ChangeRequest
		var mig, bc string
		if err := rows.Scan(
			&item.ID, &item.TenantID, &item.RequestType, &item.Title, &item.AffectedTenants, &item.AffectedEntities,
			&mig, &item.RiskLevel, &item.ComplexityScore, &bc, &item.Status, &item.CreatedAt,
		); err != nil {
			return nil, err
		}
		item.ID = strings.ToLower(item.ID)
		item.TenantID = strings.ToLower(item.TenantID)
		item.MigrationRequired = mig == "Y"
		item.BackwardCompatible = bc == "Y"
		out = append(out, item)
	}
	if out == nil {
		out = []ChangeRequest{}
	}
	return out, rows.Err()
}

func (r *Repository) InsertAudit(ctx context.Context, tenantID, action, details string) error {
	_, err := r.db.ExecContext(ctx, `
INSERT INTO audit_logs (tenant_id, action, entity_type, details_json)
VALUES (HEXTORAW(:1), :2, 'CHANGE_REQUEST', :3)`,
		strings.ToUpper(tenantID), action, details)
	return err
}

func (r *Repository) CountRequirements(ctx context.Context) (int, error) {
	var count int
	err := r.db.QueryRowContext(ctx, `SELECT COUNT(*) FROM customer_requirements`).Scan(&count)
	return count, err
}

func (r *Repository) EnsureDemoRequirements(ctx context.Context) error {
	count, err := r.CountRequirements(ctx)
	if err != nil {
		return err
	}
	if count > 0 {
		return nil
	}

	demos := []struct {
		tenantCode string
		code       string
		title      string
		desc       string
		priority   string
		tags       int
		eps        float64
		retention  int
	}{
		{"warehouse-s", "WH-BASE", "Small warehouse baseline", "Few tags, short retention", "MEDIUM", 50, 5, 30},
		{"hospital-m", "HOSP-SAFE", "Hospital safety tracking", "More tags, room-level zones", "HIGH", 500, 40, 90},
		{"factory-l", "FAC-SCALE", "Large factory scale", "High event rate and long history", "CRITICAL", 5000, 400, 365},
	}

	for _, demo := range demos {
		var tenantID string
		err := r.db.QueryRowContext(ctx, `
SELECT RAWTOHEX(id) FROM tenants WHERE UPPER(code) = UPPER(:1)`, demo.tenantCode).Scan(&tenantID)
		if err != nil {
			return fmt.Errorf("tenant %s: %w", demo.tenantCode, err)
		}
		_, err = r.db.ExecContext(ctx, `
INSERT INTO customer_requirements (
  tenant_id, code, title, description, priority, expected_tags, expected_eps, retention_days
) VALUES (HEXTORAW(:1), :2, :3, :4, :5, :6, :7, :8)`,
			tenantID, demo.code, demo.title, demo.desc, demo.priority, demo.tags, demo.eps, demo.retention)
		if err != nil {
			return err
		}
	}
	return nil
}
