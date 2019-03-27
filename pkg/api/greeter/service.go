package greeter

import (
	"context"
	"fmt"

	grpc "google.golang.org/grpc"
)

//go:generate mockgen -source=api.pb.go -destination=mock/api.go -package=mock

func NewService() *GreeterService {
	return &GreeterService{}
}

type GreeterService struct {
}

func (s *GreeterService) Bind(server *grpc.Server) {
	RegisterServiceServer(server, s)
}

func (s *GreeterService) SayHello(ctx context.Context, req *HelloRequest) (*HelloReply, error) {
	msg := "Hello " + req.Name
	defer fmt.Println(msg)

	return &HelloReply{
		Message: msg,
	}, nil
}
