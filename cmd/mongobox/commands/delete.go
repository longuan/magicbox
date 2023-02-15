package commands

import (
	"fmt"

	"github.com/spf13/cobra"

	"github.com/longuan/magicbox/internal/mongobox"
)

var DeleteCmd = cobra.Command{
	Use:   "delete name",
	Short: "delete a mongodb cluster",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		err := mongobox.DestroyCluster(args[0])
		if err != nil {
			fmt.Println(err)
		} else {
			fmt.Println("delete success")
		}
	},
}
