package main

import (
	"fmt"
	"net"

	"go.uber.org/dig"
	"google.golang.org/grpc"

	"github.com/alanchchen/go-project-skeleton/pkg/api/greeter"
)

type EndpointConfig struct {
	dig.In

	Host string `name:"api.host"`
	Port int    `name:"api.port"`
}

func (cfg EndpointConfig) Endpoint() string {
	return net.JoinHostPort(cfg.Host, fmt.Sprintf("%d", cfg.Port))
}

func NewTCPSocket(cfg EndpointConfig) (net.Listener, error) {
	return net.Listen("tcp", cfg.Endpoint())
}

func NewRPCServer() *grpc.Server {
	svc := greeter.NewService()
	server := grpc.NewServer()
	svc.Bind(server)
	return server
}
