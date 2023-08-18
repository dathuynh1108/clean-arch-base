package cmd

import (
	"fmt"
	"os"
	"os/signal"

	httpdelivery "github.com/dathuynh1108/clean-arch-base/internal/v1/delivery/http_delivery"
	"github.com/spf13/cobra"
)

var httpCMD = &cobra.Command{
	Use:   "http",
	Short: "Serve HTTP API service",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("HTTP started")
		httpdelivery.StartHTTPServer("0.0.0.0", "3000")
		quitChan := make(chan os.Signal)
		signal.Notify(quitChan, os.Interrupt)
		<-quitChan
	},
}
