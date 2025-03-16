package services

import (
	pb "product/proto"
	"product/storage"
)

type ProductService struct {
}

func NewProductService() *ProductService {
	return &ProductService{}
}

func (p *ProductService) GetProduct(id uint64) (*pb.Product, error) {
	product, err := storage.StorageInstance.Product.Get(id)
	if err != nil {
		return nil, err
	}
	return product, nil
}

func (p *ProductService) CreateProduct(product *pb.Product) (*pb.Product, error) {
	// product.Id = uuid.New().String()
	product, err := storage.StorageInstance.Product.CreateProduct(product)
	if err != nil {
		return nil, err
	}
	return product, nil
}

func (p *ProductService) UpdateProduct(product *pb.Product) (*pb.Product, error) {
	product, err := storage.StorageInstance.Product.Update(product)
	if err != nil {
		return product, err
	}
	return product, nil
}

func (p *ProductService) DeleteProduct(id uint64) error {
	err := storage.StorageInstance.Product.Delete(id)
	if err != nil {
		return err
	}
	return nil
}
