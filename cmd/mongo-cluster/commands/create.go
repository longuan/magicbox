package commands

import (
	"context"
	"fmt"

	"github.com/spf13/cobra"

	"github.com/longuan/magicbox/internal/cluster"
)

var (
	isShardSet bool
	memberNum  uint8
	shardNum   uint8
	mongodFile string
	mongosFile string
	hiddenNum  uint8
	keyFile    string
	foreground bool
	destdir    string
	mongosNum  uint8
)

func init() {
	CreateCmd.Flags().BoolVar(&isShardSet, "s", false, "Create a ShardSet if true, otherwise create a ReplSet")
	CreateCmd.Flags().Uint8Var(&memberNum, "member", 3, "Number of members per ReplSet")
	CreateCmd.Flags().Uint8Var(&shardNum, "shard", 2, "Number of shards, only valid for ShardSet")
	CreateCmd.Flags().Uint8Var(&mongosNum, "mongosNum", 2, "Number of mongos, only valid for ShardSet")
	CreateCmd.Flags().StringVar(&mongodFile, "mongod", "mongod", "The binary file path for mongod. If not specified, searches for mongod from $PATH")
	CreateCmd.Flags().StringVar(&mongosFile, "mongos", "mongos", "The binary file path for mongos. If not specified, searches for mongod from $PATH")
	CreateCmd.Flags().Uint8Var(&hiddenNum, "hidden", 0, "Each ReplSet is created with n hidden member if true. need member>=3 ")
	CreateCmd.Flags().StringVar(&keyFile, "keyfile", "", "The file path for key file.")
	CreateCmd.Flags().BoolVar(&foreground, "foreground", false, "Create a cluster in the foreground. Delete the cluster after receiving CTRL+C")
	CreateCmd.Flags().StringVar(&destdir, "destdir", "/tmp/mongoboxdata", "The temporary dir for storing all data and log files.")
}

var CreateCmd = cobra.Command{
	Use:   "create name [--s] [--shard n] [--member n] [--mongod path] [--mongos path]",
	Short: "create a MongoDB cluster, a replica set or a shard set",
	Example: `	create test-cls --member 5
	create test-cls --s --shard 2`,
	Args: cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		if isShardSet {
			opt := cluster.CreateShardSetOption{
				CreateReplicaSetOption: cluster.CreateReplicaSetOption{
					MemNum:     uint32(memberNum),
					MongodFile: mongodFile,
					Hidden:     uint32(hiddenNum),
					KeyFile:    keyFile,
				},
				MongosFile: mongosFile,
				ShardNum:   uint32(shardNum),
				MongosNum:  uint32(mongosNum),
			}

			manager := cluster.NewNativeClusterManager(destdir)
			_, err := manager.CreateSS(context.Background(), args[0], &opt)
			if err != nil {
				fmt.Println(err)
				return
			}
		} else {
			opt := cluster.CreateReplicaSetOption{
				MemNum:     uint32(memberNum),
				MongodFile: mongodFile,
				Hidden:     uint32(hiddenNum),
				KeyFile:    keyFile,
			}

			manager := cluster.NewNativeClusterManager(destdir)
			_, err := manager.CreateRS(context.Background(), args[0], &opt)
			if err != nil {
				fmt.Println(err)
				return
			}
		}
	},
}
