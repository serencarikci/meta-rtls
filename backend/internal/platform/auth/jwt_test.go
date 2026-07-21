package auth

import "testing"

func TestIssueAndParseToken(t *testing.T) {
	tokens := NewTokenService("test-secret-at-least-32-characters-long", 60)
	signed, exp, err := tokens.Issue("user-1", "tenant-1", "admin@demo.test", "ADMIN")
	if err != nil {
		t.Fatal(err)
	}
	if signed == "" {
		t.Fatal("expected token string")
	}
	if exp.IsZero() {
		t.Fatal("expected expiry")
	}

	claims, err := tokens.Parse(signed)
	if err != nil {
		t.Fatal(err)
	}
	if claims.UserID != "user-1" || claims.TenantID != "tenant-1" {
		t.Fatalf("unexpected claims: %+v", claims)
	}
	if claims.Email != "admin@demo.test" || claims.Role != "ADMIN" {
		t.Fatalf("unexpected claims: %+v", claims)
	}
}

func TestParseRejectsBadToken(t *testing.T) {
	tokens := NewTokenService("test-secret-at-least-32-characters-long", 60)
	if _, err := tokens.Parse("not-a-jwt"); err == nil {
		t.Fatal("expected parse error")
	}
}
