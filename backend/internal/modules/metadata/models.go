package metadata

import "time"

type Definition struct {
	ID             string    `json:"id"`
	TenantID       string    `json:"tenantId"`
	EntityType     string    `json:"entityType"`
	Code           string    `json:"code"`
	Name           string    `json:"name"`
	Description    string    `json:"description"`
	CurrentVersion int       `json:"currentVersion"`
	Status         string    `json:"status"`
	CreatedAt      time.Time `json:"createdAt"`
}

type SchemaVersion struct {
	ID           string    `json:"id"`
	TenantID     string    `json:"tenantId"`
	DefinitionID string    `json:"definitionId"`
	VersionNo    int       `json:"versionNo"`
	Changelog    string    `json:"changelog"`
	IsCurrent    bool      `json:"isCurrent"`
	CreatedAt    time.Time `json:"createdAt"`
}

type Field struct {
	ID              string   `json:"id"`
	TenantID        string   `json:"tenantId"`
	DefinitionID    string   `json:"definitionId"`
	SchemaVersionID string   `json:"schemaVersionId"`
	FieldKey        string   `json:"fieldKey"`
	Label           string   `json:"label"`
	DataType        string   `json:"dataType"`
	IsRequired      bool     `json:"isRequired"`
	MinValue        *float64 `json:"minValue,omitempty"`
	MaxValue        *float64 `json:"maxValue,omitempty"`
	RegexPattern    string   `json:"regexPattern,omitempty"`
	EnumValues      []string `json:"enumValues,omitempty"`
	DefaultValue    string   `json:"defaultValue,omitempty"`
	SortOrder       int      `json:"sortOrder"`
}

type CreateDefinitionRequest struct {
	EntityType  string `json:"entityType" binding:"required"`
	Code        string `json:"code" binding:"required"`
	Name        string `json:"name" binding:"required"`
	Description string `json:"description"`
}

type CreateFieldRequest struct {
	FieldKey     string   `json:"fieldKey" binding:"required"`
	Label        string   `json:"label" binding:"required"`
	DataType     string   `json:"dataType" binding:"required"`
	IsRequired   bool     `json:"isRequired"`
	MinValue     *float64 `json:"minValue"`
	MaxValue     *float64 `json:"maxValue"`
	RegexPattern string   `json:"regexPattern"`
	EnumValues   []string `json:"enumValues"`
	DefaultValue string   `json:"defaultValue"`
	SortOrder    int      `json:"sortOrder"`
}

type CreateVersionRequest struct {
	Changelog string `json:"changelog"`
}

type ValidateRequest struct {
	DefinitionID string         `json:"definitionId" binding:"required"`
	Values       map[string]any `json:"values" binding:"required"`
}

type ValidateResponse struct {
	Valid  bool     `json:"valid"`
	Errors []string `json:"errors"`
}

type TenantFeature struct {
	Code        string `json:"code"`
	Name        string `json:"name"`
	Description string `json:"description"`
	Category    string `json:"category"`
	Enabled     bool   `json:"enabled"`
}
