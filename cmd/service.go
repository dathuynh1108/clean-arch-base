package cmd

import (
	"github.com/spf13/cobra"
)

var serviceCMD = &cobra.Command{
	Use:   "service",
	Short: "Service",
}

func init() {
	rootCMD.AddCommand(serviceCMD)
	rootCMD.PersistentFlags().StringP("config", "c", "./configs/config.toml", "Config file path")

	serviceCMD.AddCommand(httpCMD)
}
