package tenant

import (
	"context"
	"errors"
)

var ERR_NOT_FOUND = errors.New("tenant not found")

type Service struct {
	repo *Repository
}

func NewService(repo *Repository) *Service {
	return &Service{repo: repo}
}

func (s *Service) List(ctx context.Context) ([]Tenant, error) {
	return s.repo.List(ctx)
}

func (s *Service) Get(ctx context.Context, id string) (*Tenant, error) {
	return s.repo.GetByID(ctx, id)
}

func (s *Service) Create(ctx context.Context, req CreateTenantRequest) (*Tenant, error) {
	return s.repo.Create(ctx, req)
}
