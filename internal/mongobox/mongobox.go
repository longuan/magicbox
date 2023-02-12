package mongobox

var startPort uint16 = 45000

type mongodRole uint8

const (
	roleUnknown mongodRole = iota
	roleStandAlone
	roleConfigSvr
	roleShardSvr
	roleReplica
)

type mongoProvider interface {
	StartMongod(binaryPath, rsName, dbPath, logFile string, port uint16, role mongodRole) error
	StartMongos(binaryPath, configRs string, configs []HostAndPort, port uint16) error
	StopMongod() error
	StopMongos() error
}
