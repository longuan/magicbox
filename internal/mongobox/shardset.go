package mongobox

import (
	"context"
	"fmt"
	"os"
	"path"

	"github.com/pkg/errors"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

const MaxShardNum = 16

type ShardSetOption struct {
	ConfigSvrOption RsOption
	ShardOptions    []RsOption
	MongosOption    []MongosOption
}

func (sso ShardSetOption) validate() error {
	if len(sso.MongosOption) == 0 {
		return errors.New("at least one mongos")
	}
	if len(sso.ShardOptions) > MaxShardNum {
		return fmt.Errorf("excced max shard num %d", MaxShardNum)
	}
	return nil
}

type MongosOption struct {
	Mongos string
}

type ShardSet struct {
	mongoss   []HostAndPort
	configsvr *pureReplicaSet
	shards    []*pureReplicaSet
	provider  mongoProvider
}

// NewShardSet 创建一个分片集群
func NewShardSet(name, keyfile string, shardSetOption ShardSetOption) (*ShardSet, error) {
	err := shardSetOption.validate()
	if err != nil {
		return nil, err
	}
	ss := &ShardSet{
		mongoss:  make([]HostAndPort, 0),
		shards:   make([]*pureReplicaSet, 0),
		provider: getProvider(),
	}

	// 创建configsvrs
	replName := genCfgName(name)
	configRs, err := newCfgsvr(replName, keyfile, shardSetOption.ConfigSvrOption)
	if err != nil {
		return nil, errors.WithMessagef(err, "new cfgsvr %s error", replName)
	}
	ss.configsvr = configRs
	err = ss.configsvr.Init()
	if err != nil {
		return nil, errors.WithMessagef(err, "init cfgsvr %s error", replName)
	}

	// 启动mongos
	mongosDir := path.Join(defaultDbPathPrefix, "mongobox-mongos")
	err = os.MkdirAll(mongosDir, os.ModePerm)
	if err != nil {
		return nil, errors.Wrapf(err, "mkdir for %s error", mongosDir)
	}
	for _, mongosOption := range shardSetOption.MongosOption {
		logFile := path.Join(mongosDir, fmt.Sprintf("mongos-%d.log", startPort))
		err = ss.provider.StartMongos(mongosOption.Mongos, replName, logFile, keyfile, configRs.Seeds(), startPort)
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

	for i, shardOption := range shardSetOption.ShardOptions {
		replName := genRsName(name, uint8(i))
		rs, err := newShard(replName, keyfile, shardOption)
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

func checkRsName(rsName, clusterName string) bool {
	if rsName == clusterName {
		return true
	}
	if rsName == genCfgName(clusterName) {
		return true
	}
	for i := 0; i < MaxShardNum; i++ {
		if rsName == genRsName(clusterName, uint8(i)) {
			return true
		}
	}
	return false
}
