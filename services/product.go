package services

import (
	"errors"
	pb "product/proto"
	"product/storage"

	"github.com/stripe/stripe-go/v81"
)

type ProductService struct {
}

func NewProductService() *ProductService {
	return &ProductService{}
}

func (p *ProductService) GetProduct(id uint64) (*pb.Product, error) {
	product_db, err := storage.StorageInstance.Product.Get(id, nil)
	if err != nil {
		return nil, err
	}
	return storage.DBToGrpc(product_db), nil
}

func (p *ProductService) CreateProduct(product *pb.Product) (*pb.Product, error) {
	stripeProduct, err := NewStripeService().CreateNewProduct(product.Name, product.Price)
	if err != nil {
		return nil, err
	}

	product.StripeProductId = stripeProduct.ID
	product.StripePriceId = stripeProduct.DefaultPrice.ID

	product_db, err := storage.StorageInstance.Product.CreateProduct(storage.GrpcToDB(product))
	if err != nil {
		return nil, err
	}
	return storage.DBToGrpc(product_db), nil
}

func (p *ProductService) UpdateProduct(updatedProduct *pb.Product) (*pb.Product, error) {
	existingProduct, err := storage.StorageInstance.Product.Get(updatedProduct.Id, nil)
	if err != nil {
		return nil, err
	}

	// check if the price has changed
	var stripeProduct *stripe.Product
	if updatedProduct.Price != 0 && updatedProduct.Price != existingProduct.Price {
		stripeProduct, err = NewStripeService().UpdateProductPrice(existingProduct.StripeProductId, existingProduct.StripePriceId, updatedProduct.Price)
		if err != nil {
			return nil, err
		}
		updatedProduct.StripePriceId = stripeProduct.DefaultPrice.ID
	}

	product_db, err := storage.StorageInstance.Product.Update(storage.GrpcToDB(updatedProduct), nil)
	if err != nil {
		return updatedProduct, err
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

func (p *ProductService) ValidateProductInventory(id, requestedQuantity uint64) (bool, error) {
	product, err := storage.StorageInstance.Product.Get(id, nil)
	if err != nil {
		return false, err
	}

	if product.Inventory < requestedQuantity {
		return false, errors.New("insufficient inventory")
	}
	return true, nil
}
