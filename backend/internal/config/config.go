package config

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

type fileConfig struct {
	AppEnv            string   `json:"appEnv"`
	AppPort           string   `json:"appPort"`
	JWTSecret         string   `json:"jwtSecret"`
	JWTTTLMinutes     int      `json:"jwtTtlMinutes"`
	OracleUser        string   `json:"oracleUser"`
	OraclePassword    string   `json:"oraclePassword"`
	OracleDSN         string   `json:"oracleDsn"`
	OracleSysPassword string   `json:"oracleSysPassword"`
	RedisAddr         string   `json:"redisAddr"`
	RedisPassword     string   `json:"redisPassword"`
	MQTTBroker        string   `json:"mqttBroker"`
	MQTTClientID      string   `json:"mqttClientId"`
	MQTTTopicLocation string   `json:"mqttTopicLocation"`
	CORSOrigins       []string `json:"corsOrigins"`
}

type Config struct {
	AppEnv            string
	AppPort           string
	JWTSecret         string
	JWTTTLMinutes     int
	OracleUser        string
	OraclePass        string
	OracleDSNHost     string
	OracleSysPassword string
	RedisAddr         string
	RedisPassword     string
	MQTTBroker        string
	MQTTClientID      string
	MQTTTopic         string
	CORSOrigins       []string
}

func Load() (*Config, error) {
	path, err := resolveConfigPath()
	if err != nil {
		return nil, err
	}

	raw, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("read config %s: %w", path, err)
	}

	var file fileConfig
	if err := json.Unmarshal(raw, &file); err != nil {
		return nil, fmt.Errorf("parse config %s: %w", path, err)
	}

	origins := make([]string, 0, len(file.CORSOrigins))
	for _, origin := range file.CORSOrigins {
		origin = strings.TrimSpace(origin)
		if origin != "" {
			origins = append(origins, origin)
		}
	}
	if len(origins) == 0 {
		origins = []string{"http://localhost:5173"}
	}

	ttl := file.JWTTTLMinutes
	if ttl <= 0 {
		ttl = 480
	}

	cfg := &Config{
		AppEnv:            defaultString(file.AppEnv, "development"),
		AppPort:           defaultString(file.AppPort, "8090"),
		JWTSecret:         defaultString(file.JWTSecret, "dev-only-secret-change-me"),
		JWTTTLMinutes:     ttl,
		OracleUser:        defaultString(file.OracleUser, "metartls"),
		OraclePass:        defaultString(file.OraclePassword, "MetaRTLS_Dev_123"),
		OracleDSNHost:     defaultString(file.OracleDSN, "localhost:1521/FREEPDB1"),
		OracleSysPassword: defaultString(file.OracleSysPassword, "Oracle_Dev_123"),
		RedisAddr:         defaultString(file.RedisAddr, "localhost:6379"),
		RedisPassword:     file.RedisPassword,
		MQTTBroker:        defaultString(file.MQTTBroker, "tcp://localhost:1883"),
		MQTTClientID:      defaultString(file.MQTTClientID, "metartls-api"),
		MQTTTopic:         defaultString(file.MQTTTopicLocation, "rtls/+/location"),
		CORSOrigins:       origins,
	}

	if err := cfg.Validate(); err != nil {
		return nil, err
	}
	return cfg, nil
}

func resolveConfigPath() (string, error) {
	if path := strings.TrimSpace(os.Getenv("CONFIG_PATH")); path != "" {
		return path, nil
	}

	candidates := []string{
		"config/config.env",
		"../config/config.env",
		"../../config/config.env",
		"config/config-temp.env",
		"../config/config-temp.env",
		"../../config/config-temp.env",
	}
	for _, candidate := range candidates {
		if _, err := os.Stat(candidate); err == nil {
			abs, absErr := filepath.Abs(candidate)
			if absErr == nil {
				return abs, nil
			}
			return candidate, nil
		}
	}
	return "", fmt.Errorf("config file not found; copy config/config-temp.env to config/config.env")
}

func (c *Config) Validate() error {
	if c.JWTSecret == "" {
		return fmt.Errorf("jwtSecret is required")
	}
	if c.AppEnv != "production" {
		return nil
	}
	if len(c.JWTSecret) < 32 {
		return fmt.Errorf("jwtSecret must be at least 32 characters in production")
	}
	switch c.JWTSecret {
	case "change-me-in-production-use-long-random-string", "dev-only-secret-change-me":
		return fmt.Errorf("jwtSecret must not use the default value in production")
	}
	return nil
}

func (c *Config) OracleDSN() string {
	return fmt.Sprintf("oracle://%s:%s@%s", c.OracleUser, c.OraclePass, c.OracleDSNHost)
}

func defaultString(value, fallback string) string {
	if strings.TrimSpace(value) == "" {
		return fallback
	}
	return value
}
