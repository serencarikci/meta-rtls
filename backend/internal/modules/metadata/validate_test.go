package metadata

import "testing"

func TestValidateValuesRequiredAndNumber(t *testing.T) {
	minV := 0.0
	maxV := 100.0
	fields := []Field{
		{FieldKey: "batteryLevel", DataType: "NUMBER", IsRequired: true, MinValue: &minV, MaxValue: &maxV},
		{FieldKey: "department", DataType: "ENUM", IsRequired: true, EnumValues: []string{"Warehouse", "Yard"}},
	}

	errs := ValidateValues(fields, map[string]any{
		"batteryLevel": 120,
		"department":   "Office",
	})
	if len(errs) != 2 {
		t.Fatalf("expected 2 errors, got %v", errs)
	}

	ok := ValidateValues(fields, map[string]any{
		"batteryLevel": 80,
		"department":   "Warehouse",
	})
	if len(ok) != 0 {
		t.Fatalf("expected no errors, got %v", ok)
	}
}
