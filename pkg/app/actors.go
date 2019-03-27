package app

import (
	"net"
	"os"
	"os/signal"
	"syscall"

	"github.com/getamis/sirius/rpc"
	"go.uber.org/dig"
)

type Logger interface {
	Println(args ...interface{})
}

type Actor interface {
	Run() error
	Interrupt(error)
}

type ActorResult struct {
	dig.Out
	Actor Actor `group:"actors"`
}

type ActorsResult struct {
	dig.In
	Actors []Actor `group:"actors"`
}

// ----------------------------------------------------------------------------

func NewRPCServerActor(server *rpc.Server, socket net.Listener, logger Logger) ActorResult {
	return ActorResult{
		Actor: &rpcServerActor{
			server: server,
			socket: socket,
			logger: logger,
		},
	}
}

type rpcServerActor struct {
	server *rpc.Server
	socket net.Listener
	logger Logger
}

func (r *rpcServerActor) Run() error {
	r.logger.Println("Starting grpc server at", r.socket.Addr())
	return r.server.Serve(r.socket)
}

func (r *rpcServerActor) Interrupt(err error) {
	r.server.Shutdown()
}

// ----------------------------------------------------------------------------

func NewSignalActor(logger Logger) ActorResult {
	return ActorResult{
		Actor: &signalActor{
			sigs:   make(chan os.Signal),
			logger: logger,
		},
	}
}

type signalActor struct {
	sigs   chan os.Signal
	logger Logger
}

func (r *signalActor) Run() error {
	signal.Notify(r.sigs, syscall.SIGTERM, syscall.SIGINT)
	defer signal.Stop(r.sigs)
	if sig := <-r.sigs; sig != nil {
		r.logger.Println("Received signal", sig, "... shutting down")
	}
	return nil
}

func (r *signalActor) Interrupt(err error) {
	close(r.sigs)
}
