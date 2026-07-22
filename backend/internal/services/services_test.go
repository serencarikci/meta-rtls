package services

import (
	"testing"

	"github.com/denizyetis/meta-rtls/internal/config"
)

func TestGetVersionAndConfig(t *testing.T) {
	cfg := &config.Config{
		AppEnv:        "development",
		AppPort:       "8090",
		JWTTTLMinutes: 480,
		MQTTBroker:    "tcp://localhost:1883",
		MQTTClientID:  "metartls-api",
		MQTTTopic:     "rtls/+/location",
		CORSOrigins:   []string{"http://localhost:5173"},
	}
	svc := NewServices(cfg)

	ver := svc.GetVersion()
	if ver.Service != "metartls-api" || ver.Version == "" {
		t.Fatalf("unexpected version: %+v", ver)
	}

	pub := svc.GetConfig()
	if pub.AppPort != "8090" || pub.AppEnv != "development" {
		t.Fatalf("unexpected config: %+v", pub)
	}
	if len(pub.CORSOrigins) != 1 {
		t.Fatalf("expected cors origins, got %+v", pub.CORSOrigins)
	}
}
