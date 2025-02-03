package models

import (
	"time"
)

type Character struct {
	ID          string    `gorm:"type:uuid;primaryKey;default:uuid_generate_v4()"`
	Name        string    `gorm:"type:varchar(255);not null"`
	Class       string    `gorm:"type:varchar(50);check:class IN ('warrior','paladin','hunter','rogue','priest','death-knight','shaman','mage','warlock','monk','druid','demon-hunter','evoker')"`
	Spec        string    `gorm:"type:varchar(50)"`
	Ilvl        int       `gorm:"type:integer"`
	LastSynced  time.Time `gorm:"type:timestamptz"`
	UserID      string    `gorm:"type:uuid;index"`
	GuildID     string    `gorm:"type:uuid;index"`
	RaidGroupID string    `gorm:"type:uuid;index"`
	CreatedAt   time.Time `gorm:"autoCreateTime"`

	User          User           `gorm:"foreignKey:UserID"`
	Guild         Guild          `gorm:"foreignKey:GuildID"`
	RaidGroup     RaidGroup      `gorm:"foreignKey:RaidGroupID"`
	Confirmations []Confirmation `gorm:"foreignKey:CharacterID"`
}
