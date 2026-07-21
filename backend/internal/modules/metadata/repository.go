package metadata

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

func (r *Repository) ListDefinitions(ctx context.Context, tenantID string) ([]Definition, error) {
	const q = `
SELECT RAWTOHEX(id), RAWTOHEX(tenant_id), entity_type, code, name,
       NVL(description, ''), current_version, status, created_at
FROM metadata_definitions
WHERE tenant_id = HEXTORAW(:1)
ORDER BY name`

	rows, err := r.db.QueryContext(ctx, q, strings.ToUpper(tenantID))
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var out []Definition
	for rows.Next() {
		var item Definition
		if err := rows.Scan(
			&item.ID, &item.TenantID, &item.EntityType, &item.Code, &item.Name,
			&item.Description, &item.CurrentVersion, &item.Status, &item.CreatedAt,
		); err != nil {
			return nil, err
		}
		item.ID = strings.ToLower(item.ID)
		item.TenantID = strings.ToLower(item.TenantID)
		out = append(out, item)
	}
	return out, rows.Err()
}

func (r *Repository) GetDefinition(ctx context.Context, tenantID, id string) (*Definition, error) {
	const q = `
SELECT RAWTOHEX(id), RAWTOHEX(tenant_id), entity_type, code, name,
       NVL(description, ''), current_version, status, created_at
FROM metadata_definitions
WHERE tenant_id = HEXTORAW(:1) AND id = HEXTORAW(:2)`

	var item Definition
	err := r.db.QueryRowContext(ctx, q, strings.ToUpper(tenantID), strings.ToUpper(id)).Scan(
		&item.ID, &item.TenantID, &item.EntityType, &item.Code, &item.Name,
		&item.Description, &item.CurrentVersion, &item.Status, &item.CreatedAt,
	)
	if err == sql.ErrNoRows {
		return nil, ERR_NOT_FOUND
	}
	if err != nil {
		return nil, err
	}
	item.ID = strings.ToLower(item.ID)
	item.TenantID = strings.ToLower(item.TenantID)
	return &item, nil
}

func (r *Repository) CreateDefinition(ctx context.Context, tenantID string, req CreateDefinitionRequest) (*Definition, error) {
	const insertDef = `
INSERT INTO metadata_definitions (tenant_id, entity_type, code, name, description, current_version, status)
VALUES (HEXTORAW(:1), :2, :3, :4, :5, 1, 'ACTIVE')`
	if _, err := r.db.ExecContext(ctx, insertDef, strings.ToUpper(tenantID), req.EntityType, req.Code, req.Name, req.Description); err != nil {
		return nil, fmt.Errorf("insert definition: %w", err)
	}

	def, err := r.findDefinitionByCode(ctx, tenantID, req.EntityType, req.Code)
	if err != nil {
		return nil, err
	}

	const insertVer = `
INSERT INTO schema_versions (tenant_id, definition_id, version_no, changelog, is_current)
VALUES (HEXTORAW(:1), HEXTORAW(:2), 1, 'initial version', 'Y')`
	if _, err := r.db.ExecContext(ctx, insertVer, strings.ToUpper(tenantID), strings.ToUpper(def.ID)); err != nil {
		return nil, fmt.Errorf("insert version: %w", err)
	}
	return def, nil
}

func (r *Repository) findDefinitionByCode(ctx context.Context, tenantID, entityType, code string) (*Definition, error) {
	const q = `
SELECT RAWTOHEX(id), RAWTOHEX(tenant_id), entity_type, code, name,
       NVL(description, ''), current_version, status, created_at
FROM metadata_definitions
WHERE tenant_id = HEXTORAW(:1) AND entity_type = :2 AND UPPER(code) = UPPER(:3)`

	var item Definition
	err := r.db.QueryRowContext(ctx, q, strings.ToUpper(tenantID), entityType, code).Scan(
		&item.ID, &item.TenantID, &item.EntityType, &item.Code, &item.Name,
		&item.Description, &item.CurrentVersion, &item.Status, &item.CreatedAt,
	)
	if err != nil {
		return nil, err
	}
	item.ID = strings.ToLower(item.ID)
	item.TenantID = strings.ToLower(item.TenantID)
	return &item, nil
}

