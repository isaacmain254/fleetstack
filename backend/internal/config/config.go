package config

import (
	"log"
	"os"
	"strconv"
	"time"

	"github.com/joho/godotenv"
)

type Config struct {
	GitHub   GitHubConfig
	Database DatabaseConfig
	Server   ServerConfig
	Docker   DockerConfig
}

type GitHubConfig struct {
	ClientID     string
	ClientSecret string
}

type DatabaseConfig struct {
	Host     string
	Port     int
	User     string
	Password string
	DBName     string
	SSLMode  string
	MaxOpenConns    int
	MaxIdleConns    int
	ConnMaxLifetime time.Duration
}

type ServerConfig struct {
	Port string
}

type DockerConfig struct {
	Socket string
}

func Load() *Config {
	if err := godotenv.Load(); err != nil {
		log.Println(".env not found, using environment variables")
	}

	// dbPort := getEnvAsInt("DB_PORT", 5432)

	return &Config{
		GitHub: GitHubConfig{
			ClientID:     getEnv("CLIENT_ID", ""),
			ClientSecret: getEnv("CLIENT_SECRET", ""),
		},

		Database: DatabaseConfig{
			Host:     getEnv("DB_HOST", "localhost"),
			Port:     getEnvAsInt("DB_PORT", 5432),
			User:     getEnv("DB_USER", "fleetstack_user"),
			Password: getEnv("DB_PASSWORD", "fleetstack_password"),
			DBName:     getEnv("DB_NAME", "fleetstack_platform"),
			SSLMode:  getEnv("DB_SSLMODE", "disable"),
		},

		Server: ServerConfig{
			Port: getEnv("PORT", "8080"),
		},

		Docker: DockerConfig{
			Socket: getEnv("DOCKER_HOST", "unix:///var/run/docker.sock"),
		},
	}
}

func getEnv(key, fallback string) string {
	if value, ok := os.LookupEnv(key); ok {
		return value
	}
	return fallback
}

func getEnvAsInt(key string, fallback int) int {
	value := getEnv(key, "")
	if value == "" {
		return fallback
	}

	n, err := strconv.Atoi(value)
	if err != nil {
		return fallback
	}

	return n
}