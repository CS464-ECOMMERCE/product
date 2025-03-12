package storage

import (
	pb "product/proto"

	"gorm.io/gorm"
)

type ProductInterface interface {
	CreateProduct(Product *pb.Product) (*pb.Product, error)
	Get(id string) (*pb.Product, error)
	Update(Product *pb.Product) error
	Delete(id string) error
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
func (i *ProductDB) Delete(id string) error {
	return i.write.Delete(&pb.Product{}, id).Error
}

// Get implements ProductInterface.
func (i *ProductDB) Get(id string) (*pb.Product, error) {
	Product := &pb.Product{}
	ret := i.read.First(Product, id)
	if ret.Error != nil {
		return nil, ret.Error
	}
	return Product, nil
}

// Update implements ProductInterface.
func (i *ProductDB) Update(Product *pb.Product) error {
	return i.write.Save(Product).Error
}

func NewProductTable(read, write *gorm.DB) ProductInterface {
	StorageInstance.AutoMigrate(&pb.Product{})
	return &ProductDB{
		read:  read,
		write: write,
	}
}
