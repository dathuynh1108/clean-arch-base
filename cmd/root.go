package cmd

import (
	"context"
	"fmt"

	"github.com/dathuynh1108/clean-arch-base/pkg/config"
	"github.com/dathuynh1108/clean-arch-base/pkg/database"
	"github.com/dathuynh1108/clean-arch-base/pkg/logger"
	redisclient "github.com/dathuynh1108/clean-arch-base/pkg/redis_client"
	"github.com/spf13/cobra"
)

var rootCMD = &cobra.Command{
	Short: "Root CMD",
	PersistentPreRun: func(cmd *cobra.Command, args []string) {
		fmt.Println("Root cmd, init base for service ...")
		configPath, err := cmd.Flags().GetString("config-path")
		if err != nil {
			panic(err)
		}

		fmt.Println("Config path:", configPath)
		err = config.InitConfig(configPath)
		if err != nil {
			panic(err)
		}

		initBaseServices()
	},
}

func Execute() error {
	return rootCMD.ExecuteContext(context.Background())
}

func initBaseServices() {
	var err error

	err = logger.InitLogger()
	if err != nil {
		panic(err)
	}

	err = redisclient.InitRedis()
	if err != nil {
		panic(err)
	}

	err = database.InitDatabase()
	if err != nil {
		panic(err)
	}
}
