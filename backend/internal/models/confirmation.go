package models

import (
	"time"
)

type Confirmation struct {
	ID          string    `gorm:"type:uuid;primaryKey;default:uuid_generate_v4()"`
	EventID     string    `gorm:"type:uuid;index"`
	CharacterID string    `gorm:"type:uuid;index"`
	Status      string    `gorm:"type:varchar(50);check:status IN ('confirmed','declined','tentative');default:'pending'"`
	Reason      string    `gorm:"type:text"`
	RespondedAt time.Time `gorm:"autoCreateTime"`

	Event     Event     `gorm:"foreignKey:EventID"`
	Character Character `gorm:"foreignKey:CharacterID"`
}
