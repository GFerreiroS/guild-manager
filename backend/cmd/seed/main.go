package main

import (
	"log"

	"github.com/GFerreiroS/guild-manager/backend/internal/database"
)

func main() {
	db, err := database.NewPostgresDB()
	if err != nil {
		log.Fatal("Failed to connect to database:", err)
	}

	// Run migrations first
	if err := database.RunMigrations(db); err != nil {
		log.Fatal("Migrations failed:", err)
	}

	// Seed test data with error handling
	log.Println("🌱 Seeding test data...")
	if err := database.SeedTestData(db); err != nil {
		log.Fatal("Seeding failed:", err)
	}

	log.Println("✅ Database seeded successfully")
}
