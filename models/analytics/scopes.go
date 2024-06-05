package analytics

import (
	"fmt"

	"gorm.io/gorm"
)

func GetLatest(slug string) func(*gorm.DB) *gorm.DB {
	path := fmt.Sprintf("/articles/%s.html", slug)
	return func(d *gorm.DB) *gorm.DB {
		return d.
			Where("path = ?", path).
			Order("capture_date DESC").
			Limit(1)
	}
}
