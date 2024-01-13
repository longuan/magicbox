package cluster

import (
	"fmt"

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
	flagset.String("keyFile", "", "key file path")
	flagset.String("pidfilepath", "", "full path to pidfile")

	return flagset
}

func defaultMongosFlagSet() *pflag.FlagSet {
	flagset := pflag.NewFlagSet("mongos", pflag.PanicOnError)
	flagset.Bool("bind_ip_all", false, "")
	flagset.String("configdb", "", "")
	flagset.Int("port", 0, "")
	flagset.String("logpath", "", "")
	flagset.Bool("fork", true, "")
	flagset.String("keyFile", "", "key file path")
	flagset.String("pidfilepath", "", "full path to pidfile")

	return flagset
}

func convertToArgs(flagset *pflag.FlagSet) []string {
	args := make([]string, 0)
	flagset.VisitAll(func(f *pflag.Flag) {
		if f.Value.Type() == "bool" {
			if f.Value.String() == "true" {
				args = append(args, "--"+f.Name)
			}
			return
		} else if f.Value.Type() == "string" {
			if f.Value.String() == "" {
				return
			}
		}
		args = append(args, "--"+f.Name, f.Value.String())
	})

	return args
}

type localProcessModule struct{}

func (l localProcessModule) startMongod(rsName string, m MongodMember) error {
	flagset := defaultMongodFlagSet()
	flagset.Set("dbpath", m.dataDir)
	flagset.Set("port", fmt.Sprintf("%d", m.port))
	flagset.Set("logpath", m.logFile)
	flagset.Set("replSet", rsName)
	flagset.Set("fork", "true")
	flagset.Set("bind_ip_all", "true")
	flagset.Set("keyFile", m.keyFile)
	flagset.Set("pidfilepath", m.pidFile)

	if m.role == roleConfigSvr {
		flagset.Set("configsvr", "true")
	} else if m.role == roleShardSvr {
		flagset.Set("shardsvr", "true")
	} else if m.role != roleReplica {
		return fmt.Errorf("%v is an invalid role", m.role)
	}

	args := convertToArgs(flagset)
	err := sys.NewProcess(m.binary, args)
	if err != nil {
		return errors.Wrapf(err, "sys.NewProcess %s %v error", m.binary, args)
	}
	return nil
}

func (l localProcessModule) startMongos(m MongosMember) error {
	flagset := defaultMongosFlagSet()
	flagset.Set("configdb", m.configs.ConnString())
	flagset.Set("port", fmt.Sprintf("%d", m.port))
	flagset.Set("logpath", m.logFile)
	flagset.Set("fork", "true")
	flagset.Set("bind_ip_all", "true")
	flagset.Set("keyFile", m.keyFile)
	flagset.Set("pidfilepath", m.pidFile)

	args := convertToArgs(flagset)
	err := sys.NewProcess(m.binary, args)
	if err != nil {
		return errors.Wrapf(err, "sys.NewProcess %s %v error", m.binary, args)
	}
	return nil
}
