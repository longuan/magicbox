package mongobox

import (
	"context"
	"fmt"
	"net"
	"os"
	"path"
	"strconv"
	"time"

	"github.com/pkg/errors"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type ReplicaSet struct {
	mongod   string
	replName string
	members  []string
}

var startPort = 45000

func NewReplicaSet(mongod, repl string, memNum uint8, role MongodRole) (*ReplicaSet, error) {
	r := &ReplicaSet{
		mongod:   mongod,
		replName: repl,
	}

	replDbDir := fmt.Sprintf("/tmp/dbfiles-%s", r.replName)
	err := os.Mkdir(replDbDir, os.ModePerm)
	if err != nil {
		return nil, errors.Wrapf(err, "mkdir for %s error", replDbDir)
	}

	for i := 0; i < int(memNum); i++ {
		for !portIsAvailable(startPort) {
			startPort += 2
		}

		host := fmt.Sprintf("127.0.0.1:%d", startPort)
		dbPath := path.Join(replDbDir, strconv.Itoa(startPort))
		logFile := path.Join(replDbDir, fmt.Sprintf("mongod-%d.log", startPort))
		err := newMongod(mongod, r.replName, dbPath, logFile, startPort, role)
		if err != nil {
			return nil, errors.WithMessagef(err, "newMongod for %s error", host)
		}
		startPort++
		r.members = append(r.members, host)
	}

	err = r.initiate()
	if err != nil {
		return nil, err
	}

	return r, nil
}

func (r *ReplicaSet) initiate() error {
	members := bson.A{}
	for i, memb := range r.members {
		members = append(members, bson.D{
			{"_id", i},
			{"host", memb},
		})
	}

	config := bson.D{
		{"_id", r.replName},
		{"members", members},
	}

	cliOpts := options.Client().SetHosts([]string{r.members[0]}).SetDirect(true)
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

func (r *ReplicaSet) PrettyPrint() {
	fmt.Println("members: ")
	fmt.Println("    ", r.members)
}

func portIsAvailable(port int) bool {
	timeout := time.Second
	conn, err := net.DialTimeout("tcp", net.JoinHostPort("127.0.0.1", fmt.Sprint(port)), timeout)
	if err != nil {
		return true
	}
	if conn != nil {
		defer conn.Close()
	}
	return false
}
