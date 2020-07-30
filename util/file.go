package util

import "os"

func FileExists(f string) bool {
	_, err := os.Stat(f)
	return err == nil || os.IsExist(err) || !os.IsNotExist(err)
}
