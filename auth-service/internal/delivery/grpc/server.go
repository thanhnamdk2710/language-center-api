package grpc

import (
	"fmt"
	"net"

	"github.com/thanhnamdk2710/auth-service/internal/config"
	"github.com/thanhnamdk2710/auth-service/internal/delivery/grpc/proto_gen"
	"github.com/thanhnamdk2710/auth-service/internal/shared/logger"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

type healthServer struct {
	proto_gen.UnimplementedHealthServiceServer
}

// Check implements the HealthCheck gRPC endpoint
func (s *healthServer) Check(ctx context.Context, req *proto_gen.HealthRequest) (*proto_gen.HealthReponse, error) {
	return &proto_gen.HealthReponse{
		Status:  "ok",
		Message: "auth-service is healthy",
	}, nil
}

func StartGRPCServer(cfg *config.Config) {
	addr := fmt.Sprintf(":%s", cfg.GRPCPort)
	lis, err := net.Listen("tcp", addr)
	if err != nil {
		logger.Error("Failed to listen on gRPC port", err)
		return
	}

	s := grpc.NewServer()
	proto_gen.RegisterHealthServiceServer(s, &healthServer{})

	reflection.Register(s)

	logger.Info(fmt.Sprintf("gRPC server running on %s", addr))
	if err := s.Serve(lis); err != nil {
		logger.Error("Failed to serve gRPC", err)
	}
}
