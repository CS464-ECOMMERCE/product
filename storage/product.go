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
	List(limit uint64, cursorID uint64) ([]*pb.Product, uint64, uint64, error)
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

func (i *ProductDB) List(limit uint64, cursorID uint64) ([]*pb.Product, uint64, uint64, error) {
	var products []*pb.Product

	query := i.read.Order("id ASC").Limit(int(limit))
	// Count the total number of products
	var totalProducts int64
	if err := i.read.Model(&pb.Product{}).Distinct("id").Count(&totalProducts).Error; err != nil {
		return nil, 0, 0, err
	}

	// Apply cursor condition if provided
	if cursorID > 0 {
		query = query.Where("id > ?", cursorID)
	}

	// Fetch the products
	if err := query.Find(&products).Error; err != nil {
		return nil, 0, 0, err
	}

	// Get the last product's ID as the next cursor
	var nextCursor uint64
	if len(products) > 0 {
		nextCursor = products[len(products)-1].Id
	}

	return products, nextCursor, uint64(totalProducts), nil
}

func NewProductTable(read, write *gorm.DB) ProductInterface {
	StorageInstance.AutoMigrate(&pb.Product{})
	return &ProductDB{
		read:  read,
		write: write,
	}
}
