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
	product_db, err := storage.StorageInstance.Product.Get(id)
	if err != nil {
		return nil, err
	}
	return storage.DBToGrpc(product_db), nil
}

func (p *ProductService) CreateProduct(product *pb.Product) (*pb.Product, error) {
	product_db, err := storage.StorageInstance.Product.CreateProduct(storage.GrpcToDB(product))
	if err != nil {
		return nil, err
	}
	return storage.DBToGrpc(product_db), nil
}

func (p *ProductService) UpdateProduct(product *pb.Product) (*pb.Product, error) {
	product_db, err := storage.StorageInstance.Product.Update(storage.GrpcToDB(product))
	if err != nil {
		return product, err
	}
	return storage.DBToGrpc(product_db), nil
}

func (p *ProductService) DeleteProduct(id uint64) error {
	err := storage.StorageInstance.Product.Delete(id)
	if err != nil {
		return err
	}
	return nil
}

func (p *ProductService) ListProducts(req *pb.ListProductsRequest) (*pb.ListProductsResponse, error) {

	if req.GetMerchantId() != 0 {
		products_db, cursor, total, err := storage.StorageInstance.Product.ListByMerchantId(req.GetMerchantId(), req.GetLimit(), req.GetCursor())
		if err != nil {
			return nil, err
		}
		return &pb.ListProductsResponse{
			Products: storage.DBsToGrpcs(products_db),
			Cursor:   cursor,
			Total:    total,
		}, nil
	}
	products_db, cursor, total, err := storage.StorageInstance.Product.List(req.GetLimit(), req.GetCursor())
	if err != nil {
		return nil, err
	}

	return &pb.ListProductsResponse{
		Products: storage.DBsToGrpcs(products_db),
		Cursor:   cursor,
		Total:    total,
	}, nil
}

func (p *ProductService) UpdateProductImages(product *pb.Product) (*pb.Product, error) {

	product_db, err := storage.StorageInstance.Product.UpdateImageUrl(storage.GrpcToDB(product))
	if err != nil {
		return product, err
	}
	return storage.DBToGrpc(product_db), nil
}
