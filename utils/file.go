package utils

import (
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"

	homedir "github.com/mitchellh/go-homedir"
)

// FileExists returns true if a file exists at the path
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

// ParseIgnoresFile parse each line of a ignores file into a slice
func ParseIgnoresFile(filename string) ([]string, error) {
	b, err := ioutil.ReadFile(filename)
	if err == nil {
		return strings.Split(string(b), "\n"), nil
	}

	return nil, err
}

// GetAbsPath returns an absolute path for the input accounting for "", and "~"
func GetAbsPath(p string) (out string, err error) {
	if p == "" {
		out, err = os.Getwd()
	} else if strings.HasPrefix(p, "~") {
		out, err = homedir.Expand(p)
	} else if !filepath.IsAbs(p) {
		out, err = filepath.Abs(p)
	} else {
		out = p
	}
	return
}
