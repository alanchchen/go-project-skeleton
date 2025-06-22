package core

import (
	"context"
	"fmt"

	"github.com/alanchchen/go-project-skeleton/pkg/greeter/api"
)

func NewService() *GreeterService {
	return &GreeterService{}
}

type GreeterService struct {
	api.UnimplementedGreeterServiceServer
}

func (s *GreeterService) SayHello(ctx context.Context, req *api.HelloRequest) (*api.HelloReply, error) {
	msg := "Hello " + req.Name
	defer fmt.Println(msg)

	return &api.HelloReply{
		Message: msg,
	}, nil
}
