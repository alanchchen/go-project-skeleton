package app

import (
	"fmt"
	"strings"

	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
	"github.com/spf13/viper"
)

func BindFlags(cmd *cobra.Command) {
	cmd.Flags().VisitAll(func(f *pflag.Flag) {
		if strings.Contains(f.Name, ".") {
			envVar := strings.ToUpper(strings.ReplaceAll(f.Name, ".", "_"))
			viper.BindEnv(f.Name, envVar)
		}

		viper.BindPFlag(f.Name, f)

		if !f.Changed && viper.IsSet(f.Name) {
			val := viper.Get(f.Name)
			cmd.Flags().Set(f.Name, fmt.Sprintf("%v", val))
		}
	})

	for _, subCmd := range cmd.Commands() {
		BindFlags(subCmd)
	}
}
