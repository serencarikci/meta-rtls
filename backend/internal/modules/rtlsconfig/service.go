package rtlsconfig

import "context"

type Service struct {
	repo *Repository
}

func NewService(repo *Repository) *Service {
	return &Service{repo: repo}
}

func (s *Service) ListSites(ctx context.Context, tenantID string) ([]Site, error) {
	return s.repo.ListSites(ctx, tenantID)
}

func (s *Service) CreateSite(ctx context.Context, tenantID string, req CreateSiteRequest) (*Site, error) {
	return s.repo.CreateSite(ctx, tenantID, req)
}

func (s *Service) ListBuildings(ctx context.Context, tenantID string) ([]Building, error) {
	return s.repo.ListBuildings(ctx, tenantID)
}

func (s *Service) ListFloors(ctx context.Context, tenantID string) ([]Floor, error) {
	return s.repo.ListFloors(ctx, tenantID)
}

func (s *Service) ListZones(ctx context.Context, tenantID, floorID string) ([]Zone, error) {
	return s.repo.ListZones(ctx, tenantID, floorID)
}
