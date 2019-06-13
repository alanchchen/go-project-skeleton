package app

import (
	"net"

	"google.golang.org/grpc"
)

func NewGRPCServerActor(server *grpc.Server, socket net.Listener, logger Logger) ActorResult {
	return ActorResult{
		Actor: &rpcServerActor{
			server: server,
			socket: socket,
			logger: logger,
		},
	}
}

type rpcServerActor struct {
	server *grpc.Server
	socket net.Listener
	logger Logger
}

func (r *rpcServerActor) Run() error {
	r.logger.Println("Starting grpc server at", r.socket.Addr())
	return r.server.Serve(r.socket)
}

func (r *rpcServerActor) Interrupt(err error) {
	r.server.GracefulStop()
}
