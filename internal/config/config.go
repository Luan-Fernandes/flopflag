package config

import (
	"cmp"
	"errors"
	"os"
	"strconv"

	"github.com/joho/godotenv"
)

type DBConfig struct {
	Host     string
	Port     string
	User     string
	Password string
	Name     string
	SSLMode  string
}

type JWTConfig struct {
	Secret         string
	ExpiresMinutes int
}

type Config struct {
	AppEnv   string
	HTTPPort string
	DB       DBConfig
	JWT      JWTConfig
}

func Load() (Config, error) {
	_ = godotenv.Load()

	cfg := Config{
		AppEnv:   getEnv("APP_ENV", "development"),
		HTTPPort: getEnv("HTTP_PORT", "8080"),
		DB: DBConfig{
			Host:     getEnv("DB_HOST", "localhost"),
			Port:     getEnv("DB_PORT", "5432"),
			User:     getEnv("DB_USER", "boli"),
			Password: getEnv("DB_PASSWORD", "boli"),
			Name:     getEnv("DB_NAME", "boli"),
			SSLMode:  getEnv("DB_SSLMODE", "disable"),
		},
		JWT: JWTConfig{
			Secret: getEnv("JWT_SECRET", ""),
		},
	}

	expMin, err := strconv.Atoi(getEnv("JWT_EXPIRES_MINUTES", "60"))
	if err != nil || expMin <= 0 {
		return Config{}, errors.New("JWT_EXPIRES_MINUTES inválido")
	}
	cfg.JWT.ExpiresMinutes = expMin

	return cfg, nil
}

func (c Config) ValidateForAPI() error {
	if c.JWT.Secret == "" {
		return errors.New("JWT_SECRET é obrigatório")
	}
	return nil
}

func getEnv(key, def string) string {
	return cmp.Or(os.Getenv(key), def)
}
