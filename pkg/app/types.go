package app

import (
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"go.uber.org/dig"
)

type Input = dig.In

type Config = viper.Viper

type Command = cobra.Command
