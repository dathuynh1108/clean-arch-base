package cmd

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"runtime"

	"github.com/dathuynh1108/clean-arch-base/pkg/config"
	"github.com/dathuynh1108/clean-arch-base/pkg/database"
	"github.com/dathuynh1108/clean-arch-base/pkg/logger"
	"github.com/dathuynh1108/clean-arch-base/pkg/minio"
	"github.com/dathuynh1108/clean-arch-base/pkg/redisclient"
	"github.com/dathuynh1108/clean-arch-base/pkg/tracer"
	"github.com/dathuynh1108/clean-arch-base/pkg/utils"

	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"golang.org/x/sys/unix"
)

func init() {
	rootCMD.PersistentFlags().StringP("config", "c", "./configs/config.toml", "Config file path")
}

var rootCMD = &cobra.Command{
	Short: "Root CMD",
	PersistentPreRun: func(cmd *cobra.Command, args []string) {
		fmt.Println("Root cmd, init base for service ...")

		configPath, err := cmd.Flags().GetString("config")
		utils.PanicOnError(err)

		fmt.Println("Config path:", configPath)
		err = config.InitConfig(configPath)
		utils.PanicOnError(err)

		initBaseServices()
		printSystemSpecs()
	},
}

func Execute() error {
	return rootCMD.ExecuteContext(context.Background())
}

func initBaseServices() {
	var err error

	err = logger.InitLogger()
	utils.PanicOnError(err)

	tracer.InitAPMEnv()

	err = redisclient.InitRedis()
	if err != nil {
		logger.GetLogger().WithError(err).Error("Failed to init redis")
		panic(err)
	}

	err = redisclient.InitLockService()
	if err != nil {
		logger.GetLogger().WithError(err).Error("Failed to init lock service")
		panic(err)
	}

	err = database.InitDatabase()
	if err != nil {
		logger.GetLogger().WithError(err).Error("Failed to init database")
		panic(err)
	}

	err = minio.InitStore()
	if err != nil {
		logger.GetLogger().WithError(err).Error("Failed to init minio store")
		panic(err)
	}
}

func printSystemSpecs() {
	var m runtime.MemStats
	runtime.ReadMemStats(&m)

	logger.GetLogger().
		WithFields(logrus.Fields{
			"CPUs": runtime.NumCPU(),
			"Memory": map[string]any{
				"Alloc":      m.Alloc,
				"TotalAlloc": m.TotalAlloc,
				"Sys":        m.Sys,
				"Num GC":     m.NumGC,
			},
		}).
		Info("Service specs")
}

func waitForSignals() {
	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, unix.SIGTERM, unix.SIGINT, unix.SIGTSTP)
	<-sigs
}
