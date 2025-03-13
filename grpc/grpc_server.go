package grpc

import (
	"log"
	"net"

	"product/controllers"
	pb "product/proto"

	"google.golang.org/grpc"
	"google.golang.org/grpc/health"
   	"google.golang.org/grpc/health/grpc_health_v1"
)

func Init() {
	lis, err := net.Listen("tcp", ":50050")
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	s := grpc.NewServer()
	healthServer := health.NewServer()
	grpc_health_v1.RegisterHealthServer(s, healthServer)
	healthServer.SetServingStatus("ProductService", grpc_health_v1.HealthCheckResponse_SERVING)
	pb.RegisterProductServiceServer(s, controllers.NewProductController())

	log.Printf("Server listening at %v", lis.Addr())
	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
