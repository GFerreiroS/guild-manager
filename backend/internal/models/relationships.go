package models

import "time"

// GuildMember represents the many-to-many relationship between users and guilds
type GuildMember struct {
	UserID   string `gorm:"type:uuid;primaryKey"`
	GuildID  string `gorm:"type:uuid;primaryKey"`
	JoinedAt time.Time
	Role     string `gorm:"type:varchar(50)"`

	User  User  `gorm:"foreignKey:UserID"`
	Guild Guild `gorm:"foreignKey:GuildID"`
}

// RaidGroupCharacter represents the many-to-many relationship between characters and raid groups
type RaidGroupCharacter struct {
	CharacterID string `gorm:"type:uuid;primaryKey"`
	RaidGroupID string `gorm:"type:uuid;primaryKey"`
	JoinedAt    time.Time

	Character Character `gorm:"foreignKey:CharacterID"`
	RaidGroup RaidGroup `gorm:"foreignKey:RaidGroupID"`
}
