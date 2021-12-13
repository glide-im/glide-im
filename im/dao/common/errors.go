package common

import (
	"errors"
	"gorm.io/gorm"
)

var ErrNoRecordFound = errors.New("no record found")
var ErrNoneUpdated = errors.New("no record updated, RowsAffected=0")

func ResolveError(db *gorm.DB) error {
	if db.Error != nil {
		return db.Error
	}
	if db.RowsAffected == 0 {
		return errors.New("RowsAffected=0")
	}
	return nil
}

func JustError(db *gorm.DB) error {
	if db.Error != nil {
		return db.Error
	}
	return nil
}

func MustFind(db *gorm.DB) error {
	if db.Error != nil {
		return db.Error
	}
	if db.RowsAffected == 0 {
		return ErrNoRecordFound
	}
	return nil
}

func MustUpdate(db *gorm.DB) error {
	if db.Error != nil {
		return db.Error
	}
	if db.RowsAffected == 0 {
		return ErrNoneUpdated
	}
	return nil
}

func ResolveUpdateErr(db *gorm.DB) error {
	if db.Error != nil {
		return db.Error
	}
	return nil
}
