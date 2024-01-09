package cmd

import (
	"fmt"

	messagqqueue "github.com/dathuynh1108/clean-arch-base/internal/message_queue"
	"github.com/dathuynh1108/clean-arch-base/pkg/config"
	"github.com/dathuynh1108/clean-arch-base/pkg/database"
	redisclient "github.com/dathuynh1108/clean-arch-base/pkg/redis_client"
	"github.com/spf13/cobra"
)

var configPath string

var serviceCMD = &cobra.Command{
	Use:   "service",
	Short: "Service",
	PersistentPreRun: func(cmd *cobra.Command, args []string) {
		fmt.Println("Pre run service, init for service here...")
		configPath, _ = cmd.Flags().GetString("config-path")
		initBaseServices()
	},
}

func init() {
	rootCMD.AddCommand(serviceCMD)

	serviceCMD.AddCommand(httpCMD)

	serviceCMD.PersistentFlags().String("config-path", "./configs/config.toml", "Config file path")
}

func initBaseServices() {
	err := config.InitConfig(configPath)
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

	messagqqueue.InitMessageQueue()
}
