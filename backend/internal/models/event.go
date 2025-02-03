package models

import (
	"time"
)

type Event struct {
	ID          string    `gorm:"type:uuid;primaryKey;default:uuid_generate_v4()"`
	RaidName    string    `gorm:"type:varchar(255);not null"`
	Difficulty  string    `gorm:"type:varchar(50);check:difficulty IN ('normal','heroic','mythic')"`
	ScheduledAt time.Time `gorm:"type:timestamptz"`
	CreatedBy   string    `gorm:"type:uuid;index"`
	GuildID     string    `gorm:"type:uuid;index"`
	CreatedAt   time.Time `gorm:"autoCreateTime"`

	Creator       User           `gorm:"foreignKey:CreatedBy"`
	Guild         Guild          `gorm:"foreignKey:GuildID"`
	Confirmations []Confirmation `gorm:"foreignKey:EventID"`
}
