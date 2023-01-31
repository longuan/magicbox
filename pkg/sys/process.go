package sys

import (
	"os/exec"
	"strings"

	"github.com/pkg/errors"
)

// NewProcess 创建一个进程
func NewProcess(command string, args []string) error {
	cmd := exec.Command(command, args...)
	return cmd.Run()
}

func RunCommand(command string, args []string) ([]byte, error) {
	cmd := exec.Command(command, args...)
	out, err := cmd.CombinedOutput()
	if err != nil {
		return nil, errors.Wrap(err, command+" "+strings.Join(args, " "))
	}
	return out, err
}
