package metadata

import (
	"fmt"
	"regexp"
	"strconv"
)

var ALLOWED_ENTITY_TYPES = map[string]bool{
	"ASSET": true, "PERSON": true, "TAG": true, "ZONE": true, "DEVICE": true,
}

var ALLOWED_DATA_TYPES = map[string]bool{
	"STRING": true, "NUMBER": true, "BOOLEAN": true, "ENUM": true, "DATE": true, "JSON": true,
}

func ValidateValues(fields []Field, values map[string]any) []string {
	var errors []string

	for _, field := range fields {
		raw, ok := values[field.FieldKey]
		if !ok || raw == nil || raw == "" {
			if field.IsRequired {
				errors = append(errors, field.FieldKey+": required")
			}
			continue
		}

		switch field.DataType {
		case "STRING", "DATE", "JSON":
			text := fmt.Sprint(raw)
			if field.RegexPattern != "" {
				matched, err := regexp.MatchString(field.RegexPattern, text)
				if err != nil || !matched {
					errors = append(errors, field.FieldKey+": regex failed")
				}
			}
		case "NUMBER":
			num, err := toFloat(raw)
			if err != nil {
				errors = append(errors, field.FieldKey+": must be a number")
				continue
			}
			if field.MinValue != nil && num < *field.MinValue {
				errors = append(errors, field.FieldKey+": below min")
			}
			if field.MaxValue != nil && num > *field.MaxValue {
				errors = append(errors, field.FieldKey+": above max")
			}
		case "BOOLEAN":
			if _, ok := raw.(bool); !ok {
				text := fmt.Sprint(raw)
				if text != "true" && text != "false" {
					errors = append(errors, field.FieldKey+": must be boolean")
				}
			}
		case "ENUM":
			text := fmt.Sprint(raw)
			found := false
			for _, item := range field.EnumValues {
				if item == text {
					found = true
					break
				}
			}
			if !found {
				errors = append(errors, field.FieldKey+": not in enum list")
			}
		}
	}

	return errors
}

func toFloat(raw any) (float64, error) {
	switch v := raw.(type) {
	case float64:
		return v, nil
	case float32:
		return float64(v), nil
	case int:
		return float64(v), nil
	case int64:
		return float64(v), nil
	case string:
		return strconv.ParseFloat(v, 64)
	default:
		return 0, fmt.Errorf("not a number")
	}
}
