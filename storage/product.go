package storage

import (
	"errors"
	"product/models"
	pb "product/proto"

	"gorm.io/gorm"
)

type ProductInterface interface {
	CreateProduct(Product *models.Product) (*models.Product, error)
	Get(id uint64, tx *gorm.DB) (*models.Product, error)
	Update(Product *models.Product, tx *gorm.DB) (*models.Product, error)
	UpdateInventory(id, inventory uint64, tx *gorm.DB) error
	Delete(id uint64) error
	List(limit uint64, cursorID uint64) ([]*models.Product, uint64, uint64, error)
	ListByMerchantId(merchantId uint64, limit uint64, cursorID uint64) ([]*models.Product, uint64, uint64, error)
	UpdateImageUrl(Product *models.Product) (*models.Product, error)
}

type ProductDB struct {
	read  *gorm.DB
	write *gorm.DB
}

// CreateProduct implements ProductInterface.
func (i *ProductDB) CreateProduct(Product *models.Product) (*models.Product, error) {
	ret := i.write.Create(Product)
	if ret.Error != nil {
		return nil, ret.Error
	}
	return Product, nil
}

// Delete implements ProductInterface.
func (i *ProductDB) Delete(id uint64) error {
	return i.write.Delete(&models.Product{}, "id = ?", id).Error
}

// Get implements ProductInterface.
func (i *ProductDB) Get(id uint64, tx *gorm.DB) (*models.Product, error) {
	Product := &models.Product{}
	db := tx
	if db == nil {
		db = i.read
	}
	ret := db.Where("id = ?", id).Where("is_deleted = false").First(Product)
	if ret.Error != nil {
		return nil, ret.Error
	}
	return Product, nil
}

// Update implements ProductInterface.
func (i *ProductDB) Update(Product *models.Product, tx *gorm.DB) (*models.Product, error) {
	db := tx

	if db == nil {
		db = i.write
	}
	result := db.Model(&models.Product{}).Where("id = ?", Product.Id).Updates(Product)

	if result.Error != nil {
		return nil, result.Error // Return the actual error
	}

	if result.RowsAffected == 0 {
		return nil, errors.New("no product found with the given ID")
	}

	return Product, nil
}

// UpdateInventory implements ProductInterface.
// This function is needed to update quantity, and handle edge case where inventory is 0.
// Using normal `Updates` with inventory = 0 will NOT update the inventory to 0.
func (i *ProductDB) UpdateInventory(id, inventory uint64, tx *gorm.DB) error {
	db := tx

	if db == nil {
		db = i.write
	}
	result := db.Model(&models.Product{}).Where("id = ?", id).Update("inventory", inventory)

	if result.Error != nil {
		return result.Error // Return the actual error
	}

	if result.RowsAffected == 0 {
		return errors.New("no product found with the given ID")
	}

	return nil
}

func (i *ProductDB) List(limit uint64, cursorID uint64) ([]*models.Product, uint64, uint64, error) {
	var products []*models.Product

	query := i.read.Order("id ASC").Where("is_deleted = false").Limit(int(limit))
	// Count the total number of products
	var totalProducts int64
	if err := i.read.Model(&models.Product{}).Where("is_deleted = false").Distinct("id").Where("is_deleted = false").Count(&totalProducts).Error; err != nil {
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

func (i *ProductDB) ListByMerchantId(merchantId uint64, limit uint64, cursorID uint64) ([]*models.Product, uint64, uint64, error) {
	var products []*models.Product

	query := i.read.Order("id ASC").Where("is_deleted = false").Where("merchant_id = ?", merchantId).Limit(int(limit))
	// Count the total number of products
	var totalProducts int64
	if err := i.read.Model(&models.Product{}).Where("is_deleted = false").Where("merchant_id = ?", merchantId).Distinct("id").Count(&totalProducts).Error; err != nil {
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

// Update implements ProductInterface.
func (i *ProductDB) UpdateImageUrl(product *models.Product) (*models.Product, error) {

	// Update the images field in the database
	result := i.write.Model(&models.Product{}).Where("id = ?", product.Id).Update(
		"images", product.Images,
	)

	if result.Error != nil {
		return nil, result.Error // Return the actual error
	}

	if result.RowsAffected == 0 {
		return nil, errors.New("no product found with the given ID")
	}

	// Update the Product object with the new images
	return product, nil
}

func NewProductTable(read, write *gorm.DB) ProductInterface {
	StorageInstance.AutoMigrate(&models.Product{})
	return &ProductDB{
		read:  read,
		write: write,
	}
}

func GrpcToDB(product *pb.Product) *models.Product {
	return &models.Product{
		Id:              product.Id,
		Name:            product.Name,
		Price:           product.Price,
		Inventory:       product.Inventory,
		Description:     product.Description,
		Images:          product.Images,
		StripePriceId:   product.StripePriceId,
		StripeProductId: product.StripeProductId,
		MerchantId:      product.MerchantId,
	}
}

func DBToGrpc(product *models.Product) *pb.Product {
	return &pb.Product{
		Id:              product.Id,
		Name:            product.Name,
		Price:           product.Price,
		Inventory:       product.Inventory,
		Description:     product.Description,
		Images:          product.Images,
		StripePriceId:   product.StripePriceId,
		StripeProductId: product.StripeProductId,
		MerchantId:      product.MerchantId,
	}
}

func DBsToGrpcs(products []*models.Product) []*pb.Product {
	var productsGrpc []*pb.Product
	for _, product := range products {
		productsGrpc = append(productsGrpc, DBToGrpc(product))
	}
	return productsGrpc
}

func GrpcsToDBs(products []*pb.Product) []*models.Product {
	var productsDB []*models.Product
	for _, product := range products {
		productsDB = append(productsDB, GrpcToDB(product))
	}
	return productsDB
}
