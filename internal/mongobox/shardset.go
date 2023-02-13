package mongobox

import (
	"context"
	"fmt"

	"github.com/pkg/errors"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type ShardSet struct {
	mongoss   []HostAndPort
	configsvr ReplicaSet
	shards    []ReplicaSet
	provider  mongoProvider
}

// NewShardSet 创建一个分片集群
func NewShardSet(mongos, mongod, clsName string, shardNum, mongosNum, mongodNum uint8, hidden bool) (*ShardSet, error) {
	ss := &ShardSet{
		mongoss:  make([]HostAndPort, 0),
		shards:   make([]ReplicaSet, 0),
		provider: &localProcessProvider{}, // 默认使用local process provider
	}

	// 创建configsvrs
	replName := genCfgName(clsName)
	configRs, err := newCfgsvr(mongod, replName, mongodNum, hidden)
	if err != nil {
		return nil, errors.WithMessagef(err, "new cfgsvr %s error", replName)
	}
	ss.configsvr = configRs
	err = ss.configsvr.Init()
	if err != nil {
		return nil, errors.WithMessagef(err, "init cfgsvr %s error", replName)
	}

	// 启动mongos
	for i := 0; i < int(mongosNum); i++ {
		err = ss.provider.StartMongos(mongos, replName, configRs.Members(), startPort)
		if err != nil {
			return nil, errors.WithMessage(err, "creating mongos error")
		}
		ss.mongoss = append(ss.mongoss, HostAndPort{"127.0.0.1", startPort})
		startPort++
	}

	cliOpts := options.Client().SetHosts([]string{ss.mongoss[0].Address()})
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
		rs, err := newShard(mongod, replName, mongodNum, hidden)
		if err != nil {
			return nil, errors.WithMessagef(err, "new shard %s error", replName)
		}
		err = rs.Init()
		if err != nil {
			return nil, errors.WithMessagef(err, "init shard %s error", replName)
		}
		var doc bson.M
		err = cli.Database("admin").RunCommand(context.Background(), bson.M{"addShard": rs.ConnString()}).Decode(&doc)
		if err != nil {
			return nil, errors.Wrapf(err, "addShard %s error", rs.ConnString())
		}
		ss.shards = append(ss.shards, rs)
	}
	return ss, nil
}

func (ss *ShardSet) PrettyPrint() {
	fmt.Println("mongos:")
	for _, mongos := range ss.mongoss {
		fmt.Println("  ", mongos)
	}
	fmt.Println("cfgsvr:")
	for _, cfgsvr := range ss.configsvr.Members() {
		fmt.Println("  ", cfgsvr)
	}
	for _, shard := range ss.shards {
		fmt.Println(shard.RsName(), ":")
		for _, mongod := range shard.Members() {
			fmt.Println("  ", mongod)
		}
	}
}

func genRsName(cls string, i uint8) string {
	return fmt.Sprintf("%s-rs_%d", cls, i)
}

func genCfgName(cls string) string {
	return fmt.Sprintf("%s-cfgsvr", cls)
}
