package database

import (
	"fmt"
	"log"
	"os"
	"time"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"

	"github.com/GFerreiroS/guild-manager/backend/internal/models"
)

type PostgresConfig struct {
	Host     string `mapstructure:"POSTGRES_HOST"`
	Port     string `mapstructure:"POSTGRES_PORT"`
	User     string `mapstructure:"POSTGRES_USER"`
	Password string `mapstructure:"POSTGRES_PASSWORD"`
	DBName   string `mapstructure:"POSTGRES_DB"`
	SSLMode  string `mapstructure:"POSTGRES_SSL_MODE"`
}

func NewPostgresDB() (*gorm.DB, error) {
	// Load configuration from environment variables
	cfg := PostgresConfig{
		Host:     getEnv("POSTGRES_HOST", "postgres"),
		Port:     getEnv("POSTGRES_PORT", "5432"),
		User:     getEnv("POSTGRES_USER", ""),
		Password: getEnv("POSTGRES_PASSWORD", ""),
		DBName:   getEnv("POSTGRES_DB", "guild_manager"),
		SSLMode:  getEnv("POSTGRES_SSL_MODE", "disable"),
	}

	// Validate required fields
	if cfg.User == "" || cfg.Password == "" {
		return nil, fmt.Errorf("database credentials are required")
	}

	// Create DSN
	dsn := fmt.Sprintf(
		"host=%s port=%s user=%s password=%s dbname=%s sslmode=%s",
		cfg.Host, cfg.Port, cfg.User, cfg.Password, cfg.DBName, cfg.SSLMode,
	)

	// Configure GORM
	gormConfig := &gorm.Config{
		Logger: logger.Default.LogMode(logger.Info),
		NowFunc: func() time.Time {
			return time.Now().UTC()
		},
	}

	// Connect to database
	db, err := gorm.Open(postgres.Open(dsn), gormConfig)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to database: %w", err)
	}

	// Enable UUID extension
	if err := db.Exec("CREATE EXTENSION IF NOT EXISTS \"uuid-ossp\"").Error; err != nil {
		return nil, fmt.Errorf("failed to create uuid extension: %w", err)
	}

	log.Println("Successfully connected to PostgreSQL database")
	return db, nil
}

func getEnv(key, defaultValue string) string {
	value := os.Getenv(key)
	if value == "" {
		return defaultValue
	}
	return value
}

func AutoMigrate(db *gorm.DB) error {
	models := []interface{}{
		&models.User{},
		&models.Guild{},
		&models.Character{},
		&models.RaidGroup{},
		&models.Event{},
		&models.Confirmation{},
		&models.GuildMember{},
		&models.RaidGroupCharacter{},
	}

	if err := db.AutoMigrate(models...); err != nil {
		return fmt.Errorf("failed to auto-migrate models: %w", err)
	}

	// Manually create indexes that GORM can't handle
	if err := db.Exec(`
		CREATE INDEX IF NOT EXISTS idx_characters_guild ON characters(guild_id);
		CREATE INDEX IF NOT EXISTS idx_events_guild ON events(guild_id);
		CREATE INDEX IF NOT EXISTS idx_confirmations_event ON confirmations(event_id);
	`).Error; err != nil {
		return fmt.Errorf("failed to create indexes: %w", err)
	}

	return nil
}
