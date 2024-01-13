package cluster

import (
	"context"
	"strings"

	"go.mongodb.org/mongo-driver/bson"

	"github.com/longuan/magicbox/pkg/mongo"
)

type ReplicaSet struct {
	createOpt CreateReplicaSetOption
	name      string
	members   []*MongodMember
	seeds     []*MongodMember
}

func (rs *ReplicaSet) init() error {
	cli, err := mongo.ConnectServer(rs.members[0].Address())
	if err != nil {
		return err
	}
	defer cli.Close()
	return cli.RsInit(context.Background(), rs.generateRsConf())
}

func (rs *ReplicaSet) generateRsConf() bson.D {
	members := bson.A{}
	hiddenNum := rs.createOpt.Hidden
	for i, memb := range rs.members {
		membConfig := bson.D{
			{"_id", i},
			{"host", memb.Address()},
		}

		if i != 0 && hiddenNum > 0 {
			// set it as hidden
			membConfig = append(membConfig, bson.E{"priority", 0})
			membConfig = append(membConfig, bson.E{"hidden", true})
			hiddenNum--
		} else {
			rs.seeds = append(rs.seeds, memb)
		}

		members = append(members, membConfig)
	}

	config := bson.D{
		{"_id", rs.name},
		{"members", members},
	}
	return config
}

func (rs *ReplicaSet) ConnString() string {
	connStr := rs.name + "/"
	for _, memb := range rs.members {
		connStr += memb.Address() + ","
	}
	connStr = strings.TrimSuffix(connStr, ",")
	return connStr
}
