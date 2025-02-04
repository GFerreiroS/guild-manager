package database

import (
	"fmt"
	"log"
	"time"

	"github.com/GFerreiroS/guild-manager/backend/internal/models"
	"gorm.io/gorm"
)

// SeedTestData now returns an error.
func SeedTestData(db *gorm.DB) error {
	// Clean existing data: drop and recreate schema.
	if err := db.Exec("DROP SCHEMA public CASCADE; CREATE SCHEMA public;").Error; err != nil {
		return fmt.Errorf("failed to clean database: %w", err)
	}

	// Re-run migrations to re-create all tables.
	if err := RunMigrations(db); err != nil {
		return fmt.Errorf("failed to re-run migrations after schema drop: %w", err)
	}

	// Create Users
	users := []models.User{
		{
			BattleNetID: "testadmin#1234",
			Username:    "AdminUser",
			Email:       "admin@example.com",
			Role:        "admin",
		},
		{
			BattleNetID: "testofficer#5678",
			Username:    "OfficerUser",
			Email:       "officer@example.com",
			Role:        "officer",
		},
		{
			BattleNetID: "testmember#9012",
			Username:    "RegularMember",
			Email:       "member@example.com",
			Role:        "member",
		},
	}

	if err := db.Create(&users).Error; err != nil {
		return fmt.Errorf("failed to seed users: %w", err)
	}
	log.Println("Seeded 3 users")

	// Create Guilds
	guilds := []models.Guild{
		{
			Name:      "Alliance Elite",
			Realm:     "Stormrage",
			Faction:   "alliance",
			CreatedBy: users[0].ID,
		},
		{
			Name:      "Horde Champions",
			Realm:     "Illidan",
			Faction:   "horde",
			CreatedBy: users[1].ID,
		},
	}

	if err := db.Create(&guilds).Error; err != nil {
		return fmt.Errorf("failed to seed guilds: %w", err)
	}
	log.Println("Seeded 2 guilds")

	// Create Guild Members
	for i := range guilds {
		guilds[i].Members = []models.User{users[0], users[1], users[2]}
		if err := db.Model(&guilds[i]).Association("Members").Append(guilds[i].Members); err != nil {
			return fmt.Errorf("failed to seed guild members: %w", err)
		}
	}
	log.Println("Added members to guilds")

	// Create Characters
	characters := []models.Character{
		{
			Name:    "FireMage",
			Realm:   "Stormrage",
			Class:   "mage",
			Spec:    "Fire",
			Ilvl:    435,
			UserID:  users[0].ID,
			GuildID: guilds[0].ID,
		},
		{
			Name:    "HolyPally",
			Realm:   "Stormrage",
			Class:   "paladin",
			Spec:    "Holy",
			Ilvl:    430,
			UserID:  users[1].ID,
			GuildID: guilds[0].ID,
		},
		{
			Name:    "ShadowPriest",
			Realm:   "Illidan",
			Class:   "priest",
			Spec:    "Shadow",
			Ilvl:    428,
			UserID:  users[2].ID,
			GuildID: guilds[1].ID,
		},
	}

	if err := db.Create(&characters).Error; err != nil {
		return fmt.Errorf("failed to seed characters: %w", err)
	}
	log.Println("Seeded 3 characters")

	// Create Raid Groups
	raidGroups := []models.RaidGroup{
		{
			Name:     "Main Raid Team",
			GuildID:  guilds[0].ID,
			Schedule: models.JSONB{"days": []string{"Tuesday", "Thursday"}, "time": "20:00"},
		},
		{
			Name:     "Weekend Warriors",
			GuildID:  guilds[1].ID,
			Schedule: models.JSONB{"days": []string{"Saturday"}, "time": "15:00"},
		},
	}

	if err := db.Create(&raidGroups).Error; err != nil {
		return fmt.Errorf("failed to seed raid groups: %w", err)
	}
	log.Println("Seeded 2 raid groups")

	// Add Characters to Raid Groups
	raidGroups[0].Characters = []models.Character{characters[0], characters[1]}
	raidGroups[1].Characters = []models.Character{characters[2]}
	for _, rg := range raidGroups {
		if err := db.Model(&rg).Association("Characters").Append(rg.Characters); err != nil {
			return fmt.Errorf("failed to add characters to raid group: %w", err)
		}
	}
	log.Println("Added characters to raid groups")

	// Create Events
	events := []models.Event{
		{
			RaidName:    "Castle Nathria",
			Difficulty:  "heroic",
			ScheduledAt: time.Now().Add(24 * time.Hour),
			CreatedBy:   users[0].ID,
			GuildID:     guilds[0].ID,
		},
		{
			RaidName:    "Sanctum of Domination",
			Difficulty:  "mythic",
			ScheduledAt: time.Now().Add(48 * time.Hour),
			CreatedBy:   users[1].ID,
			GuildID:     guilds[1].ID,
		},
	}

	if err := db.Create(&events).Error; err != nil {
		return fmt.Errorf("failed to seed events: %w", err)
	}
	log.Println("Seeded 2 events")

	// Create Confirmations
	confirmations := []models.Confirmation{
		{
			EventID:     events[0].ID,
			CharacterID: characters[0].ID,
			Status:      "confirmed",
		},
		{
			EventID:     events[0].ID,
			CharacterID: characters[1].ID,
			Status:      "tentative",
			Reason:      "Might be late",
		},
		{
			EventID:     events[1].ID,
			CharacterID: characters[2].ID,
			Status:      "declined",
			Reason:      "Out of town",
		},
	}

	if err := db.Create(&confirmations).Error; err != nil {
		return fmt.Errorf("failed to seed confirmations: %w", err)
	}
	log.Println("Seeded 3 confirmations")

	log.Println("✅ All test data seeded successfully!")
	return nil
}
