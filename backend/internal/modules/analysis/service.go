package analysis

import (
	"context"
	"fmt"
	"strings"
)

type Service struct {
	repo *Repository
}

func NewService(repo *Repository) *Service {
	return &Service{repo: repo}
}

func (s *Service) ListRequirements(ctx context.Context) ([]Requirement, error) {
	items, err := s.repo.ListRequirements(ctx)
	if err != nil {
		return nil, err
	}
	if items == nil {
		items = []Requirement{}
	}
	return items, nil
}

func (s *Service) CompareProfiles(ctx context.Context) (*CompareResponse, error) {
	profiles, err := s.repo.ListProfiles(ctx)
	if err != nil {
		return nil, err
	}
	if profiles == nil {
		profiles = []TenantProfile{}
	}

	notes := []string{
		"SMALL tenants need fewer tags and shorter retention.",
		"MEDIUM tenants add more features and higher event rate.",
		"LARGE tenants need impact analysis and long history.",
	}
	return &CompareResponse{Profiles: profiles, Notes: notes}, nil
}

func (s *Service) AnalyzeImpact(ctx context.Context, tenantID string, req ImpactRequest) (*ImpactResult, error) {
	req.RequestType = strings.ToUpper(strings.TrimSpace(req.RequestType))
	if req.RequestType == "" {
		req.RequestType = "ADD_METADATA_FIELD"
	}

	tenantCount, err := s.repo.CountTenants(ctx)
	if err != nil {
		return nil, err
	}
	liveEntities, err := s.repo.CountTagsAndAssets(ctx)
	if err != nil {
		return nil, err
	}
	expectedTags, err := s.repo.SumExpectedTags(ctx)
	if err != nil {
		return nil, err
	}

	affectedEntities := liveEntities
	if expectedTags > affectedEntities {
		affectedEntities = expectedTags
	}
	if affectedEntities == 0 {
		affectedEntities = 100
	}

	score := 10.0
	if req.IsRequired {
		score += 25
	}
	if req.RequestType == "ADD_METADATA_FIELD" {
		score += 15
	}
	if req.RequestType == "REMOVE_METADATA_FIELD" {
		score += 35
	}
	if req.DataType == "JSON" {
		score += 10
	}
	if affectedEntities > 1000 {
		score += 20
	}
	if affectedEntities > 10000 {
		score += 30
	}

	risk := "LOW"
	if score >= 40 {
		risk = "MEDIUM"
	}
	if score >= 70 {
		risk = "HIGH"
	}
	if score >= 90 {
		risk = "CRITICAL"
	}

	migration := req.IsRequired || req.RequestType == "REMOVE_METADATA_FIELD" || req.RequestType == "CHANGE_FIELD_TYPE"
	backward := !migration

	result := &ImpactResult{
		Title:              req.Title,
		RequestType:        req.RequestType,
		AffectedTenants:    tenantCount,
		AffectedEntities:   affectedEntities,
		MigrationRequired:  migration,
		RiskLevel:          risk,
		ComplexityScore:    score,
		BackwardCompatible: backward,
		Summary: fmt.Sprintf(
			"Change '%s' touches about %d tenants and %d entities. Risk=%s, score=%.0f.",
			req.Title, tenantCount, affectedEntities, risk, score,
		),
	}

	if req.Save {
		id, err := s.repo.SaveChangeRequest(ctx, tenantID, *result, req)
		if err != nil {
			return nil, err
		}
		result.SavedID = id
		_ = s.repo.InsertAudit(ctx, tenantID, "IMPACT_ANALYZED", fmt.Sprintf(`{"title":%q,"risk":%q}`, req.Title, risk))
	}

	return result, nil
}

func (s *Service) ListChangeRequests(ctx context.Context, tenantID string) ([]ChangeRequest, error) {
	return s.repo.ListChangeRequests(ctx, tenantID)
}

func (s *Service) Bootstrap(ctx context.Context) error {
	return s.repo.EnsureDemoRequirements(ctx)
}
