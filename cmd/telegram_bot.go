package cmd

import (
	"fmt"
	"os"
	"os/signal"

	"github.com/dathuynh1108/clean-arch-base/internal"
	"github.com/spf13/cobra"
)

var telegramBotCMD = &cobra.Command{
	Use:   "telegram_bot",
	Short: "Start telegram bot delivery",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("Telegram bot started")

		internal.StartTelegramBot()

		quitChan := make(chan os.Signal)
		signal.Notify(quitChan, os.Interrupt)
		<-quitChan
	},
}
