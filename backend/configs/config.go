package configs

import (
	"fmt"
	"os"
	"strconv"
)

// Config holds all configuration for the application
type Config struct {
	Server   ServerConfig
	Database DatabaseConfig
	Redis    RedisConfig
	JWT      JWTConfig
}

// ServerConfig holds server configuration
type ServerConfig struct {
	Port         string
	ReadTimeout  int
	WriteTimeout int
	Mode         string // debug, release, test
}

// DatabaseConfig holds database configuration
type DatabaseConfig struct {
	URL            string
	MaxConnections int
	MinConnections int
	MaxLifetime    int // in minutes
}

// RedisConfig holds Redis configuration
type RedisConfig struct {
	URL      string
	Addr     string
	Password string
	DB       int
}

// JWTConfig holds JWT configuration
type JWTConfig struct {
	Secret     string
	Expiration int // in hours
}

// Load loads configuration from environment variables
func Load() *Config {
	return &Config{
		Server: ServerConfig{
			Port:         getEnv("PORT", "8080"),
			ReadTimeout:  getEnvInt("SERVER_READ_TIMEOUT", 30),
			WriteTimeout: getEnvInt("SERVER_WRITE_TIMEOUT", 30),
			Mode:         getEnv("GIN_MODE", "release"),
		},
		Database: DatabaseConfig{
			URL:            getEnv("DATABASE_URL", "postgres://postgres:password@localhost:5432/ai_learning?sslmode=disable"),
			MaxConnections: getEnvInt("DB_MAX_CONNECTIONS", 25),
			MinConnections: getEnvInt("DB_MIN_CONNECTIONS", 5),
			MaxLifetime:    getEnvInt("DB_MAX_LIFETIME", 60),
		},
		Redis: RedisConfig{
			URL:      getEnv("REDIS_URL", "redis://localhost:6379"),
			Addr:     getEnv("REDIS_ADDR", "localhost:6379"),
			Password: getEnv("REDIS_PASSWORD", ""),
			DB:       getEnvInt("REDIS_DB", 0),
		},
		JWT: JWTConfig{
			Secret:     getEnv("JWT_SECRET", "your-secret-key-change-in-production"),
			Expiration: getEnvInt("JWT_EXPIRATION", 24),
		},
	}
}

// getEnv gets environment variable or returns default value
func getEnv(key, defaultValue string) string {
	value := os.Getenv(key)
	if value == "" {
		return defaultValue
	}
	return value
}

// getEnvInt gets environment variable as integer or returns default value
func getEnvInt(key string, defaultValue int) int {
	value := os.Getenv(key)
	if value == "" {
		return defaultValue
	}
	intValue, err := strconv.Atoi(value)
	if err != nil {
		return defaultValue
	}
	return intValue
}

// Validate validates the configuration
func (c *Config) Validate() error {
	if c.JWT.Secret == "your-secret-key-change-in-production" {
		fmt.Println("WARNING: Using default JWT secret. Change JWT_SECRET in production!")
	}
	
	if c.Server.Port == "" {
		return fmt.Errorf("server port cannot be empty")
	}
	
	if c.Database.URL == "" {
		return fmt.Errorf("database URL cannot be empty")
	}
	
	return nil
}
