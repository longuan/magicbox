package mongobox

import (
	"context"
	"fmt"
	"os"
	"path"
	"strconv"

	"github.com/pkg/errors"
	"go.mongodb.org/mongo-driver/bson"

	"github.com/longuan/magicbox/pkg/mongo"
	"github.com/longuan/magicbox/pkg/sys"
)

type ReplicaSet interface {
	Init() error
	RsName() string
	Members() []HostAndPort
	ConnString() string
}

func NewReplicaSet(mongod, repl string, memNum uint8, hidden bool) (ReplicaSet, error) {
	rs, err := newRs(mongod, repl, memNum, hidden, roleReplica)
	if err != nil {
		return nil, errors.WithMessage(err, "new rs")
	}
	err = rs.Init()
	return rs, err
}

func newShard(mongod, repl string, memNum uint8, hidden bool) (ReplicaSet, error) {
	return newRs(mongod, repl, memNum, hidden, roleShardSvr)
}

func newCfgsvr(mongod, repl string, memNum uint8, hidden bool) (ReplicaSet, error) {
	return newRs(mongod, repl, memNum, hidden, roleConfigSvr)
}

func newRs(mongod, repl string, memNum uint8, hidden bool, role mongodRole) (ReplicaSet, error) {
	if hidden {
		return newReplicaSetWithHidden(mongod, repl, memNum, role, &localProcessProvider{})
	} else {
		return newPureReplicaSet(mongod, repl, memNum, role, &localProcessProvider{})
	}
}

type pureReplicaSet struct {
	replName string
	members  []HostAndPort
	provider mongoProvider
}

var _ ReplicaSet = (*pureReplicaSet)(nil)

func newPureReplicaSet(mongod, repl string, memNum uint8, role mongodRole, p mongoProvider) (
	*pureReplicaSet, error) {
	r := &pureReplicaSet{
		replName: repl,
		members:  make([]HostAndPort, 0),
		provider: p,
	}

	replDbDir := fmt.Sprintf("/tmp/dbfiles-%s", r.replName)
	err := os.Mkdir(replDbDir, os.ModePerm)
	if err != nil {
		return nil, errors.Wrapf(err, "mkdir for %s error", replDbDir)
	}

	for i := 0; i < int(memNum); i++ {
		for !sys.PortIsAvailable(startPort) {
			startPort += 2
		}

		dbPath := path.Join(replDbDir, strconv.FormatUint(uint64(startPort), 10))
		err := os.Mkdir(dbPath, os.ModePerm)
		if err != nil {
			return nil, errors.Wrapf(err, "mkdir for %s error", replDbDir)
		}
		logFile := path.Join(replDbDir, fmt.Sprintf("mongod-%d.log", startPort))
		err = r.provider.StartMongod(mongod, r.replName, dbPath, logFile, startPort, role)
		if err != nil {
			return nil, errors.WithMessagef(err, "newMongod for %s error", repl)
		}
		r.members = append(r.members, HostAndPort{"127.0.0.1", startPort})
		startPort++
	}

	return r, nil
}

func (r *pureReplicaSet) generateRsConf() bson.D {
	members := bson.A{}
	for i, memb := range r.members {
		membConfig := bson.D{
			{"_id", i},
			{"host", memb.Address()},
		}

		members = append(members, membConfig)
	}

	config := bson.D{
		{"_id", r.replName},
		{"members", members},
	}
	return config
}

func (r *pureReplicaSet) Init() error {
	cli, err := mongo.ConnectServer(r.members[0].Address())
	if err != nil {
		return err
	}
	defer cli.Close()
	return cli.RsInit(context.Background(), r.generateRsConf())
}

func (r *pureReplicaSet) RsName() string {
	return r.replName
}

func (r *pureReplicaSet) Members() []HostAndPort {
	return r.members // TODO: return a copy
}

func (r *pureReplicaSet) ConnString() string {
	return ConnStringForRs(r.replName, r.members)
}

type replicaSetWithHidden struct {
	*pureReplicaSet
	seeds []HostAndPort
}

var _ ReplicaSet = (*replicaSetWithHidden)(nil)

func newReplicaSetWithHidden(mongod, repl string, memNum uint8, role mongodRole, p mongoProvider) (
	*replicaSetWithHidden, error) {
	if memNum < 3 {
		return nil, errors.New("members number should greate than 3 if setting hidden")
	}

	pr, err := newPureReplicaSet(mongod, repl, memNum, role, p)
	if err != nil {
		return nil, err
	}
	r := &replicaSetWithHidden{
		pureReplicaSet: pr,
		seeds:          make([]HostAndPort, 0),
	}

	for i := 0; i < int(memNum); i++ {
		if i != int(memNum) {
			r.seeds = append(r.seeds, r.members[i])
		}
	}

	return r, nil
}

func (r *replicaSetWithHidden) generateRsConf() bson.D {
	members := bson.A{}
	for i, memb := range r.members {
		membConfig := bson.D{
			{"_id", i},
			{"host", memb.Address()},
		}

		isSeed := false
		for _, seed := range r.seeds {
			if memb == seed {
				isSeed = true
				break
			}
		}

		if !isSeed {
			// set it as hidden
			membConfig = append(membConfig, bson.E{"priority", 0})
			membConfig = append(membConfig, bson.E{"hidden", true})
		}

		members = append(members, membConfig)
	}

	config := bson.D{
		{"_id", r.replName},
		{"members", members},
	}

	return config
}

func (r *replicaSetWithHidden) Init() error {
	cli, err := mongo.ConnectServer(r.members[0].Address())
	if err != nil {
		return err
	}
	defer cli.Close()
	return cli.RsInit(context.Background(), r.generateRsConf())
}

func (r *replicaSetWithHidden) ConnString() string {
	return ConnStringForRs(r.RsName(), r.seeds)
}
