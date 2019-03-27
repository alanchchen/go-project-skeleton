package cmd

import (
	"strings"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// Viper uses the following precedence order. Each item takes precedence over the item below it:
//
//   explicit call to SeT
//   flag
//   env
//   config
//   key/value store
//   default
//
// Viper configuration keys are case insensitive.

// InitViper reads in config file and ENV variables if set.
func InitViper(app *cobra.Command, args []string) (err error) {
	appName := app.Use
	var cfgFile string
	if cfgFlag := app.Flags().Lookup("config"); cfgFlag != nil {
		cfgFile = cfgFlag.Value.String()
	}

	if cfgFile != "" { // enable ability to specify config file via flag
		viper.SetConfigFile(cfgFile)
	} else {
		viper.SetConfigName(appName)
		viper.SetConfigType("yaml")
		viper.AddConfigPath(".")
		viper.AddConfigPath("$HOME/." + appName)
		viper.AddConfigPath("$HOME")
		viper.AddConfigPath("/")
	}

	if err = bindFlags(app); err != nil {
		return err
	}

	viper.AutomaticEnv() // read in environment variables that match
	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_", "-", "_"))

	// If a config file is found, read it.
	if err = viper.ReadInConfig(); err == nil {
		app.Println("Using config file:", viper.ConfigFileUsed())
	}

	// Ignore the error if config file is not found
	return nil
}

func bindFlags(cmd *cobra.Command) error {
	if err := viper.BindPFlags(cmd.Flags()); err != nil {
		return err
	}

	for _, subcmd := range cmd.Commands() {
		if err := bindFlags(subcmd); err != nil {
			return err
		}
	}

	return nil
}
