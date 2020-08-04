// +build !windows

package util

import (
	"os"
	"syscall"
)

func CheckUID(dir string, sys interface{}, uid int) error {
	stat := sys.(*syscall.Stat_t)
	if uint32(uid) != stat.Uid {
		return os.Chown(dir, uid, -1)
	}
	return nil
}
