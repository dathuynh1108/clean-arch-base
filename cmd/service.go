package cmd

import (
	"fmt"
	"time"

	messagqqueue "github.com/dathuynh1108/clean-arch-base/internal/message_queue"
	"github.com/dathuynh1108/clean-arch-base/pkg/config"
	redisclient "github.com/dathuynh1108/clean-arch-base/pkg/redis_client"
	"github.com/spf13/cobra"
)

var serviceCMD = &cobra.Command{
	Use:   "service",
	Short: "Service",
	PersistentPreRun: func(cmd *cobra.Command, args []string) {
		fmt.Println("Pre run service, init for service here...")
		configPath, _ := cmd.Flags().GetString("config-path")
		config.InitConfig(configPath)

		redisclient.InitRedis()

		messagqqueue.InitMessageQueue()

		testQueue, err := messagqqueue.GetMessageQueue().OpenQueue("test")
		if err != nil {
			panic(err)
		}
		testQueue.StartConsuming(10, time.Second)
		testConsumer := messagqqueue.NewConsumer(func(payload string) error {
			fmt.Println("test queue", payload)
			return nil
		})
		testQueue.AddConsumer("test", testConsumer)
		testQueue.Publish("test payload")
	},
}

func init() {
	rootCMD.AddCommand(serviceCMD)
	serviceCMD.AddCommand(httpCMD)
	serviceCMD.PersistentFlags().String("config-path", "./configs/config.toml", "Config file path")
}
