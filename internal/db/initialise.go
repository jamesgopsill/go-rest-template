package db

import (
	"github.com/rs/zerolog/log"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

var Connection *gorm.DB

type User struct {
	gorm.Model
	Name         string
	Email        string
	PasswordHash []byte
}

func Initialise(dbPath string) {
	var err error
	log.Info().Msg("Connecting to database")
	Connection, err = gorm.Open(sqlite.Open(dbPath), &gorm.Config{})
	if err != nil {
		panic("failed to connect database")
	}
	log.Info().Msg("Connected to database")
	Connection.AutoMigrate(&User{})
}
