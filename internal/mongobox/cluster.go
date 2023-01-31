package mongobox

import (
	"context"
	"fmt"
	"strings"

	"github.com/pkg/errors"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type ShardSet struct {
	mongoss   []string
	configsvr *ReplicaSet
	shards    []*ReplicaSet
}

// NewShardSet 创建一个分片集群
func NewShardSet(mongos, mongod, clsName string, shardNum, mongosNum, mongodNum uint8) (*ShardSet, error) {
	ss := &ShardSet{
		mongoss: make([]string, 0),
		shards:  make([]*ReplicaSet, 0),
	}

	// 创建configsvrs
	replName := genCfgName(clsName)
	rs, err := NewReplicaSet(mongod, replName, mongodNum, RoleConfigSvr)
	if err != nil {
		return nil, errors.WithMessagef(err, "creating cfgsvr %s error", replName)
	}
	ss.configsvr = rs
	configStr := replName + "/" + strings.Join(rs.members, ",")

	// 启动mongos
	for i := 0; i < int(mongosNum); i++ {
		err = newMongos(mongos, startPort, configStr)
		if err != nil {
			return nil, errors.WithMessage(err, "creating mongos error")
		}
		ss.mongoss = append(ss.mongoss, fmt.Sprintf("127.0.0.1:%d", startPort))
		startPort++
	}

	cliOpts := options.Client().SetHosts([]string{ss.mongoss[0]})
	cli, err := mongo.NewClient(cliOpts)
	if err != nil {
		return nil, errors.Wrap(err, "mongo.NewClient error")
	}
	err = cli.Connect(context.Background())
	if err != nil {
		return nil, errors.Wrap(err, "cli.Connect error")
	}
	defer cli.Disconnect(context.Background())

	for i := 0; i < int(shardNum); i++ {
		replName := genRsName(clsName, uint8(i))
		rs, err := NewReplicaSet(mongod, replName, mongodNum, RoleShardSvr)
		if err != nil {
			return nil, errors.WithMessagef(err, "creating replicaset %s error", replName)
		}
		shardStr := replName + "/" + strings.Join(rs.members, ",")
		var doc bson.M
		err = cli.Database("admin").RunCommand(context.Background(), bson.M{"addShard": shardStr}).Decode(&doc)
		if err != nil {
			return nil, errors.Wrapf(err, "addShard %s error", shardStr)
		}
		ss.shards = append(ss.shards, rs)
	}
	return ss, nil
}

func genRsName(cls string, i uint8) string {
	return fmt.Sprintf("%s-rs_%d", cls, i)
}

func genCfgName(cls string) string {
	return fmt.Sprintf("%s-cfgsvr", cls)
}
