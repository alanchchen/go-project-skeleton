package app

import (
	"strings"
	"time"

	"github.com/oklog/run"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"go.uber.org/dig"
)

type Runner struct {
	container *dig.Container
	config    *viper.Viper
}

func Run(app *cobra.Command, args []string, initializers ...interface{}) error {
	runner, err := newRunner(app, args)
	if err != nil {
		return err
	}

	runnable := func(r ActorsResult) error {
		runGroup := run.Group{}

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

	return runner.execute(runnable, initializers...)
}

func RunCustom(app *cobra.Command, args []string, runnable interface{}, initializers ...interface{}) error {
	runner, err := newRunner(app, args)
	if err != nil {
		return err
	}

	return runner.execute(runnable, initializers...)
}

func newRunner(app *cobra.Command, args []string) (*Runner, error) {
	r := &Runner{
		container: dig.New(),
		config:    viper.New(),
	}

	if err := r.bindCobraCommand(app, args); err != nil {
		return nil, err
	}

	return r, nil
}

func (r *Runner) execute(runnable interface{}, initializers ...interface{}) error {
	container := r.container

	initializers = append(initializers, func() *Config {
		return r.config
	})

	for _, initFn := range initializers {
		if err := container.Provide(initFn); err != nil {
			return err
		}
	}

	// Invoke actors
	return dig.RootCause(container.Invoke(runnable))
}

func (r *Runner) bindCobraCommand(app *cobra.Command, args []string) (err error) {
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

	if err = cfg.BindPFlags(app.Flags()); err != nil {
		return err
	}

	cfg.AutomaticEnv() // read in environment variables that match
	cfg.SetEnvKeyReplacer(strings.NewReplacer(".", "_", "-", "_"))

	// If a config file is found, read it.
	if err = cfg.ReadInConfig(); err == nil {
		app.Println("Using config file:", cfg.ConfigFileUsed())
		// Ignore the error if config file is not found
	}

	return r.bindContainer()
}

func (r *Runner) bindContainer() error {
	cfg := r.config
	container := r.container

	for _, key := range cfg.AllKeys() {
		k := key
		val := cfg.Get(k)

		var getter interface{}
		switch val.(type) {
		case bool:
			getter = func() bool {
				return cfg.GetBool(k)
			}
		case int:
			getter = func() int {
				return cfg.GetInt(k)
			}
		case int32:
			getter = func() int32 {
				return cfg.GetInt32(k)
			}
		case int64:
			getter = func() int64 {
				return cfg.GetInt64(k)
			}
		case float64:
			getter = func() float64 {
				return cfg.GetFloat64(k)
			}
		case string:
			getter = func() string {
				return cfg.GetString(k)
			}
		case []string:
			getter = func() []string {
				return cfg.GetStringSlice(k)
			}
		case map[string]string:
			getter = func() map[string]string {
				return cfg.GetStringMapString(k)
			}
		case map[string][]string:
			getter = func() map[string][]string {
				return cfg.GetStringMapStringSlice(k)
			}
		case map[string]interface{}:
			getter = func() map[string]interface{} {
				return cfg.GetStringMap(k)
			}
		case time.Time:
			getter = func() time.Time {
				return cfg.GetTime(k)
			}
		case time.Duration:
			getter = func() time.Duration {
				return cfg.GetDuration(k)
			}
		default:
			continue
		}

		if err := container.Provide(getter, dig.Name(k)); err != nil {
			return err
		}
	}

	return nil
}
