package config

import "testing"

func TestValidateAllowsDevDefaultSecret(t *testing.T) {
	cfg := &Config{
		AppEnv:    "development",
		JWTSecret: "dev-only-secret-change-me",
	}
	if err := cfg.Validate(); err != nil {
		t.Fatal(err)
	}
}

func TestValidateRejectsWeakProductionSecret(t *testing.T) {
	cfg := &Config{
		AppEnv:    "production",
		JWTSecret: "short",
	}
	if err := cfg.Validate(); err == nil {
		t.Fatal("expected error for short production secret")
	}
}

func TestValidateRejectsDefaultProductionSecret(t *testing.T) {
	cfg := &Config{
		AppEnv:    "production",
		JWTSecret: "change-me-in-production-use-long-random-string",
	}
	if err := cfg.Validate(); err == nil {
		t.Fatal("expected error for default production secret")
	}
}
