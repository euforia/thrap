package utils

import (
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
)

func FileExists(fpath string) bool {
	_, err := os.Stat(fpath)
	return err == nil
}

// GetLocalPath computes the path from the user specified args.  Uses the
// current directory if none is supplied in args
func GetLocalPath(in string) (dirpath string, err error) {
	// Assume cwd
	if len(in) == 0 {
		return os.Getwd()
	}

	// Assume cwd + supplied path if not an absolute path
	if !filepath.IsAbs(in) {
		var wd string
		if wd, err = os.Getwd(); err == nil {
			dirpath = filepath.Join(wd, in)
		}
	}

	return
}

func ParseIgnoresFile(filename string) ([]string, error) {
	b, err := ioutil.ReadFile(filename)
	if err == nil {
		return strings.Split(string(b), "\n"), nil
	}

	return nil, err
}
