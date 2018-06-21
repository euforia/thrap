package utils

import "os"

func FileExists(fpath string) bool {
	_, err := os.Stat(fpath)
	return err == nil
}
