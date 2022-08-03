package mongobox

import (
	"fmt"

	"github.com/pkg/errors"

	"github.com/longuan/magicbox/pkg/sys"
)

//  configStr形如 <config replset name>/<host1:port>,<host2:port>
func newMongos(mongos string, port int, configStr string) error {
	args := []string{}
	args = append(args, "--port")
	args = append(args, fmt.Sprintf("%d", port))
	args = append(args, "--logpath")
	args = append(args, fmt.Sprintf("/tmp/mongos-%d.log", port))
	args = append(args, "--configdb")
	args = append(args, configStr)
	args = append(args, "--fork")
	err := sys.NewProcess(mongos, args)
	if err != nil {
		return errors.Wrapf(err, "sys.NewProcess %s %v error", mongos, args)
	}
	return nil
}
