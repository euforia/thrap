package manifest

import (
	"errors"
	"io"
	"io/ioutil"
	"path/filepath"
	"strings"

	"github.com/euforia/thrap/consts"
	"github.com/euforia/thrap/thrapb"
	"github.com/euforia/thrap/utils"
	"github.com/euforia/thrap/vcs"
	"gopkg.in/yaml.v2"
)

// LoadManifest loads a hcl or yaml manifest
func LoadManifest(mfile string) (*thrapb.Stack, error) {

	if mfile == "" {
		if utils.FileExists(consts.DefaultManifestFile) {
			mfile = consts.DefaultManifestFile
		} else if utils.FileExists("thrap.hcl") {
			// This is deprecated
			mfile = "thrap.hcl"
		} else {
			return nil, errors.New("no manifest found")
		}
	}

	mpath, err := utils.GetLocalPath(mfile)
	if err != nil {
		return nil, err
	}

	var st *thrapb.Stack
	if strings.HasSuffix(mfile, ".hcl") {
		st, err = ParseHCL(mpath)
	} else {
		st, err = ParseYAML(mpath)
	}

	if err == nil {
		st.Version = vcs.GetRepoVersion(filepath.Dir(mpath)).String()
	}

	return st, err
}

// WriteYAMLManifest writes a manifest as yaml to the Writer
func WriteYAMLManifest(st *thrapb.Stack, w io.Writer) error {
	b, err := yaml.Marshal(st)
	if err == nil {
		_, err = w.Write(b)
	}
	return err
}

// ConvertHCLToYAML converts an HCL stack to a YAML stack
func ConvertHCLToYAML(infile string) error {
	stack, err := ParseHCL(infile)
	if err != nil {
		return err
	}

	i := strings.LastIndex(infile, ".")
	ymlfile := infile[:i+1] + "yml"

	b, err := yaml.Marshal(stack)
	if err != nil {
		return err
	}

	return ioutil.WriteFile(ymlfile, b, 0644)
}
