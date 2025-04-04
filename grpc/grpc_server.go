package grpc

import (
	"log"
	"net"

	"product/controllers"
	pb "product/proto"
	"product/services"

	"google.golang.org/grpc"
	"google.golang.org/grpc/health"
	"google.golang.org/grpc/health/grpc_health_v1"
)

func Init() {
	ClientInit()

	lis, err := net.Listen("tcp", ":50051")
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	s := grpc.NewServer()
	healthServer := health.NewServer()
	grpc_health_v1.RegisterHealthServer(s, healthServer)
	healthServer.SetServingStatus("ProductService", grpc_health_v1.HealthCheckResponse_SERVING)

	cartService := services.NewCartService(ApiServerInstance.CartServiceConn)

	pb.RegisterProductServiceServer(s, controllers.NewProductController(cartService))

	log.Printf("Server listening at %v", lis.Addr())
	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
