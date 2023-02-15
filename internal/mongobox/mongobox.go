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
	StopMongod(cluster string) error
	StopMongos(cluster string) error
}

func DestroyCluster(name string) error {
	provider := getProvider()
	err := provider.StopMongod(name)
	if err != nil {
		return err
	}
	err = provider.StopMongos(name)
	if err != nil {
		return err
	}
	return nil
}

func getProvider() mongoProvider {
	// 目前只有localProcessProvider
	return &localProcessProvider{}
}
