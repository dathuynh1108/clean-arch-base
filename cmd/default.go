package cmd

import (
	"github.com/dathuynh1108/clean-arch-base/internal"
	"github.com/dathuynh1108/clean-arch-base/pkg/utils"

	"github.com/spf13/cobra"
)

func init() {
	rootCMD.AddCommand(defaultCMD)
}

var defaultCMD = &cobra.Command{
	Use:   "default",
	Short: "Serve default service",
	Run: func(cmd *cobra.Command, args []string) {
		startStopper := internal.StartHTTPServer(cmd.Context())

		waitForSignals()
		utils.PanicOnError(startStopper.Stop(cmd.Context()))
	},
}
