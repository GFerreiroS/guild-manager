package models

import (
	"time"
)

type User struct {
	ID          string    `gorm:"type:uuid;primaryKey;default:uuid_generate_v4()"`
	BattleNetID string    `gorm:"type:varchar(255);uniqueIndex;not null"`
	Username    string    `gorm:"type:varchar(255);not null"`
	Email       string    `gorm:"type:varchar(255);unique"`
	Role        string    `gorm:"type:varchar(50);not null;default:'member'"`
	CreatedAt   time.Time `gorm:"autoCreateTime"`
	UpdatedAt   time.Time `gorm:"autoUpdateTime"`

	Characters    []Character `gorm:"foreignKey:UserID"`
	CreatedGuilds []Guild     `gorm:"foreignKey:CreatedBy"`
	Guilds        []Guild     `gorm:"many2many:guild_members;"`
}
