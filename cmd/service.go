package cmd

import (
	"fmt"

	messagqqueue "github.com/dathuynh1108/clean-arch-base/internal/message_queue"
	"github.com/dathuynh1108/clean-arch-base/pkg/config"
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
	serviceCMD.AddCommand(telegramBotCMD)

	serviceCMD.PersistentFlags().String("config-path", "./configs/config.toml", "Config file path")
}

func initBaseServices() {
	config.InitConfig(configPath)
	redisclient.InitRedis()
	messagqqueue.InitMessageQueue()
}
