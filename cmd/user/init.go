package main

import (
	"fmt"
	"net"

	"github.com/spf13/viper"
	"google.golang.org/grpc"

	"github.com/alanchchen/go-project-skeleton/pkg/api/user"
)

func NewConnection(cfg *viper.Viper) (*grpc.ClientConn, error) {
	return grpc.Dial(APIEndpoint(cfg), grpc.WithInsecure())
}

func NewClient(conn *grpc.ClientConn) user.ServiceClient {
	return user.NewServiceClient(conn)
}

func NewTCPSocket(cfg *viper.Viper) (net.Listener, error) {
	return net.Listen("tcp", APIEndpoint(cfg))
}

func APIEndpoint(cfg *viper.Viper) string {
	return fmt.Sprintf("%s:%d", cfg.GetString("api.host"), cfg.GetInt("api.port"))
}

func NewRPCServer() *grpc.Server {
	svc := user.NewService()
	server := grpc.NewServer()
	svc.Bind(server)
	return server
}
