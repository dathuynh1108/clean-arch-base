package main

import (
	"fmt"

	"github.com/dathuynh1108/clean-arch-base/cmd"
	"github.com/spf13/cobra"
)

func InitCobra() {
	cobra.OnInitialize(func() {
		fmt.Println("Init cobra...")
	})
}

func main() {
	InitCobra()
	cmd.Execute()
}
