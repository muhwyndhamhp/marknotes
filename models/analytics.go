package models

import (
	"time"

	"github.com/muhwyndhamhp/marknotes/utils/typeext"
	"gorm.io/gorm"
)

type Analytics struct {
	gorm.Model
	CaptureDate time.Time
	Path        string `gorm:"index"`
	Data        typeext.JSONB
}
