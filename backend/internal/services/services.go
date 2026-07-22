package services

import (
	"github.com/denizyetis/meta-rtls/internal/config"
	"github.com/denizyetis/meta-rtls/internal/version"
)

type VersionInfo struct {
	Service string `json:"service"`
	Version string `json:"version"`
}

type PublicConfig struct {
	Service           string   `json:"service"`
	Version           string   `json:"version"`
	AppEnv            string   `json:"appEnv"`
	AppPort           string   `json:"appPort"`
	JWTTTLMinutes     int      `json:"jwtTtlMinutes"`
	MQTTBroker        string   `json:"mqttBroker"`
	MQTTClientID      string   `json:"mqttClientId"`
	MQTTTopicLocation string   `json:"mqttTopicLocation"`
	CORSOrigins       []string `json:"corsOrigins"`
}

type Services struct {
	cfg *config.Config
}

func NewServices(cfg *config.Config) *Services {
	return &Services{cfg: cfg}
}

func (s *Services) GetVersion() VersionInfo {
	return VersionInfo{
		Service: version.SERVICE_NAME,
		Version: version.VERSION,
	}
}

func (s *Services) GetConfig() PublicConfig {
	origins := append([]string{}, s.cfg.CORSOrigins...)
	return PublicConfig{
		Service:           version.SERVICE_NAME,
		Version:           version.VERSION,
		AppEnv:            s.cfg.AppEnv,
		AppPort:           s.cfg.AppPort,
		JWTTTLMinutes:     s.cfg.JWTTTLMinutes,
		MQTTBroker:        s.cfg.MQTTBroker,
		MQTTClientID:      s.cfg.MQTTClientID,
		MQTTTopicLocation: s.cfg.MQTTTopic,
		CORSOrigins:       origins,
	}
}
