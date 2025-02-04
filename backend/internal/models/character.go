package models

import (
	"time"
)

type Character struct {
	ID          string    `gorm:"type:uuid;primaryKey;default:uuid_generate_v4()"`
	Name        string    `gorm:"type:varchar(255);not null"`
	Realm       string    `gorm:"type:varchar(255);not null"`
	Class       string    `gorm:"type:varchar(50);not null;check:class IN ('warrior','paladin','hunter','rogue','priest','death-knight','shaman','mage','warlock','monk','druid','demon-hunter','evoker')"`
	Spec        string    `gorm:"type:varchar(50)"`
	Ilvl        int       `gorm:"not null"`
	LastSynced  time.Time `gorm:"type:timestamptz"`
	UserID      string    `gorm:"type:uuid;not null"`
	GuildID     string    `gorm:"type:uuid;not null"`
	RaidGroupID *string   `gorm:"type:uuid"` // Changed to pointer so that nil (NULL) is stored if not set
	CreatedAt   time.Time `gorm:"autoCreateTime"`
	UpdatedAt   time.Time `gorm:"autoUpdateTime"`

	User          User           `gorm:"foreignKey:UserID"`
	Guild         Guild          `gorm:"foreignKey:GuildID"`
	RaidGroup     RaidGroup      `gorm:"foreignKey:RaidGroupID"`
	Confirmations []Confirmation `gorm:"foreignKey:CharacterID"`
}
