package mongobox

import (
	"context"
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/pkg/errors"
	"github.com/spf13/pflag"

	"github.com/longuan/magicbox/pkg/sys"
)

func defaultMongodFlagSet() *pflag.FlagSet {
	flagset := pflag.NewFlagSet("mongod", pflag.PanicOnError)
	flagset.Bool("bind_ip_all", false, "")
	flagset.String("dbpath", "", "")
	flagset.Int("port", 0, "")
	flagset.String("logpath", "", "")
	flagset.String("replSet", "", "")
	flagset.Bool("fork", true, "")
	flagset.Bool("configsvr", false, "")
	flagset.Bool("shardsvr", false, "")

	return flagset
}

func defaultMongosFlagSet() *pflag.FlagSet {
	flagset := pflag.NewFlagSet("mongos", pflag.PanicOnError)
	flagset.Bool("bind_ip_all", false, "")
	flagset.String("configdb", "", "")
	flagset.Int("port", 0, "")
	flagset.String("logpath", "", "")
	flagset.Bool("fork", true, "")

	return flagset
}

func convertToArgs(flagset *pflag.FlagSet) []string {
	args := make([]string, 0)
	flagset.VisitAll(func(f *pflag.Flag) {
		if f.Value.Type() == "bool" {
			if f.Value.String() == "true" {
				args = append(args, "--"+f.Name)
			}
		} else {
			args = append(args, "--"+f.Name, f.Value.String())
		}
	})

	return args
}

type localProcessProvider struct {
}

var _ mongoProvider = (*localProcessProvider)(nil)

// StartMongod 创建一个本地的mongod进程
func (l *localProcessProvider) StartMongod(binaryPath, rsName, dbPath, logFile string, port uint16,
	role mongodRole) error {
	if role == roleUnknown || role == roleStandAlone {
		return errors.Errorf("not support")
	}

	flagset := defaultMongodFlagSet()
	flagset.Set("dbpath", dbPath)
	flagset.Set("port", fmt.Sprintf("%d", port))
	flagset.Set("logpath", logFile)
	flagset.Set("replSet", rsName)
	flagset.Set("fork", "true")

	if role == roleConfigSvr {
		flagset.Set("configsvr", "true")
	} else if role == roleShardSvr {
		flagset.Set("shardsvr", "true")
	}

	args := convertToArgs(flagset)
	err := sys.NewProcess(binaryPath, args)
	if err != nil {
		return errors.Wrapf(err, "sys.NewProcess %s %v error", binaryPath, args)
	}
	return nil
}

// StartMongos 创建一个本地的mongos进程
func (l *localProcessProvider) StartMongos(binaryPath, configRs string, configs []HostAndPort, port uint16) error {
	flagset := defaultMongosFlagSet()
	flagset.Set("configdb", ConnStringForRs(configRs, configs))
	flagset.Set("port", fmt.Sprintf("%d", port))
	flagset.Set("logpath", fmt.Sprintf("/tmp/mongos-%d.log", port))
	flagset.Set("fork", "true")

	args := convertToArgs(flagset)
	err := sys.NewProcess(binaryPath, args)
	if err != nil {
		return errors.Wrapf(err, "sys.NewProcess %s %v error", binaryPath, args)
	}
	return nil
}

func (l *localProcessProvider) StopMongod(cluster string) error {
	getPidCmd := fmt.Sprintf("ps -ef | awk '/[-]-replSet %s-/{print $2}'", cluster)
	out, err := sys.RunCommand(getPidCmd)
	if err != nil {
		return errors.WithMessagef(err, "run command %s error", getPidCmd)
	}
	outStr := string(out)
	pidStrs := strings.Split(strings.TrimSpace(outStr), "\n")
	pids := make([]int, 0)
	for _, s := range pidStrs {
		p, err := strconv.Atoi(s)
		if err != nil {
			return errors.Wrapf(err, "strconv %s error", s)
		}
		pids = append(pids, p)
	}

	for _, pid := range pids {
		getCmdlineCmd := fmt.Sprintf("ps -p %d -o cmd=", pid)
		out, err = sys.RunCommand(getCmdlineCmd)
		if err != nil {
			return errors.WithMessagef(err, "run command %s error", getCmdlineCmd)
		}
		cmdStrs := strings.Split(strings.TrimSpace(string(out)), " ")
		flagset := defaultMongodFlagSet()
		err = flagset.Parse(cmdStrs[1:])
		if err != nil {
			return errors.Wrapf(err, "flagset mongod parse %s error", string(out))
		}

		waitCtx, cancel := context.WithTimeout(context.Background(), time.Minute)
		defer cancel()
		err = sys.StopProcess(waitCtx, pid)
		if err != nil {
			return errors.WithMessagef(err, "stop process %d error", pid)
		}

		dbDir, err := flagset.GetString("dbpath")
		if err != nil {
			return err
		}
		// TODO: 如果dbpath是相对路径，就要先获取进程的当前路径，然后拼成绝对路径
		err = sys.RemoveDir(dbDir)
		if err != nil {
			return err
		}

		// TODO: 删除mongodb的logpath，要考虑到mongodb的日志rotating
	}

	return nil
}

func (l *localProcessProvider) StopMongos(cluster string) error {
	getPidCmd := fmt.Sprintf("ps -ef | awk '/[-]-configdb %s\\//{print $2}'", genCfgName(cluster))
	out, err := sys.RunCommand(getPidCmd)
	if err != nil {
		return errors.WithMessagef(err, "run command %s error", getPidCmd)
	}
	outStr := string(out)
	pidStrs := strings.Split(strings.TrimSpace(outStr), "\n")
	pids := make([]int, 0)
	for _, s := range pidStrs {
		p, err := strconv.Atoi(s)
		if err != nil {
			return errors.Wrapf(err, "strconv %s error", s)
		}
		pids = append(pids, p)
	}

	for _, pid := range pids {
		getCmdlineCmd := fmt.Sprintf("ps -p %d -o cmd=", pid)
		out, err = sys.RunCommand(getCmdlineCmd)
		if err != nil {
			return errors.WithMessagef(err, "run command %s error", getCmdlineCmd)
		}
		cmdStrs := strings.Split(strings.TrimSpace(string(out)), " ")
		flagset := defaultMongosFlagSet()
		err = flagset.Parse(cmdStrs[1:])
		if err != nil {
			return errors.Wrapf(err, "flagset mongos parse %s error", string(out))
		}

		waitCtx, cancel := context.WithTimeout(context.Background(), time.Minute)
		defer cancel()
		err = sys.StopProcess(waitCtx, pid)
		if err != nil {
			return errors.WithMessagef(err, "stop process %d error", pid)
		}

		// TODO: 删除mongodb的logpath，要考虑到mongodb的日志rotating
	}

	return nil
}
