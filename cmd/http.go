package cmd

import (
	"fmt"
	"os"
	"os/signal"

	"github.com/dathuynh1108/clean-arch-base/internal"
	"github.com/spf13/cobra"
)

var httpCMD = &cobra.Command{
	Use:   "http",
	Short: "Serve HTTP API service",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("HTTP started")

		internal.StartHTTPServer()

		quitChan := make(chan os.Signal, 1)
		signal.Notify(quitChan, os.Interrupt)
		<-quitChan
	},
}
