package db

import (
	"database/sql/driver"
	"encoding/json"
	"fmt"
	"strings"

	"github.com/rs/zerolog/log"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

var Connection *gorm.DB

type Base struct {
	gorm.Model
	ID string
}

type User struct {
	ID           string
	Name         string
	Email        string
	Scopes       Scopes
	PasswordHash string
	Thumbnail    string
}

// https://stackoverflow.com/questions/41375563/unsupported-scan-storing-driver-value-type-uint8-into-type-string

type Scopes []string

func (s Scopes) Value() (driver.Value, error) {
	if len(s) == 0 {
		return "[]", nil
	}
	return fmt.Sprintf(`["%s"]`, strings.Join(s, `","`)), nil
}

func (s *Scopes) Scan(src interface{}) (err error) {
	var scopes []string
	err = json.Unmarshal([]byte(src.(string)), &scopes)
	if err != nil {
		return
	}
	*s = scopes
	return nil
}

const (
	SYS_ADMIN_SCOPE = "sysadmin"
	ADMIN_SCOPE     = "admin"
	USER_SCOPE      = "user"
)

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
