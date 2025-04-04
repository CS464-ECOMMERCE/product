package storage

import (
	"product/models"

	"gorm.io/gorm"
)

type OrderInterface interface {
	CreateOrder(order *models.Order, tx *gorm.DB) (*models.Order, error)
	UpdateOrder(order *models.Order, tx *gorm.DB) error
}

type OrderDB struct {
	write *gorm.DB
}

func NewOrderTable(write *gorm.DB) OrderInterface {
	StorageInstance.AutoMigrate(&models.Order{})
	return &OrderDB{
		write: write,
	}
}

func (i *OrderDB) CreateOrder(order *models.Order, tx *gorm.DB) (*models.Order, error) {
	// Use the provided transaction if available, otherwise use the default write DB
	db := tx
	if db == nil {
		db = i.write
	}

	// Perform the create operation
	ret := db.Create(order)
	if ret.Error != nil {
		return nil, ret.Error
	}
	return order, nil
}

func (i *OrderDB) UpdateOrder(order *models.Order, tx *gorm.DB) error {
	db := tx
	if db == nil {
		db = i.write
	}

	// Perform the create operation
	ret := db.Model(&models.Order{}).Where("id = ?", order.Id).Updates(&order)
	if ret.Error != nil {
		return ret.Error
	}
	return nil
}
