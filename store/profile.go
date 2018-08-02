package store

import (
	"io/ioutil"
	"path/filepath"

	"github.com/euforia/hclencoder"
	"github.com/euforia/thrap/consts"
	"github.com/euforia/thrap/thrapb"
	"github.com/hashicorp/hcl"
)

type HCLFileProfileStorage struct {
	file string
	m    map[string]*thrapb.Profile
}

func NewHCLFileProfileStorage(fpath string) *HCLFileProfileStorage {
	return &HCLFileProfileStorage{
		m: map[string]*thrapb.Profile{
			"default": thrapb.DefaultProfile(),
		},
		file: fpath,
	}
}

// LoadHCLFileProfileStorage loads all profiles in a directory containing .thrap dir
func LoadHCLFileProfileStorage(dir string) (*HCLFileProfileStorage, error) {
	fpath := filepath.Join(dir, consts.WorkDir, consts.ProfilesFile)
	m, err := parseProfiles(fpath)
	if err == nil {
		return &HCLFileProfileStorage{m: m, file: fpath}, nil
	}
	return nil, err
}

func (st *HCLFileProfileStorage) Get(id string) *thrapb.Profile {
	return st.m[id]
}

func (st *HCLFileProfileStorage) Default() *thrapb.Profile {
	return st.m["default"]
}

func (st *HCLFileProfileStorage) Sync() error {
	profs := map[string]map[string]*thrapb.Profile{
		"profiles": st.m,
	}

	b, err := hclencoder.Encode(profs)
	if err != nil {
		return err
	}

	return ioutil.WriteFile(st.file, b, 0644)
}

// ParseProfiles parses profiles at the given path
func parseProfiles(profPath string) (map[string]*thrapb.Profile, error) {
	wprofs := make(map[string]map[string]*thrapb.Profile)
	b, err := ioutil.ReadFile(profPath)
	if err != nil {
		return nil, err
	}

	err = hcl.Unmarshal(b, &wprofs)
	if err != nil {
		return nil, err
	}

	profs := wprofs["profiles"]
	for k, v := range profs {
		v.ID = k
	}
	return profs, err
}
