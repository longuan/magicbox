package mongobox

import (
	"fmt"
	"os"
	"path"
	"strconv"
	"strings"

	"github.com/pkg/errors"

	"github.com/longuan/magicbox/pkg/sys"
)

type localStorageModule struct {
	root string
}

func (l localStorageModule) allocateForMongod(cluster string, members []*MongodMember) error {
	clusterDir := l.getClusterDir(cluster)

	for _, mem := range members {
		base := path.Join(clusterDir, mem.role.String())

		dbPath := path.Join(base, strconv.FormatUint(uint64(mem.port), 10))
		err := os.MkdirAll(dbPath, os.ModePerm)
		if err != nil {
			return errors.Wrapf(err, "mkdir for %s error", dbPath)
		}
		mem.dataDir = dbPath
		mem.logFile = path.Join(base, fmt.Sprintf("mongod-%d.log", mem.port))
		mem.pidFile = path.Join(base, fmt.Sprintf("mongod-%d.pid", mem.port))
	}
	return nil
}

func (l localStorageModule) allocateForMongos(cluster string, members []*MongosMember) error {
	clusterDir := l.getClusterDir(cluster)

	for _, mem := range members {
		base := path.Join(clusterDir, "mongos")
		err := os.MkdirAll(base, os.ModePerm)
		if err != nil {
			return errors.Wrapf(err, "mkdir for %s error", base)
		}

		mem.logFile = path.Join(base, fmt.Sprintf("mongos-%d.log", mem.port))
		mem.pidFile = path.Join(base, fmt.Sprintf("mongos-%d.pid", mem.port))
	}
	return nil
}

func (l localStorageModule) isEmpty(cluster string) (bool, error) {
	clusterDir := l.getClusterDir(cluster)
	entries, err := os.ReadDir(clusterDir)
	if os.IsNotExist(err) {
		return true, nil
	}
	if err != nil {
		return false, err
	}
	return len(entries) == 0, nil
}

func (l localStorageModule) getClusterDir(cluster string) string {
	return path.Join(l.root, cluster)
}

func (l localStorageModule) deallocate(cluster string) error {
	clusterDir := l.getClusterDir(cluster)
	return sys.RemoveDir(clusterDir)
}

func (l localStorageModule) getAllPids(cluster string) ([]int, error) {
	clusterDir := l.getClusterDir(cluster)

	roles, err := os.ReadDir(clusterDir)
	if err != nil {
		return nil, err
	}
	pids := make([]int, 0)
	for _, role := range roles {
		if !role.IsDir() {
			continue
		}
		nodes, err := os.ReadDir(path.Join(clusterDir, role.Name()))
		if err != nil {
			return nil, err
		}
		for _, node := range nodes {
			if strings.HasPrefix(node.Name(), ".pid") {
				content, err := os.ReadFile(path.Join(clusterDir, role.Name(), node.Name()))
				if err != nil {
					return nil, err
				}
				pid, err := strconv.Atoi(strings.TrimSpace(string(content)))
				if err != nil {
					return nil, err
				}
				pids = append(pids, pid)
			}
		}
	}
	return pids, nil
}
