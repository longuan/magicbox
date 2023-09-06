package mongobox

type CreateReplicaSetOption struct {
	MemNum     uint32 // 节点数量
	MongodFile string
	Hidden     uint32 // hidden节点数量
	KeyFile    string
}

type CreateShardSetOption struct {
	CreateReplicaSetOption
	MongosFile string
	ShardNum   uint32
	MongosNum  uint32
}

var (
	defaultCreateReplicaSetOption = CreateReplicaSetOption{
		MemNum:     3,
		MongodFile: "mongod", // mongod command
		Hidden:     0,        // no hidden
		KeyFile:    "",       // will not use key file
	}

	defaultCreateShardSetOption = CreateShardSetOption{
		CreateReplicaSetOption: defaultCreateReplicaSetOption,
		MongosFile:             "mongos",
		ShardNum:               2,
		MongosNum:              1,
	}
)
