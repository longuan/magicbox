package mongobox

import (
	"fmt"

	"github.com/pkg/errors"

	"github.com/longuan/magicbox/pkg/sys"
)

type mongodRole uint8

const (
	roleUnknown mongodRole = iota
	roleStandAlone
	roleConfigSvr
	roleShardSvr
	roleReplica
)

// newMongod 创建一个本地的mongod进程
func newMongodProcess(mongod, replName, dbPath, logFile string, port int, role mongodRole) error {
	if role == roleUnknown || role == roleStandAlone {
		return errors.Errorf("not support")
	}

	args := []string{}
	args = append(args, "--bind_ip_all")
	args = append(args, "--dbpath")
	args = append(args, dbPath)
	args = append(args, "--port")
	args = append(args, fmt.Sprintf("%d", port))
	args = append(args, "--logpath")
	args = append(args, logFile)
	args = append(args, "--replSet")
	args = append(args, replName)
	args = append(args, "--fork")

	if role == roleConfigSvr {
		args = append(args, "--configsvr")
	} else if role == roleShardSvr {
		args = append(args, "--shardsvr")
	}

	err := sys.NewProcess(mongod, args)
	if err != nil {
		return errors.Wrapf(err, "sys.NewProcess %s %v error", mongod, args)
	}
	return nil
}

//  configStr形如 <config replset name>/<host1:port>,<host2:port>
func newMongosProcess(mongos string, port int, configStr string) error {
	args := []string{}
	args = append(args, "--port")
	args = append(args, fmt.Sprintf("%d", port))
	args = append(args, "--logpath")
	args = append(args, fmt.Sprintf("/tmp/mongos-%d.log", port))
	args = append(args, "--configdb")
	args = append(args, configStr)
	args = append(args, "--fork")
	err := sys.NewProcess(mongos, args)
	if err != nil {
		return errors.Wrapf(err, "sys.NewProcess %s %v error", mongos, args)
	}
	return nil
}
