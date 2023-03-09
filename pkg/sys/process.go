package sys

import (
	"context"
	"os"
	"os/exec"
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
