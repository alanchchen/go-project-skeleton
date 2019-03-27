package main

import (
	"fmt"
	"log"
	"os"

	_ "github.com/joho/godotenv/autoload"
	"github.com/oklog/run"
	"github.com/spf13/cobra"
	"go.uber.org/dig"

	"github.com/alanchchen/go-project-skeleton/cmd"
	"github.com/alanchchen/go-project-skeleton/pkg/app"
)

func init() {
	// Setup flags
	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", fmt.Sprintf("config file (default is ./%s.yaml)", appName))

	rootCmd.PersistentFlags().String("api.host", "", "the grpc server listening host")
	rootCmd.PersistentFlags().Int("api.port", 8088, "the grpc server listening port")
}

const (
	appName = "user"
)

var cfgFile string
var logger = log.New(os.Stderr, "", log.LstdFlags|log.Lshortfile)

// rootCmd is the root command
var rootCmd = &cobra.Command{
	Use:               appName,
	Short:             appName + " is an RPC server",
	Long:              appName + " is an RPC server",
	PersistentPreRunE: cmd.InitViper,
	RunE: func(cmd *cobra.Command, args []string) error {
		container := dig.New()

		initializers := []interface{}{
			// actors
			app.NewRPCServerActor,
			app.NewSignalActor,

			// actors' dependencies
			NewRunGroup,
			NewRPCServer,
			NewTCPSocket,
			func() app.Logger {
				return logger
			},
		}

		for _, initFn := range initializers {
			if err := container.Provide(initFn); err != nil {
				return err
			}
		}

		// Invoke actors
		return container.Invoke(func(runGroup *run.Group, r app.ActorsResult) error {
			for _, actor := range r.Actors {
				runGroup.Add(actor.Run, actor.Interrupt)
			}

			// Run blocks until all the actors return. In the normal case, that’ll be when someone hits ctrl-C,
			// triggering the signal handler. If something breaks, its error will be propegated through. In all
			// cases, the first returned error triggers the interrupt function for all actors. And in this way,
			// we can reliably and coherently ensure that every goroutine that’s Added to the group is stopped,
			// when Run returns.
			return runGroup.Run()
		})
	},
}

func main() {
	if err := rootCmd.Execute(); err != nil {
		logger.Println(err)
		os.Exit(1)
	}
}
