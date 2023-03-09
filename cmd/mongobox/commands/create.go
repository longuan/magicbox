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
		hidden, err := cmd.Flags().GetBool("hidden")
		if err != nil {
			fmt.Println(err)
			return
		}
		if hidden && member < 3 {
			fmt.Println("member should greate than 3 if setting hidden")
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
		keyfile, err := cmd.Flags().GetString("keyfile")
		if err != nil {
			fmt.Println(err)
			return
		}
		defaultMo := mongobox.MongodOption{
			Mongod: mongod,
		}

		defaultRsOption := mongobox.RsOption{
			Members: make([]mongobox.MongodOption, 0),
		}
		for i := 0; uint8(i) < member; i++ {
			defaultRsOption.Members = append(defaultRsOption.Members, defaultMo.Dup())
		}
		if hidden {
			defaultRsOption.Members[0].Hidden = true
		}

		if isShardSet {
			ssoOption := mongobox.ShardSetOption{
				ConfigSvrOption: defaultRsOption.Dup(),
				ShardOptions:    make([]mongobox.RsOption, 0),
				MongosOption:    make([]mongobox.MongosOption, 0),
			}
			for i := 0; i < int(shardNum); i++ {
				ssoOption.ShardOptions = append(ssoOption.ShardOptions, defaultRsOption.Dup())
			}
			ssoOption.MongosOption = append(ssoOption.MongosOption, mongobox.MongosOption{
				Mongos: mongos,
			})
			_, err = mongobox.NewShardSet(args[0], keyfile, ssoOption)
			if err != nil {
				fmt.Println(err)
				return
			}
		} else {
			_, err = mongobox.NewReplicaSet(args[0], keyfile, defaultRsOption.Dup())
			if err != nil {
				fmt.Println(err)
				return
			}
		}
	},
}

func init() {
	CreateCmd.Flags().Bool("s", false, "create a shard set if true")
	CreateCmd.Flags().Uint8("member", 3, "member number for each replica set")
	CreateCmd.Flags().Uint8("shard", 2, "shard number for shard set")
	CreateCmd.Flags().String("mongod", "mongod", "mongod binary file path, default is 'mongod' system command")
	CreateCmd.Flags().String("mongos", "mongos", "mongos binary file path, default is 'mongos' system command")
	CreateCmd.Flags().Bool("hidden", false, "create a replica set with a hidden member. need member>=3 ")
	CreateCmd.Flags().String("keyfile", "", "create a cluster with keyfile. file path with file name")
}
