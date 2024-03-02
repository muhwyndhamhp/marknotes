package scopes

import (
	"fmt"

	"gorm.io/gorm"
)

type Direction string

const (
	Ascending  Direction = "asc"
	Descending Direction = "desc"
)

func OrderBy(field string, direction Direction) QueryScope {
	return func(db *gorm.DB) *gorm.DB {
		if field == "" || direction == "" {
			return db
		}
		return db.Order(fmt.Sprintf("%s %s", field, direction))
	}
}
