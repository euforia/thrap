package manifest

import (
	"errors"
	"io"
	"path/filepath"
	"strings"

	"github.com/euforia/hclencoder"
	"github.com/euforia/thrap/consts"
	"github.com/euforia/thrap/thrapb"
	"github.com/euforia/thrap/utils"
	"github.com/euforia/thrap/vcs"
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

	if err == nil {
		st.Version = vcs.GetRepoVersion(filepath.Dir(mpath)).String()
	}

	return st, err
}

func WriteManifest(st *thrapb.Stack, w io.Writer) error {
	key := `manifest "` + st.ID + `"`
	out := map[string]interface{}{
		key: st,
	}

	b, err := hclencoder.Encode(&out)
	if err == nil {
		b = append(append([]byte("\n"), b...), []byte("\n")...)
		_, err = w.Write(b)
	}

	return err
}
