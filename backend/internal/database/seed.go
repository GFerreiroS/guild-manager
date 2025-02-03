package database

import (
	"log"

	"github.com/GFerreiroS/guild-manager/backend/internal/models"

	"gorm.io/gorm"
)

func SeedTestData(db *gorm.DB) {
	// Create test user
	user := models.User{
		BattleNetID: "testuser#1234",
		Username:    "TestUser",
		Email:       "test@example.com",
		Role:        "admin",
	}

	if err := db.Create(&user).Error; err != nil {
		log.Fatalf("Failed to seed user: %v", err)
	}

	// Create test guild
	guild := models.Guild{
		Name:      "Test Guild",
		Realm:     "Stormrage",
		Faction:   "alliance",
		CreatedBy: user.ID,
	}

	if err := db.Create(&guild).Error; err != nil {
		log.Fatalf("Failed to seed guild: %v", err)
	}

	// Create test character
	character := models.Character{
		Name:    "Testchar",
		Class:   "mage",
		Ilvl:    425,
		UserID:  user.ID,
		GuildID: guild.ID,
	}

	if err := db.Create(&character).Error; err != nil {
		log.Fatalf("Failed to seed character: %v", err)
	}

	log.Println("Test data seeded successfully")
}
