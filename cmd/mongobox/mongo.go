package main

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"

	"github.com/longuan/magicbox/cmd/mongobox/commands"
)

var rootCmd = cobra.Command{
	Use:     "mongobox",
	Short:   "mongobox is a tool for creating mongodb cluster quickly",
	Version: "v0.0.1",
	Args:    cobra.ExactArgs(1),
}

func init() {
	rootCmd.AddCommand(&commands.CreateCmd)
}

func main() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
