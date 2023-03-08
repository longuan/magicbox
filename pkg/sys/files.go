package sys

import (
	"os"
)

func RemoveDir(dir string) error {
	return os.RemoveAll(dir)
}
