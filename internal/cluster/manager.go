package cluster

import (
	"context"
	"fmt"

	"github.com/longuan/magicbox/pkg/sys"
	"github.com/pkg/errors"
	"golang.org/x/sync/errgroup"
)

type nativeClusterManager struct {
	storage localStorageModule
	process localProcessModule
}

func NewNativeClusterManager(desdir string) *nativeClusterManager {
	return &nativeClusterManager{
		storage: localStorageModule{desdir},
		process: localProcessModule{},
	}
}

func (l *nativeClusterManager) CreateRS(ctx context.Context, name string, opt *CreateReplicaSetOption) (
	*ReplicaSet, error) {
	empty, err := l.storage.isEmpty(name)
	if err != nil {
		return nil, err
	}
	if !empty {
		return nil, fmt.Errorf("%s is not empty", name)
	}
	return l.createRS(ctx, name, name, roleReplica, opt)
}

func (l *nativeClusterManager) createRS(ctx context.Context, clusterName, replName string,
	role mongodRole, opt *CreateReplicaSetOption) (
	*ReplicaSet, error) {
	if opt == nil {
		opt = &defaultCreateReplicaSetOption
	}
	if opt.Hidden >= opt.MemNum {
		return nil, fmt.Errorf("hidden should not greater than member num")
	}

	rs := &ReplicaSet{
		createOpt: *opt,
		name:      replName,
		members:   make([]*MongodMember, 0, opt.MemNum),
	}

	for i := uint32(0); i < opt.MemNum; i++ {
		member := &MongodMember{
			host:    "localhost",
			port:    allocatePort(),
			role:    role,
			binary:  opt.MongodFile,
			keyFile: opt.KeyFile,
		}

		rs.members = append(rs.members, member)
	}

	err := l.storage.allocateForMongod(clusterName, rs.members)
	if err != nil {
		return nil, err
	}

	for _, memb := range rs.members {
		err = l.process.startMongod(rs.name, *memb)
		if err != nil {
			return nil, err
		}
	}

	return rs, rs.init()
}

func (l *nativeClusterManager) CreateSS(ctx context.Context, name string, opt *CreateShardSetOption) (
	*ShardSet, error) {
	if opt == nil {
		opt = &defaultCreateShardSetOption
	}
	if opt.Hidden >= opt.MemNum {
		return nil, fmt.Errorf("hidden should not greater than member num")
	}
	if opt.MongosNum < 1 {
		return nil, errors.New("at least one mongos")
	}
	empty, err := l.storage.isEmpty(name)
	if err != nil {
		return nil, err
	}
	if !empty {
		return nil, fmt.Errorf("%s is not empty", name)
	}

	ss := &ShardSet{
		createOpt: *opt,
		name:      name,
		shards:    make([]*ReplicaSet, 0, opt.ShardNum),
		mongoss:   make([]*MongosMember, 0, opt.MongosNum),
	}

	replName := genCfgName(name)
	configRs, err := l.createRS(ctx, name, replName, roleConfigSvr, &opt.CreateReplicaSetOption)
	if err != nil {
		return nil, errors.WithMessagef(err, "new cfgsvr %s error", replName)
	}
	ss.configRs = configRs
	if err != nil {
		return nil, errors.WithMessagef(err, "init cfgsvr %s error", replName)
	}

	for i := uint32(0); i < opt.MongosNum; i++ {
		member := &MongosMember{
			host:    "localhost",
			port:    allocatePort(),
			binary:  opt.MongosFile,
			keyFile: opt.KeyFile,
			configs: ss.configRs,
		}

		ss.mongoss = append(ss.mongoss, member)
	}

	err = l.storage.allocateForMongos(name, ss.mongoss)
	if err != nil {
		return nil, err
	}

	for _, mem := range ss.mongoss {
		err = l.process.startMongos(*mem)
		if err != nil {
			return nil, err
		}
	}

	for i := 0; i < int(ss.createOpt.ShardNum); i++ {
		replName := genRsName(ss.name, uint8(i))
		rs, err := l.createRS(ctx, name, replName, roleShardSvr, &ss.createOpt.CreateReplicaSetOption)
		if err != nil {
			return nil, errors.WithMessagef(err, "new shard %s error", replName)
		}
		ss.shards = append(ss.shards, rs)
	}

	return ss, ss.init(ctx)
}

func (l *nativeClusterManager) Delete(ctx context.Context, cluster string) error {
	allPids, err := l.storage.getAllPids(cluster)
	if err != nil {
		return err
	}
	// stop all processes
	wg, _ := errgroup.WithContext(ctx)
	for _, pid := range allPids {
		pid := pid
		wg.Go(func() error {
			return sys.StopProcess(ctx, pid)
		})
	}

	err = wg.Wait()
	if err != nil {
		return err
	}

	return l.storage.deallocate(cluster)
}
