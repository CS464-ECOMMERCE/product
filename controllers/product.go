package controllers

import (
	"context"

	pb "product/proto"
	"product/services"
)

type ProductController struct {
	pb.UnimplementedProductServiceServer
}

func NewProductController() *ProductController {
	return &ProductController{}
}

func (p *ProductController) GetProduct(ctx context.Context, message *pb.GetProductRequest) (*pb.Product, error) {
	product, err := services.NewProductService().GetProduct(message.GetId())
	if err != nil {
		return nil, err
	}
	return product, nil
}

func (p *ProductController) CreateProduct(ctx context.Context, message *pb.CreateProductRequest) (*pb.Product, error) {
	product, err := services.NewProductService().CreateProduct(&pb.Product{
		Name:        message.GetName(),
		Description: message.GetDescription(),
		Price:       message.GetPrice(),
		Inventory:   message.GetInventory(),
	})
	if err != nil {
		return nil, err
	}
	return product, nil
}

func (p *ProductController) UpdateProduct(ctx context.Context, message *pb.UpdateProductRequest) (*pb.Product, error) {
	product, err := services.NewProductService().UpdateProduct(&pb.Product{
		Id:          message.GetId(),
		Name:        message.GetName(),
		Description: message.GetDescription(),
		Price:       message.GetPrice(),
		Inventory:   message.GetInventory(),
	})
	if err != nil {
		return nil, err
	}
	return product, nil
}

func (p *ProductController) DeleteProduct(ctx context.Context, message *pb.DeleteProductRequest) (*pb.Empty, error) {
	err := services.NewProductService().DeleteProduct(message.GetId())
	if err != nil {
		return nil, err
	}
	return &pb.Empty{}, nil
}
