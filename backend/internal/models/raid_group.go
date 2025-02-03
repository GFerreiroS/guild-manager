package models

import (
	"database/sql/driver"
	"encoding/json"
	"errors"
	"time"
)

// Add this custom JSONB type at the top of the file
type JSONB map[string]interface{}

func (j JSONB) Value() (driver.Value, error) {
	return json.Marshal(j)
}

func (j *JSONB) Scan(value interface{}) error {
	b, ok := value.([]byte)
	if !ok {
		return errors.New("type assertion to []byte failed")
	}
	return json.Unmarshal(b, &j)
}

// Update RaidGroup struct
type RaidGroup struct {
	ID        string    `gorm:"type:uuid;primaryKey;default:uuid_generate_v4()"`
	Name      string    `gorm:"type:varchar(255);not null"`
	GuildID   string    `gorm:"type:uuid;index"`
	Schedule  JSONB     `gorm:"type:jsonb"` // Use custom type here
	CreatedAt time.Time `gorm:"autoCreateTime"`

	Guild      Guild       `gorm:"foreignKey:GuildID"`
	Characters []Character `gorm:"many2many:raid_group_characters;"`
}
