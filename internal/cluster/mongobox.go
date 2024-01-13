package cluster

import "github.com/longuan/magicbox/pkg/sys"

var startingPort uint16 = 6666

func allocatePort() uint16 {
	startingPort++
	for !sys.PortIsAvailable(startingPort) {
		if startingPort > 65535 {
			startingPort = 6666
		}
		startingPort++
	}
	return startingPort
}

type mongodRole uint8

const (
	roleConfigSvr mongodRole = iota + 1
	roleShardSvr
	roleReplica
)

func (r mongodRole) String() string {
	if r == roleConfigSvr {
		return "configsvr"
	} else if r == roleShardSvr {
		return "shardsvr"
	} else if r == roleReplica {
		return "replica"
	} else {
		return "unknown"
	}
}
