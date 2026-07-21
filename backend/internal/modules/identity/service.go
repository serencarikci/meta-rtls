package identity

import (
	"context"
	"errors"
	"fmt"

	"github.com/denizyetis/meta-rtls/internal/platform/auth"
	"golang.org/x/crypto/bcrypt"
)

const DEMO_PASSWORD = "MetaRTLS!2026"

var (
	ERR_INVALID_CREDENTIALS = errors.New("invalid credentials")
	ERR_NOT_FOUND           = errors.New("user not found")

	DEMO_USERS = []struct {
		tenant string
		email  string
		name   string
		role   string
	}{
		{"warehouse-s", "admin@warehouse-s.demo", "Warehouse Admin", "TENANT_ADMIN"},
		{"hospital-m", "admin@hospital-m.demo", "Hospital Admin", "TENANT_ADMIN"},
		{"factory-l", "admin@factory-l.demo", "Factory Admin", "TENANT_ADMIN"},
	}
)

type Service struct {
	repo   *Repository
	tokens *auth.TokenService
}

func NewService(repo *Repository, tokens *auth.TokenService) *Service {
	return &Service{repo: repo, tokens: tokens}
}

func (s *Service) Login(ctx context.Context, req LoginRequest) (*LoginResponse, error) {
	user, _, _, err := s.repo.FindByTenantCodeAndEmail(ctx, req.TenantCode, req.Email)
	if err != nil {
		return nil, err
	}
	if err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(req.Password)); err != nil {
		return nil, ERR_INVALID_CREDENTIALS
	}
	token, exp, err := s.tokens.Issue(user.ID, user.TenantID, user.Email, user.Role)
	if err != nil {
		return nil, fmt.Errorf("issue token: %w", err)
	}
	user.PasswordHash = ""
	return &LoginResponse{AccessToken: token, ExpiresAt: exp, User: *user}, nil
}

func (s *Service) Me(ctx context.Context, userID string) (*MeResponse, error) {
	user, tenantCode, tenantName, err := s.repo.FindByID(ctx, userID)
	if err != nil {
		return nil, err
	}
	user.PasswordHash = ""
	return &MeResponse{User: *user, TenantCode: tenantCode, TenantName: tenantName}, nil
}

func (s *Service) BootstrapDemoUsers(ctx context.Context) error {
	hash, err := bcrypt.GenerateFromPassword([]byte(DEMO_PASSWORD), bcrypt.DefaultCost)
	if err != nil {
		return err
	}
	for _, demoUser := range DEMO_USERS {
		if err := s.repo.EnsureDemoUser(ctx, demoUser.tenant, demoUser.email, demoUser.name, demoUser.role, string(hash)); err != nil {
			return err
		}
	}
	return nil
}
