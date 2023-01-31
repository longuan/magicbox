package mongobox

import (
	"fmt"

	"github.com/pkg/errors"

	"github.com/longuan/magicbox/pkg/sys"
)

type MongodRole uint8

const (
	RoleUnknown = iota
	RoleStandAlone
	RoleConfigSvr
	RoleShardSvr
	RoleReplica
)

// newMongod 创建一个本地的mongod进程
func newMongod(mongod, replName, dbPath, logFile string, port int, role MongodRole) error {
	if role == RoleUnknown || role == RoleStandAlone {
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

	if role == RoleConfigSvr {
		args = append(args, "--configsvr")
	} else if role == RoleShardSvr {
		args = append(args, "--shardsvr")
	}

	err := sys.NewProcess(mongod, args)
	if err != nil {
		return errors.Wrapf(err, "sys.NewProcess %s %v error", mongod, args)
	}
	return nil
}
