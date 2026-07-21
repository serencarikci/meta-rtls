package tenant

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

func (r *Repository) List(ctx context.Context) ([]Tenant, error) {
	const q = `
SELECT RAWTOHEX(id), code, name, profile_scale, status, created_at
FROM tenants
ORDER BY name`

	rows, err := r.db.QueryContext(ctx, q)
	if err != nil {
		return nil, fmt.Errorf("list tenants: %w", err)
	}
	defer rows.Close()

	var out []Tenant
	for rows.Next() {
		var t Tenant
		if err := rows.Scan(&t.ID, &t.Code, &t.Name, &t.ProfileScale, &t.Status, &t.CreatedAt); err != nil {
			return nil, err
		}
		t.ID = strings.ToLower(t.ID)
		out = append(out, t)
	}
	return out, rows.Err()
}

func (r *Repository) GetByID(ctx context.Context, id string) (*Tenant, error) {
	const q = `
SELECT RAWTOHEX(id), code, name, profile_scale, status, created_at
FROM tenants WHERE id = HEXTORAW(:1)`

	var t Tenant
	err := r.db.QueryRowContext(ctx, q, strings.ToUpper(id)).Scan(
		&t.ID, &t.Code, &t.Name, &t.ProfileScale, &t.Status, &t.CreatedAt,
	)
	if err == sql.ErrNoRows {
		return nil, ERR_NOT_FOUND
	}
	if err != nil {
		return nil, err
	}
	t.ID = strings.ToLower(t.ID)
	return &t, nil
}

func (r *Repository) Create(ctx context.Context, req CreateTenantRequest) (*Tenant, error) {
	const insertQ = `INSERT INTO tenants (code, name, profile_scale, status) VALUES (:1, :2, :3, 'ACTIVE')`
	if _, err := r.db.ExecContext(ctx, insertQ, req.Code, req.Name, req.ProfileScale); err != nil {
		return nil, fmt.Errorf("create tenant: %w", err)
	}

	const selectQ = `
SELECT RAWTOHEX(id), code, name, profile_scale, status, created_at
FROM tenants WHERE UPPER(code) = UPPER(:1)`
	var t Tenant
	if err := r.db.QueryRowContext(ctx, selectQ, req.Code).Scan(
		&t.ID, &t.Code, &t.Name, &t.ProfileScale, &t.Status, &t.CreatedAt,
	); err != nil {
		return nil, err
	}
	t.ID = strings.ToLower(t.ID)
	return &t, nil
}
