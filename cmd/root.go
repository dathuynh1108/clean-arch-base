package cmd

import (
	"context"

	"github.com/spf13/cobra"
)

var rootCMD = &cobra.Command{
	Short: "Root CMD",
}

func Execute() error {
	return rootCMD.ExecuteContext(context.Background())
}
