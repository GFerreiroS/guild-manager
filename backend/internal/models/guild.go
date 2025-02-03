package models

import (
	"time"
)

type Guild struct {
	ID        string    `gorm:"type:uuid;primaryKey;default:uuid_generate_v4()"`
	Name      string    `gorm:"type:varchar(255);not null;uniqueIndex"`
	Realm     string    `gorm:"type:varchar(255);not null"`
	Faction   string    `gorm:"type:varchar(50);check:faction IN ('alliance','horde')"`
	CreatedBy string    `gorm:"type:uuid;index"`
	CreatedAt time.Time `gorm:"autoCreateTime"`

	Members    []User      `gorm:"many2many:guild_members;"`
	Characters []Character `gorm:"foreignKey:GuildID"`
	RaidGroups []RaidGroup `gorm:"foreignKey:GuildID"`
	Events     []Event     `gorm:"foreignKey:GuildID"`
}
