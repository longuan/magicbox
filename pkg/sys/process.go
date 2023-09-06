package sys

import (
	"context"
	"fmt"
	"os"
	"os/exec"
	"strconv"
	"strings"
	"syscall"
	"time"

	"github.com/pkg/errors"
)

// NewProcess 创建一个进程
func NewProcess(command string, args []string) error {
	cmd := exec.Command(command, args...)
	return cmd.Run()
}

func RunCommand(command string) ([]byte, error) {
	cmd := exec.Command("bash", "-c", command)
	out, err := cmd.CombinedOutput()
	if err != nil {
		return nil, err
	}
	return out, nil
}

func StopProcess(ctx context.Context, pid int) error {
	proc, err := os.FindProcess(pid)
	if err != nil {
		return errors.Wrapf(err, "find process by %d error", pid)
	}
	err = proc.Signal(syscall.SIGTERM)
	if err != nil {
		return errors.Wrapf(err, "signal SIGTERM to %d error", pid)
	}

	for {
		time.Sleep(time.Second)
		select {
		case <-ctx.Done():
			return errors.Errorf("context done %s with last error %v", ctx.Err(), err)
		default:
			p, _ := os.FindProcess(pid)
			// If sig is 0, then no signal is sent,
			// but existence and permission checks are still performed;
			// this can be used to check  for  the  existence  of
			// a process ID or process group ID that the caller is permitted to signal.
			err = p.Signal(syscall.Signal(0))
			if err == os.ErrProcessDone {
				return nil
			}
		}
	}
}

func GetAllPid(pattern string) ([]int, error) {
	getPidCmd := fmt.Sprintf("ps -ef | awk '/[%c]%s/{print $2}'", pattern[0], pattern[1:])
	out, err := RunCommand(getPidCmd)
	if err != nil {
		return nil, errors.WithMessagef(err, "run command %s error", getPidCmd)
	}
	outStr := strings.TrimSpace(string(out))
	if outStr == "" {
		return nil, errors.Errorf("there is no process belongs to [%s]", pattern)
	}
	pidStrs := strings.Split(outStr, "\n")
	pids := make([]int, 0)
	for _, s := range pidStrs {
		p, err := strconv.Atoi(s)
		if err != nil {
			return nil, errors.Wrapf(err, "strconv %s error", s)
		}
		pids = append(pids, p)
	}
	return pids, nil
}
