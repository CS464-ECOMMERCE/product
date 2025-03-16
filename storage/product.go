package storage

import (
	"errors"
	pb "product/proto"

	"gorm.io/gorm"
)

type ProductInterface interface {
	CreateProduct(Product *pb.Product) (*pb.Product, error)
	Get(id uint64) (*pb.Product, error)
	Update(Product *pb.Product) (*pb.Product, error)
	Delete(id uint64) error
}

type ProductDB struct {
	read  *gorm.DB
	write *gorm.DB
}

// CreateProduct implements ProductInterface.
func (i *ProductDB) CreateProduct(Product *pb.Product) (*pb.Product, error) {
	ret := i.write.Create(Product)
	if ret.Error != nil {
		return nil, ret.Error
	}
	return Product, nil
}

// Delete implements ProductInterface.
func (i *ProductDB) Delete(id uint64) error {
	return i.write.Delete(&pb.Product{}, "id = ?", id).Error
}

// Get implements ProductInterface.
func (i *ProductDB) Get(id uint64) (*pb.Product, error) {
	Product := &pb.Product{}
	ret := i.read.Where("id = ?", id).First(Product)
	if ret.Error != nil {
		return nil, ret.Error
	}
	return Product, nil
}

// Update implements ProductInterface.
func (i *ProductDB) Update(Product *pb.Product) (*pb.Product, error) {
	result := i.write.Model(&pb.Product{}).Where("id = ?", Product.Id).Where("merchant_id = ?", Product.MerchantId).Updates(Product)

	if result.Error != nil {
		return nil, result.Error // Return the actual error
	}

	if result.RowsAffected == 0 {
		return nil, errors.New("no product found with the given ID")
	}

	return Product, nil
}

func NewProductTable(read, write *gorm.DB) ProductInterface {
	StorageInstance.AutoMigrate(&pb.Product{})
	return &ProductDB{
		read:  read,
		write: write,
	}
}
