package config

import (
	"fmt"
	"time"

	"github.com/spf13/viper"
)

// Config holds all configuration for the application.
type Config struct {
	// PostgreSQL settings
	Database struct {
		Host     string
		Port     string
		User     string
		Password string
		DBName   string
		SSLMode  string
	}
	// Redis settings
	Redis struct {
		Addr     string
		Password string
		DB       int
		Timeout  time.Duration
	}
	// Rate limiting settings
	RateLimit struct {
		RequestsPerMinute int
	}
}

// LoadConfig loads configuration using Viper.
func LoadConfig() (*Config, error) {
	// Set environment variable prefix if you like (optional)
	// viper.SetEnvPrefix("APP")

	viper.AutomaticEnv()

	// Set defaults for Database
	viper.SetDefault("database.host", "postgres")
	viper.SetDefault("database.port", "5432")
	viper.SetDefault("database.user", "admin")
	viper.SetDefault("database.password", "secret")
	viper.SetDefault("database.dbname", "guild_manager")
	viper.SetDefault("database.sslmode", "disable")

	// Set defaults for Redis
	viper.SetDefault("redis.addr", "redis:6379")
	viper.SetDefault("redis.password", "")
	viper.SetDefault("redis.db", 0)
	viper.SetDefault("redis.timeout", 5) // seconds

	// Set defaults for Rate Limiting
	viper.SetDefault("ratelimit.requestspersminute", 60)

	var cfg Config
	if err := viper.Unmarshal(&cfg); err != nil {
		return nil, fmt.Errorf("failed to unmarshal configuration: %w", err)
	}

	// Convert timeout to time.Duration (if needed)
	cfg.Redis.Timeout = cfg.Redis.Timeout * time.Second

	return &cfg, nil
}
