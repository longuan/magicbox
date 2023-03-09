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

const defaultDbPathPrefix = "/tmp"

type RsOption struct {
	Members []MongodOption
}

func (rp RsOption) Dup() RsOption {
	return rp
}

func (rp RsOption) validate() error {
	if len(rp.Members) == 0 {
		return errors.New("at least one member")
	}
	hidden := 0
	for _, mem := range rp.Members {
		if mem.Mongod == "" {
			return errors.New("mongod should be valid file path")
		}
		if mem.Hidden {
			hidden++
		}
	}

	if hidden > 1 {
		return errors.New("only support one hidden")
	}
	if hidden > 0 && len(rp.Members) < 3 {
		return errors.New("at least 3 members when having hidden")
	}
	return nil
}

type MongodOption struct {
	Mongod string
	Hidden bool
}

func (mo MongodOption) Dup() MongodOption {
	return mo
}

func NewReplicaSet(repl, keyfile string, rp RsOption) (*pureReplicaSet, error) {
	if repl == "" {
		return nil, errors.New("repl should be empty")
	}
	err := rp.validate()
	if err != nil {
		return nil, err
	}
	rs, err := newPureReplicaSet(repl, keyfile, rp, roleReplica, getProvider())
	if err != nil {
		return nil, errors.WithMessage(err, "new rs")
	}
	err = rs.Init()
	return rs, err
}

func newShard(name, keyfile string, rp RsOption) (*pureReplicaSet, error) {
	return newPureReplicaSet(name, keyfile, rp, roleShardSvr, getProvider())
}

func newCfgsvr(name, keyfile string, rp RsOption) (*pureReplicaSet, error) {
	return newPureReplicaSet(name, keyfile, rp, roleConfigSvr, getProvider())
}

type pureReplicaSet struct {
	replName string
	members  []HostAndPort
	seeds    []HostAndPort
	provider mongoProvider
	options  []MongodOption
}

func newPureReplicaSet(repl, keyfile string, rp RsOption, role mongodRole, p mongoProvider) (
	*pureReplicaSet, error) {
	r := &pureReplicaSet{
		replName: repl,
		members:  make([]HostAndPort, 0),
		provider: p,
		options:  rp.Members,
	}

	replDbDir := path.Join(defaultDbPathPrefix, "mongobox-"+r.replName)
	err := os.MkdirAll(replDbDir, os.ModePerm)
	if err != nil {
		return nil, errors.Wrapf(err, "mkdir for %s error", replDbDir)
	}

	for _, mem := range rp.Members {
		for !sys.PortIsAvailable(startPort) {
			startPort += 2
		}

		dbPath := path.Join(replDbDir, strconv.FormatUint(uint64(startPort), 10))
		err := os.Mkdir(dbPath, os.ModePerm)
		if err != nil {
			return nil, errors.Wrapf(err, "mkdir for %s error", replDbDir)
		}
		logFile := path.Join(replDbDir, fmt.Sprintf("mongod-%d.log", startPort))
		err = r.provider.StartMongod(mem.Mongod, r.replName, dbPath, logFile, keyfile, startPort, role)
		if err != nil {
			return nil, errors.WithMessagef(err, "newMongod for %s error", r.replName)
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

		if r.options[i].Hidden {
			// set it as hidden
			membConfig = append(membConfig, bson.E{"priority", 0})
			membConfig = append(membConfig, bson.E{"hidden", true})
		} else {
			r.seeds = append(r.seeds, memb)
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

func (r *pureReplicaSet) Seeds() []HostAndPort {
	return r.seeds // TODO: return a copy
}

func (r *pureReplicaSet) ConnString() string {
	return ConnStringForRs(r.replName, r.seeds)
}
