package main

import (
	"fmt"
	"net"

	"github.com/spf13/viper"
	"google.golang.org/grpc"

	"github.com/alanchchen/go-project-skeleton/pkg/api/greeter"
)

func NewTCPSocket(cfg *viper.Viper) (net.Listener, error) {
	return net.Listen("tcp", APIEndpoint(cfg))
}

func APIEndpoint(cfg *viper.Viper) string {
	return fmt.Sprintf("%s:%d", cfg.GetString("api.host"), cfg.GetInt("api.port"))
}

func NewRPCServer() *grpc.Server {
	svc := greeter.NewService()
	server := grpc.NewServer()
	svc.Bind(server)
	return server
}
