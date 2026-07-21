package analysis

import "testing"

func TestBuildImpactResultLowRiskAdd(t *testing.T) {
	result := buildImpactResult(ImpactRequest{
		RequestType: "ADD_METADATA_FIELD",
		Title:       "Add color",
		DataType:    "STRING",
		IsRequired:  false,
	}, 3, 10, 20)

	if result.RiskLevel != "LOW" {
		t.Fatalf("expected LOW, got %s score=%.0f", result.RiskLevel, result.ComplexityScore)
	}
	if result.MigrationRequired {
		t.Fatal("optional add should not need migration")
	}
	if result.AffectedEntities != 20 {
		t.Fatalf("expected entities from expected tags, got %d", result.AffectedEntities)
	}
}

func TestBuildImpactResultHighRiskRemove(t *testing.T) {
	result := buildImpactResult(ImpactRequest{
		RequestType: "REMOVE_METADATA_FIELD",
		Title:       "Remove battery",
		DataType:    "JSON",
		IsRequired:  true,
	}, 3, 50, 2000)

	if result.RiskLevel != "HIGH" && result.RiskLevel != "CRITICAL" {
		t.Fatalf("expected HIGH or CRITICAL, got %s score=%.0f", result.RiskLevel, result.ComplexityScore)
	}
	if !result.MigrationRequired {
		t.Fatal("remove required field needs migration")
	}
	if result.BackwardCompatible {
		t.Fatal("migration change is not backward compatible")
	}
}
