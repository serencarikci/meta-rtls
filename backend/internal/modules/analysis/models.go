package analysis

import "time"

type Requirement struct {
	ID            string    `json:"id"`
	TenantID      string    `json:"tenantId"`
	TenantCode    string    `json:"tenantCode"`
	ProfileScale  string    `json:"profileScale"`
	Code          string    `json:"code"`
	Title         string    `json:"title"`
	Description   string    `json:"description"`
	Priority      string    `json:"priority"`
	ExpectedTags  int       `json:"expectedTags"`
	ExpectedEPS   float64   `json:"expectedEps"`
	RetentionDays int       `json:"retentionDays"`
	CreatedAt     time.Time `json:"createdAt"`
}

type TenantProfile struct {
	TenantID         string   `json:"tenantId"`
	Code             string   `json:"code"`
	Name             string   `json:"name"`
	ProfileScale     string   `json:"profileScale"`
	FeatureCodes     []string `json:"featureCodes"`
	MetadataFields   int      `json:"metadataFields"`
	ExpectedTags     int      `json:"expectedTags"`
	ExpectedEPS      float64  `json:"expectedEps"`
	RetentionDays    int      `json:"retentionDays"`
	RequirementTitle string   `json:"requirementTitle"`
}

type CompareResponse struct {
	Profiles []TenantProfile `json:"profiles"`
	Notes    []string        `json:"notes"`
}

type ImpactRequest struct {
	RequestType string `json:"requestType" binding:"required"`
	Title       string `json:"title" binding:"required"`
	FieldKey    string `json:"fieldKey"`
	DataType    string `json:"dataType"`
	IsRequired  bool   `json:"isRequired"`
	Save        bool   `json:"save"`
}

type ImpactResult struct {
	Title              string  `json:"title"`
	RequestType        string  `json:"requestType"`
	AffectedTenants    int     `json:"affectedTenants"`
	AffectedEntities   int     `json:"affectedEntities"`
	MigrationRequired  bool    `json:"migrationRequired"`
	RiskLevel          string  `json:"riskLevel"`
	ComplexityScore    float64 `json:"complexityScore"`
	BackwardCompatible bool    `json:"backwardCompatible"`
	Summary            string  `json:"summary"`
	SavedID            string  `json:"savedId,omitempty"`
}

type ChangeRequest struct {
	ID                 string    `json:"id"`
	TenantID           string    `json:"tenantId"`
	RequestType        string    `json:"requestType"`
	Title              string    `json:"title"`
	AffectedTenants    int       `json:"affectedTenants"`
	AffectedEntities   int       `json:"affectedEntities"`
	MigrationRequired  bool      `json:"migrationRequired"`
	RiskLevel          string    `json:"riskLevel"`
	ComplexityScore    float64   `json:"complexityScore"`
	BackwardCompatible bool      `json:"backwardCompatible"`
	Status             string    `json:"status"`
	CreatedAt          time.Time `json:"createdAt"`
}
