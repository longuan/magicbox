package main

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"

	"github.com/longuan/magicbox/cmd/mongo-cluster/commands"
)

var rootCmd = cobra.Command{
	Use:     "mongobox",
	Short:   "mongobox is a tool to quickly create/delete MongoDB cluster",
	Version: "v0.0.1",
	Args:    cobra.ExactArgs(1),
}

func init() {
	rootCmd.AddCommand(&commands.CreateCmd)
	rootCmd.AddCommand(&commands.ListCmd)
	rootCmd.AddCommand(&commands.DeleteCmd)
}

func main() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
