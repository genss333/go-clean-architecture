package config

import (
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	Server   ServerConfig
	Postgres DatabaseConfig
	MariaDB  DatabaseConfig
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
			User:     getEnv("PG_USER", "benchmark"),
			Password: getEnv("PG_PASSWORD", "benchmark"),
			DBName:   getEnv("PG_DBNAME", "benchmark_db"),
			SSLMode:  getEnv("PG_SSLMODE", "disable"),
		},
		MariaDB: DatabaseConfig{
			Host:     getEnv("MARIA_HOST", "localhost"),
			Port:     getEnv("MARIA_PORT", "3306"),
			User:     getEnv("MARIA_USER", "benchmark"),
			Password: getEnv("MARIA_PASSWORD", "benchmark"),
			DBName:   getEnv("MARIA_DBNAME", "benchmark_db"),
		},
		RabbitMQ: RabbitMQConfig{
			URL: getEnv("RABBITMQ_URL", "amqp://guest:guest@localhost:5672/"),
		},
	}, nil
}

func getEnv(key, fallback string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return fallback
}
