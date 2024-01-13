package commands

import (
	"github.com/spf13/cobra"
)

var ListCmd = cobra.Command{
	Use:   "list",
	Short: "list all MongoDB cluster created by mongobox",
	Args:  cobra.ExactArgs(0),
	Run: func(cmd *cobra.Command, args []string) {
	},
}
