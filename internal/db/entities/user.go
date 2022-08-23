package entities

import "time"

type User struct {
	ID           string `gorm:"primaryKey"`
	CreatedAt    time.Time
	UpdatedAt    time.Time
	Name         string
	Email        string `gorm:"uniqueIndex"`
	Scopes       SerialisableStringArray
	PasswordHash string `json:"-"`
	Thumbnail    string
}