func (r *Repository) ListVersions(ctx context.Context, tenantID, definitionID string) ([]SchemaVersion, error) {
	const q = `
SELECT RAWTOHEX(id), RAWTOHEX(tenant_id), RAWTOHEX(definition_id), version_no,
       NVL(changelog, ''), is_current, created_at
FROM schema_versions
WHERE tenant_id = HEXTORAW(:1) AND definition_id = HEXTORAW(:2)
ORDER BY version_no`

	rows, err := r.db.QueryContext(ctx, q, strings.ToUpper(tenantID), strings.ToUpper(definitionID))
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var out []SchemaVersion
	for rows.Next() {
		var item SchemaVersion
		var isCurrent string
		if err := rows.Scan(
			&item.ID, &item.TenantID, &item.DefinitionID, &item.VersionNo,
			&item.Changelog, &isCurrent, &item.CreatedAt,
		); err != nil {
			return nil, err
		}
		item.ID = strings.ToLower(item.ID)
		item.TenantID = strings.ToLower(item.TenantID)
		item.DefinitionID = strings.ToLower(item.DefinitionID)
		item.IsCurrent = isCurrent == "Y"
		out = append(out, item)
	}
	return out, rows.Err()
}

func (r *Repository) GetCurrentVersion(ctx context.Context, tenantID, definitionID string) (*SchemaVersion, error) {
	const q = `
SELECT RAWTOHEX(id), RAWTOHEX(tenant_id), RAWTOHEX(definition_id), version_no,
       NVL(changelog, ''), is_current, created_at
FROM schema_versions
WHERE tenant_id = HEXTORAW(:1) AND definition_id = HEXTORAW(:2) AND is_current = 'Y'`

	var item SchemaVersion
	var isCurrent string
	err := r.db.QueryRowContext(ctx, q, strings.ToUpper(tenantID), strings.ToUpper(definitionID)).Scan(
		&item.ID, &item.TenantID, &item.DefinitionID, &item.VersionNo,
		&item.Changelog, &isCurrent, &item.CreatedAt,
	)
	if err == sql.ErrNoRows {
		return nil, ERR_NOT_FOUND
	}
	if err != nil {
		return nil, err
	}
	item.ID = strings.ToLower(item.ID)
	item.TenantID = strings.ToLower(item.TenantID)
	item.DefinitionID = strings.ToLower(item.DefinitionID)
	item.IsCurrent = isCurrent == "Y"
	return &item, nil
}

func (r *Repository) CreateNextVersion(ctx context.Context, tenantID, definitionID, changelog string) (*SchemaVersion, error) {
	current, err := r.GetCurrentVersion(ctx, tenantID, definitionID)
	if err != nil {
		return nil, err
	}

	nextNo := current.VersionNo + 1

	_, err = r.db.ExecContext(ctx, `
UPDATE schema_versions SET is_current = 'N'
WHERE tenant_id = HEXTORAW(:1) AND definition_id = HEXTORAW(:2)`,
		strings.ToUpper(tenantID), strings.ToUpper(definitionID))
	if err != nil {
		return nil, err
	}

	_, err = r.db.ExecContext(ctx, `
INSERT INTO schema_versions (tenant_id, definition_id, version_no, changelog, is_current)
VALUES (HEXTORAW(:1), HEXTORAW(:2), :3, :4, 'Y')`,
		strings.ToUpper(tenantID), strings.ToUpper(definitionID), nextNo, changelog)
	if err != nil {
		return nil, err
	}

	_, err = r.db.ExecContext(ctx, `
UPDATE metadata_definitions
SET current_version = :1, updated_at = SYSTIMESTAMP
WHERE id = HEXTORAW(:2) AND tenant_id = HEXTORAW(:3)`,
		nextNo, strings.ToUpper(definitionID), strings.ToUpper(tenantID))
	if err != nil {
		return nil, err
	}

	oldFields, err := r.ListFields(ctx, tenantID, current.ID)
	if err != nil {
		return nil, err
	}
	newVersion, err := r.GetCurrentVersion(ctx, tenantID, definitionID)
	if err != nil {
		return nil, err
	}
	for _, field := range oldFields {
		req := CreateFieldRequest{
			FieldKey:     field.FieldKey,
			Label:        field.Label,
			DataType:     field.DataType,
			IsRequired:   field.IsRequired,
			MinValue:     field.MinValue,
			MaxValue:     field.MaxValue,
			RegexPattern: field.RegexPattern,
			EnumValues:   field.EnumValues,
			DefaultValue: field.DefaultValue,
			SortOrder:    field.SortOrder,
		}
		if _, err := r.CreateField(ctx, tenantID, field.DefinitionID, newVersion.ID, req); err != nil {
			return nil, err
		}
	}
	return newVersion, nil
}

