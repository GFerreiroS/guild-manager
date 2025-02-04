package main

import (
	"flag"
	"log"
	"strconv"

	"github.com/GFerreiroS/guild-manager/backend/internal/database"
)

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

	// Handle migration action
	switch action {
	case "up":
		if err := database.RunMigrations(db); err != nil {
			log.Fatalf("Migration failed: %v", err)
		}
	case "force":
		ver, err := strconv.Atoi(version)
		if err != nil {
			log.Fatalf("Invalid version: %v", err)
		}
		if err := database.ForceVersion(db, ver); err != nil {
			log.Fatalf("Force version failed: %v", err)
		}
		log.Printf("Forced version to %d", ver)
	default:
		log.Fatalf("Invalid action: %s", action)
	}
}
