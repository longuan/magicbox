package sys

import (
	"os"
	"path"
)

func RemoveDir(dir string) error {
	d, err := os.Open(dir)
	if err != nil {
		if os.ErrNotExist == err {
			return nil
		}
		return err
	}
	defer d.Close()
	names, err := d.Readdirnames(-1)
	if err != nil {
		return err
	}
	for _, name := range names {
		err = os.RemoveAll(path.Join(dir, name))
		if err != nil {
			return err
		}
	}
	return nil
}
