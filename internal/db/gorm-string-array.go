package db

import (
	"database/sql/driver"
	"encoding/json"
	"fmt"
	"strings"
)

// https://stackoverflow.com/questions/41375563/unsupported-scan-storing-driver-value-type-uint8-into-type-string

type GormStringArray []string

func (s GormStringArray) Value() (driver.Value, error) {
	if len(s) == 0 {
		return "[]", nil
	}
	return fmt.Sprintf(`["%s"]`, strings.Join(s, `","`)), nil
}

func (s *GormStringArray) Scan(src interface{}) (err error) {
	var array []string
	err = json.Unmarshal([]byte(src.(string)), &array)
	if err != nil {
		return
	}
	*s = array
	return nil
}
