package sys

import "os/exec"

// NewProcess 创建一个进程
func NewProcess(command string, args []string) error {
	cmd := exec.Command(command, args...)
	return cmd.Run()
}
