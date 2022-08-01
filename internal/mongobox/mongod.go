package mongobox

import (
	"fmt"
	"os"
	"time"

	"github.com/pkg/errors"

	"github.com/longuan/magicbox/pkg/sys"
)

// NewMongod 创建一个本地的mongod进程
func newMongod(mongod string, replName string, port int) error {
	dbDir := fmt.Sprintf("/tmp/dbfiles%d", time.Now().UnixMilli())
	err := os.Mkdir(dbDir, os.ModePerm)
	if err != nil {
		return errors.Wrapf(err, "mkdir for %s error", dbDir)
	}
	args := []string{}
	args = append(args, "--dbpath")
	args = append(args, dbDir)
	args = append(args, "--port")
	args = append(args, fmt.Sprintf("%d", port))
	args = append(args, "--logpath")
	args = append(args, fmt.Sprintf("/tmp/mongod-%d.log", port))
	args = append(args, "--replSet")
	args = append(args, replName)
	args = append(args, "--fork")
	err = sys.NewProcess(mongod, args)
	if err != nil {
		return errors.Wrapf(err, "sys.NewProcess %s %v error", mongod, args)
	}
	return nil
}
