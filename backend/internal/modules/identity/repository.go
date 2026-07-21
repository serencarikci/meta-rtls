package identity

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"strings"
)

type Repository struct {
	db *sql.DB
}

func NewRepository(db *sql.DB) *Repository {
	return &Repository{db: db}
}

func (r *Repository) FindByTenantCodeAndEmail(ctx context.Context, tenantCode, email string) (*User, string, string, error) {
	const q = `
SELECT
  RAWTOHEX(u.id),
  RAWTOHEX(u.tenant_id),
  u.email,
  u.display_name,
  u.role,
  u.status,
  u.password_hash,
  u.created_at,
  t.code,
  t.name
FROM users u
JOIN tenants t ON t.id = u.tenant_id
WHERE UPPER(t.code) = UPPER(:1)
  AND LOWER(u.email) = LOWER(:2)
  AND u.status = 'ACTIVE'
  AND t.status = 'ACTIVE'`

	var u User
	var tenantCodeOut, tenantName string
	err := r.db.QueryRowContext(ctx, q, tenantCode, email).Scan(
		&u.ID, &u.TenantID, &u.Email, &u.DisplayName, &u.Role, &u.Status, &u.PasswordHash, &u.CreatedAt,
		&tenantCodeOut, &tenantName,
	)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, "", "", ERR_INVALID_CREDENTIALS
	}
	if err != nil {
		return nil, "", "", fmt.Errorf("find user: %w", err)
	}
	u.ID = strings.ToLower(u.ID)
	u.TenantID = strings.ToLower(u.TenantID)
	return &u, tenantCodeOut, tenantName, nil
}

func (r *Repository) FindByID(ctx context.Context, userID string) (*User, string, string, error) {
	const q = `
SELECT
  RAWTOHEX(u.id),
  RAWTOHEX(u.tenant_id),
  u.email,
  u.display_name,
  u.role,
  u.status,
  u.password_hash,
  u.created_at,
  t.code,
  t.name
FROM users u
JOIN tenants t ON t.id = u.tenant_id
WHERE u.id = HEXTORAW(:1)`

	var u User
	var tenantCode, tenantName string
	err := r.db.QueryRowContext(ctx, q, strings.ToUpper(userID)).Scan(
		&u.ID, &u.TenantID, &u.Email, &u.DisplayName, &u.Role, &u.Status, &u.PasswordHash, &u.CreatedAt,
		&tenantCode, &tenantName,
	)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, "", "", ERR_NOT_FOUND
	}
	if err != nil {
		return nil, "", "", fmt.Errorf("find user by id: %w", err)
	}
	u.ID = strings.ToLower(u.ID)
	u.TenantID = strings.ToLower(u.TenantID)
	return &u, tenantCode, tenantName, nil
}

func (r *Repository) EnsureDemoUser(ctx context.Context, tenantCode, email, displayName, role, passwordHash string) error {
	var tenantID string
	err := r.db.QueryRowContext(ctx,
		`SELECT RAWTOHEX(id) FROM tenants WHERE UPPER(code) = UPPER(:1)`, tenantCode,
	).Scan(&tenantID)
	if errors.Is(err, sql.ErrNoRows) {
		return fmt.Errorf("tenant %s not found for demo seed", tenantCode)
	}
	if err != nil {
		return err
	}

	var existing string
	err = r.db.QueryRowContext(ctx, `
SELECT RAWTOHEX(id) FROM users
WHERE tenant_id = HEXTORAW(:1) AND LOWER(email) = LOWER(:2)`,
		tenantID, email,
	).Scan(&existing)

	if errors.Is(err, sql.ErrNoRows) {
		_, err = r.db.ExecContext(ctx, `
INSERT INTO users (tenant_id, email, password_hash, display_name, role, status)
VALUES (HEXTORAW(:1), :2, :3, :4, :5, 'ACTIVE')`,
			tenantID, email, passwordHash, displayName, role)
		return err
	}
	if err != nil {
		return err
	}

	_, err = r.db.ExecContext(ctx, `
UPDATE users
SET password_hash = :1, display_name = :2, role = :3, status = 'ACTIVE', updated_at = SYSTIMESTAMP
WHERE id = HEXTORAW(:4)`,
		passwordHash, displayName, role, existing)
	return err
}
