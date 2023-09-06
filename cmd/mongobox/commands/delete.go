package commands

import (
	"context"
	"fmt"

	"github.com/spf13/cobra"

	"github.com/longuan/magicbox/internal/mongobox"
)

func init() {
	DeleteCmd.Flags().StringVar(&destdir, "destdir", "/tmp/mongoboxdata", "The temporary dir for storing all data and log files.")
}

var DeleteCmd = cobra.Command{
	Use:   "delete <cluster name>",
	Short: "delete a MongoDB cluster",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		manager := mongobox.NewNativeClusterManager(destdir)
		err := manager.Delete(context.Background(), args[0])
		if err != nil {
			fmt.Println(err)
			return
		}
	},
}
