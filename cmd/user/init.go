package main

import (
	"fmt"
	"net"

	"google.golang.org/grpc"

	"github.com/alanchchen/go-project-skeleton/pkg/api/user"
	"github.com/alanchchen/go-project-skeleton/pkg/app"
)

type EndpointConfig struct {
	app.Input

	Host string `name:"api.host"`
	Port int    `name:"api.port"`
}

func (cfg EndpointConfig) Endpoint() string {
	return net.JoinHostPort(cfg.Host, fmt.Sprintf("%d", cfg.Port))
}

func NewConnection(cfg EndpointConfig) (*grpc.ClientConn, error) {
	return grpc.Dial(cfg.Endpoint(), grpc.WithInsecure())
}

func NewClient(conn *grpc.ClientConn) user.ServiceClient {
	return user.NewServiceClient(conn)
}

func NewTCPSocket(cfg EndpointConfig) (net.Listener, error) {
	return net.Listen("tcp", cfg.Endpoint())
}

func NewRPCServer() *grpc.Server {
	svc := user.NewService()
	server := grpc.NewServer()
	svc.Bind(server)
	return server
}
