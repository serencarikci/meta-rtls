package config

import (
	"fmt"
	"os"
	"strconv"
	"strings"

	"github.com/joho/godotenv"
)

type Config struct {
	AppEnv        string
	AppPort       string
	JWTSecret     string
	JWTTTLMinutes int
	OracleUser    string
	OraclePass    string
	OracleDSNHost string
	RedisAddr     string
	RedisPassword string
	MQTTBroker    string
	MQTTClientID  string
	MQTTTopic     string
	CORSOrigins   []string
}

func Load() (*Config, error) {
	_ = godotenv.Load("../.env")
	_ = godotenv.Load(".env")

	ttl, err := strconv.Atoi(getEnv("JWT_TTL_MINUTES", "480"))
	if err != nil {
		return nil, fmt.Errorf("JWT_TTL_MINUTES: %w", err)
	}

	cfg := &Config{
		AppEnv:        getEnv("APP_ENV", "development"),
		AppPort:       getEnv("APP_PORT", "8080"),
		JWTSecret:     getEnv("JWT_SECRET", "dev-only-secret-change-me"),
		JWTTTLMinutes: ttl,
		OracleUser:    getEnv("ORACLE_USER", "metartls"),
		OraclePass:    getEnv("ORACLE_PASSWORD", "MetaRTLS_Dev_123"),
		OracleDSNHost: getEnv("ORACLE_DSN", "localhost:1521/FREEPDB1"),
		RedisAddr:     getEnv("REDIS_ADDR", "localhost:6379"),
		RedisPassword: getEnv("REDIS_PASSWORD", ""),
		MQTTBroker:    getEnv("MQTT_BROKER", "tcp://localhost:1883"),
		MQTTClientID:  getEnv("MQTT_CLIENT_ID", "metartls-api"),
		MQTTTopic:     getEnv("MQTT_TOPIC_LOCATION", "rtls/+/location"),
		CORSOrigins:   strings.Split(getEnv("CORS_ORIGINS", "http://localhost:5173"), ","),
	}

	if cfg.JWTSecret == "" {
		return nil, fmt.Errorf("JWT_SECRET is required")
	}
	return cfg, nil
}

func (c *Config) OracleDSN() string {
	return fmt.Sprintf("oracle://%s:%s@%s", c.OracleUser, c.OraclePass, c.OracleDSNHost)
}

func getEnv(key, fallback string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return fallback
}
