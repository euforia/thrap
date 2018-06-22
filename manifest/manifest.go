package manifest

import (
	"errors"
	"strings"

	"github.com/euforia/thrap/consts"
	"github.com/euforia/thrap/thrapb"
	"github.com/euforia/thrap/utils"
)

// LoadManifest loads a thrap manifest.  A manifest begins with
// `manifest "id" {` followed by the remainder definition
func LoadManifest(mfile string) (*thrapb.Stack, error) {

	if mfile == "" {
		if utils.FileExists(consts.DefaultManifestFile) {
			mfile = "thrap.hcl"
		} else if utils.FileExists("thrap.yml") {
			mfile = "thrap.yml"
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

	return st, err
}
