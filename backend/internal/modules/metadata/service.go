package metadata

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

func (s *Service) ListDefinitions(ctx context.Context, tenantID string) ([]Definition, error) {
	items, err := s.repo.ListDefinitions(ctx, tenantID)
	if err != nil {
		return nil, err
	}
	if items == nil {
		items = []Definition{}
	}
	return items, nil
}

func (s *Service) GetDefinition(ctx context.Context, tenantID, id string) (*Definition, error) {
	return s.repo.GetDefinition(ctx, tenantID, id)
}

func (s *Service) CreateDefinition(ctx context.Context, tenantID string, req CreateDefinitionRequest) (*Definition, error) {
	req.EntityType = strings.ToUpper(strings.TrimSpace(req.EntityType))
	if !ALLOWED_ENTITY_TYPES[req.EntityType] {
		return nil, ERR_INVALID_ENTITY
	}
	return s.repo.CreateDefinition(ctx, tenantID, req)
}

func (s *Service) ListVersions(ctx context.Context, tenantID, definitionID string) ([]SchemaVersion, error) {
	if _, err := s.repo.GetDefinition(ctx, tenantID, definitionID); err != nil {
		return nil, err
	}
	items, err := s.repo.ListVersions(ctx, tenantID, definitionID)
	if err != nil {
		return nil, err
	}
	if items == nil {
		items = []SchemaVersion{}
	}
	return items, nil
}

func (s *Service) CreateVersion(ctx context.Context, tenantID, definitionID string, req CreateVersionRequest) (*SchemaVersion, error) {
	if _, err := s.repo.GetDefinition(ctx, tenantID, definitionID); err != nil {
		return nil, err
	}
	changelog := req.Changelog
	if changelog == "" {
		changelog = "new schema version"
	}
	return s.repo.CreateNextVersion(ctx, tenantID, definitionID, changelog)
}

func (s *Service) ListFields(ctx context.Context, tenantID, versionID string) ([]Field, error) {
	items, err := s.repo.ListFields(ctx, tenantID, versionID)
	if err != nil {
		return nil, err
	}
	if items == nil {
		items = []Field{}
	}
	return items, nil
}

func (s *Service) CreateField(ctx context.Context, tenantID, versionID string, req CreateFieldRequest) (*Field, error) {
	req.DataType = strings.ToUpper(strings.TrimSpace(req.DataType))
	if !ALLOWED_DATA_TYPES[req.DataType] {
		return nil, ERR_INVALID_TYPE
	}
	if req.DataType == "ENUM" && len(req.EnumValues) == 0 {
		return nil, fmt.Errorf("%w: enum needs values", ERR_BAD_REQUEST)
	}

	versionsOwned, err := s.findDefinitionIDForVersion(ctx, tenantID, versionID)
	if err != nil {
		return nil, err
	}
	return s.repo.CreateField(ctx, tenantID, versionsOwned, versionID, req)
}

func (s *Service) findDefinitionIDForVersion(ctx context.Context, tenantID, versionID string) (string, error) {
	defs, err := s.repo.ListDefinitions(ctx, tenantID)
	if err != nil {
		return "", err
	}
	for _, def := range defs {
		versions, err := s.repo.ListVersions(ctx, tenantID, def.ID)
		if err != nil {
			return "", err
		}
		for _, version := range versions {
			if strings.EqualFold(version.ID, versionID) {
				return def.ID, nil
			}
		}
	}
	return "", ERR_NOT_FOUND
}

func (s *Service) Validate(ctx context.Context, tenantID string, req ValidateRequest) (*ValidateResponse, error) {
	if _, err := s.repo.GetDefinition(ctx, tenantID, req.DefinitionID); err != nil {
		return nil, err
	}
	current, err := s.repo.GetCurrentVersion(ctx, tenantID, req.DefinitionID)
	if err != nil {
		return nil, err
	}
	fields, err := s.repo.ListFields(ctx, tenantID, current.ID)
	if err != nil {
		return nil, err
	}
	errs := ValidateValues(fields, req.Values)
	if errs == nil {
		errs = []string{}
	}
	return &ValidateResponse{Valid: len(errs) == 0, Errors: errs}, nil
}

func (s *Service) ListFeatures(ctx context.Context, tenantID string) ([]TenantFeature, error) {
	items, err := s.repo.ListTenantFeatures(ctx, tenantID)
	if err != nil {
		return nil, err
	}
	if items == nil {
		items = []TenantFeature{}
	}
	return items, nil
}

func (s *Service) BootstrapDemoMetadata(ctx context.Context) error {
	tenantID, err := s.repo.GetTenantIDByCode(ctx, "warehouse-s")
	if err != nil {
		return err
	}
	count, err := s.repo.CountDefinitions(ctx, tenantID)
	if err != nil {
		return err
	}
	if count > 0 {
		return nil
	}

	def, err := s.CreateDefinition(ctx, tenantID, CreateDefinitionRequest{
		EntityType:  "ASSET",
		Code:        "FORKLIFT",
		Name:        "Forklift asset fields",
		Description: "Demo metadata for warehouse forklifts",
	})
	if err != nil {
		return err
	}
	current, err := s.repo.GetCurrentVersion(ctx, tenantID, def.ID)
	if err != nil {
		return err
	}

	minLoad := 0.0
	maxLoad := 5000.0
	minBattery := 0.0
	maxBattery := 100.0

	fields := []CreateFieldRequest{
		{FieldKey: "loadCapacity", Label: "Load capacity (kg)", DataType: "NUMBER", IsRequired: true, MinValue: &minLoad, MaxValue: &maxLoad, SortOrder: 1},
		{FieldKey: "batteryLevel", Label: "Battery level (%)", DataType: "NUMBER", IsRequired: true, MinValue: &minBattery, MaxValue: &maxBattery, SortOrder: 2},
		{FieldKey: "department", Label: "Department", DataType: "ENUM", IsRequired: true, EnumValues: []string{"Warehouse", "Shipping", "Yard"}, SortOrder: 3},
	}
	for _, field := range fields {
		if _, err := s.repo.CreateField(ctx, tenantID, def.ID, current.ID, field); err != nil {
			return err
		}
	}
	return nil
}
