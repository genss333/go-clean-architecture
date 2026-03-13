package config

import (
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	Server   ServerConfig
	Postgres DatabaseConfig
	RabbitMQ RabbitMQConfig
}

type ServerConfig struct {
	Port string
	Mode string // debug, release, test
}

type DatabaseConfig struct {
	Host     string
	Port     string
	User     string
	Password string
	DBName   string
	SSLMode  string
}

type RabbitMQConfig struct {
	URL string
}

func Load() (*Config, error) {
	_ = godotenv.Load()

	return &Config{
		Server: ServerConfig{
			Port: getEnv("SERVER_PORT", "8080"),
			Mode: getEnv("GIN_MODE", "debug"),
		},
		Postgres: DatabaseConfig{
			Host:     getEnv("PG_HOST", "localhost"),
			Port:     getEnv("PG_PORT", "5432"),
			User:     getEnv("PG_USER", ""),
			Password: getEnv("PG_PASSWORD", ""),
			DBName:   getEnv("PG_DBNAME", ""),
			SSLMode:  getEnv("PG_SSLMODE", "disable"),
		},
		RabbitMQ: RabbitMQConfig{
			URL: getEnv("RABBITMQ_URL", ""),
		},
	}, nil
}

func getEnv(key, fallback string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return fallback
}
