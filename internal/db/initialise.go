package db

import (
	"jamesgopsill/go-rest-template/internal/config"

	"github.com/rs/zerolog/log"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

var Connection *gorm.DB

func Initialise() {
	var err error
	log.Info().Msg("Connecting to database")
	Connection, err = gorm.Open(sqlite.Open(config.DBPath), &gorm.Config{})
	if err != nil {
		panic("failed to connect database")
	}
	log.Info().Msg("Connected to database")
	Connection.AutoMigrate(&User{})
}
