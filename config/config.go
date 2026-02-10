package config

import (
	"fmt"
	"os"
	"strconv"

	"github.com/joho/godotenv"
)

type AuthConfig struct {
	Secret    string
	ExpiryMin int
}

type Config struct {
	Port            int
	DBUrl           string
	ExchangeRateAPI string
	AuthConfig      AuthConfig
}

func LoadConfig() (Config, error) {
	if err := godotenv.Load(); err != nil {
		return Config{}, fmt.Errorf("error in loading .env file: %w", err)
	}

	appPort, err := strconv.Atoi(getEnv("APP_PORT", "8080"))
	if err != nil {
		return Config{}, fmt.Errorf("invalid APP_PORT: %w", err)
	}

	expiryMin, err := strconv.Atoi(getEnv("AUTH_EXPIRY_MIN", "60"))
	if err != nil {
		return Config{}, fmt.Errorf("invalid AUTH_EXPIRY_MIN: %w", err)
	}

	cfg := Config{
		Port:  appPort,
		DBUrl: getEnv("DB_URL", ""),
		ExchangeRateAPI: getEnv("EXCHANGE_RATE_API", ""),
		AuthConfig: AuthConfig{
			Secret:    getEnv("AUTH_SECRET", ""),
			ExpiryMin: expiryMin,
		},
	}

	// required fiels
	if cfg.DBUrl == "" {
		return Config{}, fmt.Errorf("DB_URL must be set")
	}

	if cfg.AuthConfig.Secret == "" {
		return Config{}, fmt.Errorf("AUTH_SECRET must be set")
	}
	if cfg.ExchangeRateAPI == "" {
		return Config{}, fmt.Errorf("EXCHANGE_RATE_API must be set")
	}

	return cfg, nil
}

func getEnv(key, defaultVal string) string {
	if v, ok := os.LookupEnv(key); ok {
		return v
	}
	return defaultVal
}
