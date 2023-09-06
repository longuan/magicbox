package mongobox

import (
	"context"
	"fmt"

	"github.com/pkg/errors"

	pkgmongo "github.com/longuan/magicbox/pkg/mongo"
)

type ShardSet struct {
	createOpt CreateShardSetOption
	name      string
	shards    []*ReplicaSet
	configRs  *ReplicaSet
	mongoss   []*MongosMember
}

func (ss *ShardSet) init(ctx context.Context) error {
	cli, err := pkgmongo.ConnectServer(ss.mongoss[0].Address())
	if err != nil {
		return err
	}
	defer cli.Close()

	for _, shard := range ss.shards {
		err = cli.AddShard(ctx, shard.ConnString())
		if err != nil {
			return errors.Wrapf(err, "addShard %s error", shard.ConnString())
		}
	}

	return nil
}

func genRsName(cls string, i uint8) string {
	return fmt.Sprintf("%s-rs_%d", cls, i)
}

func genCfgName(cls string) string {
	return fmt.Sprintf("%s-cfgsvr", cls)
}
