package commands

import (
	"fmt"

	"github.com/spf13/cobra"

	"github.com/longuan/magicbox/internal/mongobox"
)

var CreateCmd = cobra.Command{
	Use:   "create name [-s] [--shard n] [--member n] [--mongod path] [--mongos path]",
	Short: "create a mongodb cluster, replica set or shard set",
	Example: `
create test-cls --member 5
create test-cls -s --shard 2
`,
	Args: cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		isShardSet, err := cmd.Flags().GetBool("s")
		if err != nil {
			fmt.Println(err)
			return
		}
		member, err := cmd.Flags().GetUint8("member")
		if err != nil {
			fmt.Println(err)
			return
		}
		shardNum, err := cmd.Flags().GetUint8("shard")
		if err != nil {
			fmt.Println(err)
			return
		}
		mongod, err := cmd.Flags().GetString("mongod")
		if err != nil {
			fmt.Println(err)
			return
		}
		mongos, err := cmd.Flags().GetString("mongos")
		if err != nil {
			fmt.Println(err)
			return
		}
		if isShardSet {
			_, err = mongobox.NewShardSet(mongos, mongod, args[0], shardNum, 1, member)
			if err != nil {
				fmt.Println(err)
				return
			}
		} else {
			_, err = mongobox.NewReplicaSet(mongod, args[0], member, mongobox.RoleReplica)
			if err != nil {
				fmt.Println(err)
				return
			}
		}
	},
}

func init() {
	CreateCmd.Flags().BoolP("s", "s", false, "create a shard set if true")
	CreateCmd.Flags().Uint8("member", 3, "member number for each replica set")
	CreateCmd.Flags().Uint8("shard", 2, "shard number for shard set")
	CreateCmd.Flags().String("mongod", "mongod", "mongod binary file path, default is 'mongod' system command")
	CreateCmd.Flags().String("mongos", "mongos", "mongos binary file path, default is 'mongos' system command")
}
