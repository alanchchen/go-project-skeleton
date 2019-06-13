package app

import (
	"os"
	"os/signal"
	"syscall"
)

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
