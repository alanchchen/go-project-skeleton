package app

import (
	"strings"

	"github.com/oklog/run"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"go.uber.org/dig"
)

type Runner struct {
	container *dig.Container
	runGroup  *run.Group
	config    *viper.Viper
}

func NewRunner() *Runner {
	return &Runner{
		container: dig.New(),
		runGroup:  &run.Group{},
		config:    viper.New(),
	}
}

func (r *Runner) Run(initializers ...interface{}) error {
	runGroup := r.runGroup

	runnable := func(r ActorsResult) error {
		for _, actor := range r.Actors {
			runGroup.Add(actor.Run, actor.Interrupt)
		}

		// Run blocks until all the actors return. In the normal case, that’ll be when someone hits ctrl-C,
		// triggering the signal handler. If something breaks, its error will be propegated through. In all
		// cases, the first returned error triggers the interrupt function for all actors. And in this way,
		// we can reliably and coherently ensure that every goroutine that’s Added to the group is stopped,
		// when Run returns.
		return runGroup.Run()
	}

	return r.RunCustom(runnable, initializers...)
}

func (r *Runner) RunCustom(runnable interface{}, initializers ...interface{}) error {
	container := r.container

	initializers = append(initializers, func() *viper.Viper {
		return r.config
	})

	for _, initFn := range initializers {
		if err := container.Provide(initFn); err != nil {
			return err
		}
	}

	// Invoke actors
	return container.Invoke(runnable)
}

func (r *Runner) BindCobraCommand(app *cobra.Command, args ...string) (err error) {
	cfg := r.config
	appName := app.Use
	var cfgFile string
	if cfgFlag := app.Flags().Lookup("config"); cfgFlag != nil {
		cfgFile = cfgFlag.Value.String()
	}

	if cfgFile != "" { // enable ability to specify config file via flag
		cfg.SetConfigFile(cfgFile)
	} else {
		cfg.SetConfigName(appName)
		cfg.SetConfigType("yaml")
		cfg.AddConfigPath(".")
		cfg.AddConfigPath("$HOME/." + appName)
		cfg.AddConfigPath("$HOME")
		cfg.AddConfigPath("/")
	}

	if err = bindFlags(app, cfg); err != nil {
		return err
	}

	cfg.AutomaticEnv() // read in environment variables that match
	cfg.SetEnvKeyReplacer(strings.NewReplacer(".", "_", "-", "_"))

	// If a config file is found, read it.
	if err = cfg.ReadInConfig(); err == nil {
		app.Println("Using config file:", cfg.ConfigFileUsed())
	}

	// Ignore the error if config file is not found
	return nil
}

func bindFlags(cmd *cobra.Command, cfg *viper.Viper) error {
	if err := cfg.BindPFlags(cmd.Flags()); err != nil {
		return err
	}

	for _, subcmd := range cmd.Commands() {
		if err := bindFlags(subcmd, cfg); err != nil {
			return err
		}
	}

	return nil
}
