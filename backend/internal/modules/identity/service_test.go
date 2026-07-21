package identity

import (
	"testing"

	"golang.org/x/crypto/bcrypt"
)

func TestDemoPasswordHashRoundTrip(t *testing.T) {
	hash, err := bcrypt.GenerateFromPassword([]byte(DEMO_PASSWORD), bcrypt.MinCost)
	if err != nil {
		t.Fatal(err)
	}
	if err := bcrypt.CompareHashAndPassword(hash, []byte(DEMO_PASSWORD)); err != nil {
		t.Fatal(err)
	}
}
