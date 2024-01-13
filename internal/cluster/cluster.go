package cluster

import (
	"context"
	"fmt"
)

type ClusterManager interface {
	List(context.Context) ([]string, error)
	// create a replica set
	CreateRS(ctx context.Context, name string, opt *CreateReplicaSetOption) (*ReplicaSet, error)
	// create a shard set
	CreateSS(ctx context.Context, name string, opt *CreateShardSetOption) (*ShardSet, error)
	Delete(ctx context.Context, name string) error
}

type MongodMember struct {
	host    string
	port    uint16
	role    mongodRole
	binary  string
	dataDir string
	logFile string
	keyFile string
	pidFile string
}

func (m MongodMember) Address() string {
	return fmt.Sprintf("%s:%d", m.host, m.port)
}

type MongosMember struct {
	host    string
	port    uint16
	binary  string
	logFile string
	keyFile string
	pidFile string
	configs *ReplicaSet
}

func (m MongosMember) Address() string {
	return fmt.Sprintf("%s:%d", m.host, m.port)
}
