package controllers

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"log"

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
		MerchantId:  message.GetMerchantId(),
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
		MerchantId:  message.GetMerchantId(),
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

func (p *ProductController) ListProducts(ctx context.Context, message *pb.ListProductsRequest) (*pb.ListProductsResponse, error) {
	products, err := services.NewProductService().ListProducts(message)
	if err != nil {
		return nil, err
	}
	return products, nil
}

// func (p *ProductController) UpdateProductImages(ctx context.Context, stream *pb.UpdateProductImagesRequest) (*pb.UpdateProductImagesResponse, error) {
// 	resp, err := services.NewProductService().UpdateProductImages(stream)
// 	if err != nil {
// 		return nil, err
// 	}
// 	return resp, nil
// }

func (p *ProductController) UpdateProductImages(stream pb.ProductService_UpdateProductImagesServer) error {
	fileBuffers := make(map[string]*bytes.Buffer) // Stores file contents
	var id uint64
	for {
		req, err := stream.Recv()
		if err == io.EOF {
			log.Println("EOF received")
			var uploadedURLs []string
			for filename, buffer := range fileBuffers {

				url, err := services.NewS3Service().UploadFile(filename, buffer)
				if err != nil {
					log.Println("error uploading image: ", err)
					return fmt.Errorf("error uploading image: %v", err)
				}

				uploadedURLs = append(uploadedURLs, url)
			}

			_, err := services.NewProductService().UpdateProductImages(&pb.Product{
				Id:     id,
				Images: uploadedURLs,
			})
			if err != nil {
				log.Println("error updating product: ", err)
				return fmt.Errorf("error updating product: %v", err)
			}
			return stream.SendAndClose(&pb.UpdateProductImagesResponse{UploadedFiles: uploadedURLs})
		}
		if err != nil {
			log.Println("error receiving stream: ", err)
			return fmt.Errorf("error receiving stream: %v", err)
		}

		id = req.Id
		filename := req.Filename
		if _, exists := fileBuffers[filename]; !exists {
			fileBuffers[filename] = &bytes.Buffer{}
		}

		_, err = fileBuffers[filename].Write(req.ImageData)
		if err != nil {
			log.Println("error writing to buffer: ", err)

			return fmt.Errorf("error writing to buffer: %v", err)
		}
	}

}

func (p *ProductController) ValidateProductInventory(ctx context.Context, message *pb.ValidateProductInventoryRequest) (*pb.ValidateProductInventoryResponse, error) {
	valid, err := services.NewProductService().ValidateProductInventory(message.GetProductId(), message.GetProductId())
	if err != nil {
		return nil, err
	}
	return &pb.ValidateProductInventoryResponse{Valid: valid}, nil
}
