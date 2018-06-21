package dockerfile

import (
	"io/ioutil"
	"path/filepath"
	"strings"

	"github.com/euforia/thrap/utils"
)

const DockerIgnoresFile = ".dockerignore"

func ReadIgnoresFile(dir string) ([]string, error) {
	fpath := filepath.Join(dir, DockerIgnoresFile)
	if !utils.FileExists(fpath) {
		return []string{}, nil
	}

	b, err := ioutil.ReadFile(fpath)
	if err != nil {
		return nil, err
	}

	return strings.Split(string(b), "\n"), nil
}
