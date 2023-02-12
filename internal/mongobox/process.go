package mongobox

import (
	"fmt"

	"github.com/pkg/errors"

	"github.com/longuan/magicbox/pkg/sys"
)

type localProcessProvider struct {
}

var _ mongoProvider = (*localProcessProvider)(nil)

// StartMongod 创建一个本地的mongod进程
func (l *localProcessProvider) StartMongod(binaryPath, rsName, dbPath, logFile string, port uint16,
	role mongodRole) error {
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
	args = append(args, rsName)
	args = append(args, "--fork")

	if role == roleConfigSvr {
		args = append(args, "--configsvr")
	} else if role == roleShardSvr {
		args = append(args, "--shardsvr")
	}

	err := sys.NewProcess(binaryPath, args)
	if err != nil {
		return errors.Wrapf(err, "sys.NewProcess %s %v error", binaryPath, args)
	}
	return nil
}

// StartMongos 创建一个本地的mongos进程
func (l *localProcessProvider) StartMongos(binaryPath, configRs string, configs []HostAndPort, port uint16) error {
	args := []string{}
	args = append(args, "--port")
	args = append(args, fmt.Sprintf("%d", port))
	args = append(args, "--logpath")
	args = append(args, fmt.Sprintf("/tmp/mongos-%d.log", port))
	args = append(args, "--configdb")
	args = append(args, ConnStringForRs(configRs, configs))
	args = append(args, "--fork")
	err := sys.NewProcess(binaryPath, args)
	if err != nil {
		return errors.Wrapf(err, "sys.NewProcess %s %v error", binaryPath, args)
	}
	return nil
}

func (l *localProcessProvider) StopMongod() error {
	return nil
}

func (l *localProcessProvider) StopMongos() error {
	return nil
}
