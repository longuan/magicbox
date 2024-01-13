package mongo

import (
	"context"

	"github.com/pkg/errors"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type server struct {
	cli *mongo.Client
}

func NewServer(cli *mongo.Client) *server {
	return &server{
		cli: cli,
	}
}

func ConnectServer(addr string, opts ...*options.ClientOptions) (*server, error) {
	cliOpts := options.Client().SetHosts([]string{addr}).SetDirect(true)
	cli, err := mongo.NewClient(cliOpts)
	if err != nil {
		return nil, errors.Wrap(err, "mongo.NewClient error")
	}
	err = cli.Connect(context.Background())
	if err != nil {
		return nil, errors.Wrap(err, "cli.Connect error")
	}
	err = cli.Ping(context.Background(), nil)
	if err != nil {
		return nil, errors.Wrap(err, "cli.Ping error")
	}
	return &server{
		cli: cli,
	}, nil
}

func (s *server) RsInit(ctx context.Context, config interface{}) error {
	result := s.cli.Database("admin").RunCommand(ctx, bson.M{"replSetInitiate": config})
	if result.Err() != nil {
		return errors.Wrapf(result.Err(), "replSetInitiate error config is %s", config)
	}
	return nil
}

func (s *server) AddShard(ctx context.Context, shard string) error {
	result := s.cli.Database("admin").RunCommand(ctx, bson.M{"addShard": shard})
	if result.Err() != nil {
		return errors.Wrapf(result.Err(), "addShard error shard is %s", shard)
	}
	return nil
}

func (s *server) Close() {
	if s.cli != nil {
		s.cli.Disconnect(context.Background())
	}
}

func (s *server) Client() *mongo.Client {
	return s.cli
}
