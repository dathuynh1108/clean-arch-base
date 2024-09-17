package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

func init() {
	rootCMD.AddCommand(debugCMD)
}

var debugCMD = &cobra.Command{
	Use:   "debug",
	Short: "Run debug internal",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("Debugging internal")
	},
}
