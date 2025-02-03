package main

import (
	"embed"
	"errors"
	"flag"
	"log"
	"strconv"

	"github.com/GFerreiroS/guild-manager/backend/internal/database"
	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/postgres"
	"github.com/golang-migrate/migrate/v4/source/iofs"
)

//go:embed migrations/*.sql
var fs embed.FS

func main() {
	var action string
	var version string
	flag.StringVar(&action, "action", "up", "Migration action (up/down/force)")
	flag.StringVar(&version, "version", "", "Version number for force action")
	flag.Parse()

	// Initialize database
	db, err := database.NewPostgresDB()
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}

	// Get generic SQL.DB instance
	sqlDB, err := db.DB()
	if err != nil {
		log.Fatalf("Failed to get SQL.DB: %v", err)
	}

	// Create migration driver
	driver, err := postgres.WithInstance(sqlDB, &postgres.Config{})
	if err != nil {
		log.Fatalf("Failed to create driver: %v", err)
	}

	// Create migration source
	source, err := iofs.New(fs, "migrations")
	if err != nil {
		log.Fatalf("Failed to create source: %v", err)
	}

	// Initialize migrator
	m, err := migrate.NewWithInstance("iofs", source, "postgres", driver)
	if err != nil {
		log.Fatalf("Failed to initialize migrator: %v", err)
	}

	// Handle migration action
	switch action {
	case "up":
		if err := m.Up(); err != nil {
			if dirtyErr, ok := err.(migrate.ErrDirty); ok {
				log.Printf("Dirty state detected at version %d", dirtyErr.Version)
				m.Force(int(dirtyErr.Version))
				log.Println("Retrying migration...")
				if err := m.Up(); err != nil {
					log.Fatal(err)
				}
			} else if err != migrate.ErrNoChange {
				log.Fatal(err)
			}
		}
		log.Println("Migrations applied successfully")
	case "down":
		if err := m.Down(); err != nil && !errors.Is(err, migrate.ErrNoChange) {
			log.Fatalf("Rollback failed: %v", err)
		}
		log.Println("Migrations rolled back successfully")
	case "force":
		ver, err := strconv.Atoi(version)
		if err != nil {
			log.Fatalf("Invalid version: %v", err)
		}
		if err := m.Force(ver); err != nil {
			log.Fatalf("Force failed: %v", err)
		}
		log.Printf("Forced version to %d", ver)
	default:
		log.Fatalf("Invalid action: %s", action)
	}
}