func (r *Repository) ListFields(ctx context.Context, tenantID, versionID string) ([]Field, error) {
	const q = `
SELECT RAWTOHEX(id), RAWTOHEX(tenant_id), RAWTOHEX(definition_id), RAWTOHEX(schema_version_id),
       field_key, label, data_type, is_required, min_value, max_value,
       NVL(regex_pattern, ''), enum_values, NVL(default_value, ''), sort_order
FROM metadata_fields
WHERE tenant_id = HEXTORAW(:1) AND schema_version_id = HEXTORAW(:2)
ORDER BY sort_order, field_key`

	rows, err := r.db.QueryContext(ctx, q, strings.ToUpper(tenantID), strings.ToUpper(versionID))
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var out []Field
	for rows.Next() {
		var item Field
		var isRequired string
		var minValue, maxValue sql.NullFloat64
		var enumRaw sql.NullString
		if err := rows.Scan(
			&item.ID, &item.TenantID, &item.DefinitionID, &item.SchemaVersionID,
			&item.FieldKey, &item.Label, &item.DataType, &isRequired, &minValue, &maxValue,
			&item.RegexPattern, &enumRaw, &item.DefaultValue, &item.SortOrder,
		); err != nil {
			return nil, err
		}
		item.ID = strings.ToLower(item.ID)
		item.TenantID = strings.ToLower(item.TenantID)
		item.DefinitionID = strings.ToLower(item.DefinitionID)
		item.SchemaVersionID = strings.ToLower(item.SchemaVersionID)
		item.IsRequired = isRequired == "Y"
		if minValue.Valid {
			v := minValue.Float64
			item.MinValue = &v
		}
		if maxValue.Valid {
			v := maxValue.Float64
			item.MaxValue = &v
		}
		if enumRaw.Valid && enumRaw.String != "" {
			_ = json.Unmarshal([]byte(enumRaw.String), &item.EnumValues)
		}
		out = append(out, item)
	}
	return out, rows.Err()
}

func (r *Repository) CreateField(ctx context.Context, tenantID, definitionID, versionID string, req CreateFieldRequest) (*Field, error) {
	required := "N"
	if req.IsRequired {
		required = "Y"
	}
	enumJSON := "[]"
	if len(req.EnumValues) > 0 {
		b, err := json.Marshal(req.EnumValues)
		if err != nil {
			return nil, err
		}
		enumJSON = string(b)
	}

	const q = `
INSERT INTO metadata_fields (
  tenant_id, definition_id, schema_version_id, field_key, label, data_type,
  is_required, min_value, max_value, regex_pattern, enum_values, default_value, sort_order
) VALUES (
  HEXTORAW(:1), HEXTORAW(:2), HEXTORAW(:3), :4, :5, :6,
  :7, :8, :9, :10, :11, :12, :13
)`
	_, err := r.db.ExecContext(ctx, q,
		strings.ToUpper(tenantID), strings.ToUpper(definitionID), strings.ToUpper(versionID),
		req.FieldKey, req.Label, req.DataType,
		required, req.MinValue, req.MaxValue, req.RegexPattern, enumJSON, req.DefaultValue, req.SortOrder,
	)
	if err != nil {
		return nil, err
	}

	fields, err := r.ListFields(ctx, tenantID, versionID)
	if err != nil {
		return nil, err
	}
	for _, field := range fields {
		if strings.EqualFold(field.FieldKey, req.FieldKey) {
			return &field, nil
		}
	}
	return nil, fmt.Errorf("field created but not found")
}

func (r *Repository) ListTenantFeatures(ctx context.Context, tenantID string) ([]TenantFeature, error) {
	const q = `
SELECT f.code, f.name, NVL(f.description, ''), f.category, NVL(tf.enabled, 'N')
FROM feature_definitions f
LEFT JOIN tenant_features tf
  ON tf.feature_id = f.id AND tf.tenant_id = HEXTORAW(:1)
ORDER BY f.code`

	rows, err := r.db.QueryContext(ctx, q, strings.ToUpper(tenantID))
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var out []TenantFeature
	for rows.Next() {
		var item TenantFeature
		var enabled string
		if err := rows.Scan(&item.Code, &item.Name, &item.Description, &item.Category, &enabled); err != nil {
			return nil, err
		}
		item.Enabled = enabled == "Y"
		out = append(out, item)
	}
	return out, rows.Err()
}

func (r *Repository) CountDefinitions(ctx context.Context, tenantID string) (int, error) {
	var count int
	err := r.db.QueryRowContext(ctx, `
SELECT COUNT(*) FROM metadata_definitions WHERE tenant_id = HEXTORAW(:1)`,
		strings.ToUpper(tenantID)).Scan(&count)
	return count, err
}

func (r *Repository) GetTenantIDByCode(ctx context.Context, tenantCode string) (string, error) {
	var id string
	err := r.db.QueryRowContext(ctx, `
SELECT RAWTOHEX(id) FROM tenants WHERE UPPER(code) = UPPER(:1)`, tenantCode).Scan(&id)
	if err == sql.ErrNoRows {
		return "", ERR_NOT_FOUND
	}
	if err != nil {
		return "", err
	}
	return strings.ToLower(id), nil
}
