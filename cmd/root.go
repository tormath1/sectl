package cmd

import (
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var (
	// output is the format desired for the output
	output string

	RootCmd = &cobra.Command{
		Use:   "sectl",
		Short: "sectl is a tool to query SELinux",
	}
)

func init() {
	cobra.OnInitialize(initConfig)
	RootCmd.PersistentFlags().StringVar(&output, "output", "json", "output format")

	RootCmd.AddCommand(statusCmd)
}

func initConfig() {
	viper.AutomaticEnv()
}
