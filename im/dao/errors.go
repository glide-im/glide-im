package dao

import (
	"errors"
	"github.com/jinzhu/gorm"
)

func resolveError(db *gorm.DB) error {
	if db.RowsAffected == 0 {
		return errors.New("update failed, no such record")
	}
	if db.Error != nil {
		return db.Error
	}
	return nil
}
