package mongobox

import (
	"context"
	"fmt"
	"os"
	"path"
	"strconv"

	"github.com/pkg/errors"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"

	"github.com/longuan/magicbox/pkg/sys"
)

type ReplicaSet interface {
	Init() error
	RsName() string
	Members() []HostAndPort
	ConnString() string
}

func NewReplicaSet(mongod, repl string, memNum uint8, hidden bool) (ReplicaSet, error) {
	return newRs(mongod, repl, memNum, hidden, roleReplica)
}

func newShard(mongod, repl string, memNum uint8, hidden bool) (ReplicaSet, error) {
	return newRs(mongod, repl, memNum, hidden, roleShardSvr)
}

func newCfgsvr(mongod, repl string, memNum uint8, hidden bool) (ReplicaSet, error) {
	return newRs(mongod, repl, memNum, hidden, roleConfigSvr)
}

func newRs(mongod, repl string, memNum uint8, hidden bool, role mongodRole) (ReplicaSet, error) {
	if hidden {
		return newReplicaSetWithHidden(mongod, repl, memNum, role)
	} else {
		return newPureReplicaSet(mongod, repl, memNum, role)
	}
}

type pureReplicaSet struct {
	replName string
	members  []HostAndPort
	provider mongoProvider
}

var _ ReplicaSet = (*pureReplicaSet)(nil)

func newPureReplicaSet(mongod, repl string, memNum uint8, role mongodRole) (*pureReplicaSet, error) {
	r := &pureReplicaSet{
		replName: repl,
		members:  make([]HostAndPort, 0),
		provider: &localProcessProvider{},
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
		startPort++
		r.members = append(r.members, HostAndPort{"127.0.0.1", startPort})
	}

	return r, nil
}

func (r *pureReplicaSet) Init() error {
	members := bson.A{}
	for i, memb := range r.members {
		membConfig := bson.D{
			{"_id", i},
			{"host", memb},
		}

		members = append(members, membConfig)
	}

	config := bson.D{
		{"_id", r.replName},
		{"members", members},
	}

	cliOpts := options.Client().SetHosts([]string{r.members[0].Address()}).SetDirect(true)
	cli, err := mongo.NewClient(cliOpts)
	if err != nil {
		return errors.Wrap(err, "mongo.NewClient error")
	}
	err = cli.Connect(context.Background())
	if err != nil {
		return errors.Wrap(err, "cli.Connect error")
	}
	defer cli.Disconnect(context.Background())

	var doc bson.M
	err = cli.Database("admin").RunCommand(context.Background(), bson.M{"replSetInitiate": config}).Decode(&doc)
	if err != nil {
		return errors.Wrapf(err, "replSetInitiate error config is %s", config)
	}
	return nil
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
	pr    *pureReplicaSet
	seeds []HostAndPort
}

var _ ReplicaSet = (*replicaSetWithHidden)(nil)

func newReplicaSetWithHidden(mongod, repl string, memNum uint8, role mongodRole) (*replicaSetWithHidden, error) {
	if memNum < 3 {
		return nil, errors.New("members number should greate than 3 if setting hidden")
	}

	pr, err := newPureReplicaSet(mongod, repl, memNum, role)
	if err != nil {
		return nil, err
	}
	r := &replicaSetWithHidden{
		pr:    pr,
		seeds: make([]HostAndPort, 0),
	}

	for i := 0; i < int(memNum); i++ {
		if i != int(memNum) {
			r.seeds = append(r.seeds, r.pr.members[i])
		}
	}

	return r, nil
}

func (r *replicaSetWithHidden) Init() error {
	members := bson.A{}
	for i, memb := range r.pr.members {
		membConfig := bson.D{
			{"_id", i},
			{"host", memb},
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
		{"_id", r.pr.replName},
		{"members", members},
	}

	cliOpts := options.Client().SetHosts([]string{r.seeds[0].Address()}).SetDirect(true)
	cli, err := mongo.NewClient(cliOpts)
	if err != nil {
		return errors.Wrap(err, "mongo.NewClient error")
	}
	err = cli.Connect(context.Background())
	if err != nil {
		return errors.Wrap(err, "cli.Connect error")
	}
	defer cli.Disconnect(context.Background())

	var doc bson.M
	err = cli.Database("admin").RunCommand(context.Background(), bson.M{"replSetInitiate": config}).Decode(&doc)
	if err != nil {
		return errors.Wrapf(err, "replSetInitiate error config is %s", config)
	}
	return nil
}

func (r *replicaSetWithHidden) RsName() string {
	return r.pr.RsName()
}

func (r *replicaSetWithHidden) Members() []HostAndPort {
	return r.pr.Members()
}

func (r *replicaSetWithHidden) ConnString() string {
	return ConnStringForRs(r.pr.RsName(), r.seeds)
}
