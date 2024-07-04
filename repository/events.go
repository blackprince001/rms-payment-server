package repository

import (
	"cashapp/models"

	"gorm.io/gorm"
)

type eventLayer struct {
	db *gorm.DB
}

func newEventLayer(db *gorm.DB) *eventLayer {
	return &eventLayer{
		db: db,
	}
}

func (el *eventLayer) Save(tx *gorm.DB, data *models.TransactionEvent) error {
	if err := tx.Create(data).Error; err != nil {
		return err
	}
	return nil
}
