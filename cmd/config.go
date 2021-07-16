package cmd

import (
	"github.com/dp1140a/geoip/config"
	"github.com/spf13/cobra"
)

func init() {
	RootCmd.AddCommand(configCmd)
}

var configCmd = &cobra.Command{
	Use:   "config",
	Short: "Generates and prints a default config",
	Run: func(cmd *cobra.Command, args []string) {
		config.Generate()
	},
}
