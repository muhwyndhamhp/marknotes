package post

import (
	"github.com/muhwyndhamhp/marknotes/pkg/models"
	"gorm.io/gorm"
)

func Upsert(p *models.Post, db *gorm.DB) error {
	if trxErr := db.Transaction(func(tx *gorm.DB) error {
		if err := db.
			Save(p).Error; err != nil {
			return err
		}

		if len(p.Tags) <= 0 {
			return nil
		}

		if err := db.
			Model(p).
			Association("Tags").
			Replace(p.Tags); err != nil {
			return err
		}
		return nil
	}); trxErr != nil {
		return trxErr
	}

	return nil
}
