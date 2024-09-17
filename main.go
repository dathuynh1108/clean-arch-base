package main

import (
	"fmt"
	"os"
	"strings"

	"github.com/dathuynh1108/clean-arch-base/cmd"

	"github.com/spf13/cobra"
)

func InitCobra() {
	cobra.OnInitialize(func() {
		printAppVersion()
	})
}

func main() {
	InitCobra()
	err := cmd.Execute()
	if err != nil {
		fmt.Println("Error executing command:", err)
	}
}

func printAppVersion() {
	// Readf from file version.txt
	versionFile, err := os.ReadFile("version.txt")
	if err != nil {
		fmt.Println("=====> API service, version.txt not found <=====")
		return
	}
	fmt.Printf("=====> API service, version: %v <===== \n", strings.TrimSpace(string(versionFile)))
}
