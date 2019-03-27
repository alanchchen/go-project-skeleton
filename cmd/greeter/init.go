package main

import (
	"fmt"
	"net"

	"github.com/getamis/sirius/rpc"
	"github.com/oklog/run"
	"github.com/spf13/viper"

	"github.com/alanchchen/go-project-skeleton/pkg/api/greeter"
)

func NewTCPSocket() (net.Listener, error) {
	return net.Listen("tcp", APIEndpoint())
}

func APIEndpoint() string {
	return fmt.Sprintf("%s:%d", viper.GetString("api.host"), viper.GetInt("api.port"))
}

func NewRunGroup() *run.Group {
	return &run.Group{}
}

func NewRPCServer() *rpc.Server {
	return rpc.NewServer(rpc.APIs(greeter.NewService()))
}
