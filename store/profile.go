package store

import (
	"errors"
	"io/ioutil"
	"path/filepath"

	"github.com/euforia/hclencoder"
	"github.com/euforia/thrap/consts"
	"github.com/euforia/thrap/thrapb"
	"github.com/hashicorp/hcl"
)

var (
	ErrProfileNotFound = errors.New("profile not found")
)

type profilesDB struct {
	Default  string
	Profiles map[string]*thrapb.Profile
}

// Get returns the profile by id.  It returns nil if it is not found
func (db *profilesDB) Get(id string) *thrapb.Profile {
	return db.Profiles[id]
}

// Default returns the default profile.  It returns nil if one has not been previously
// declared
func (db *profilesDB) GetDefault() *thrapb.Profile {
	return db.Profiles[db.Default]
}

// SetDefault sets the id to the default profile.  It returns an errProfileNotFound
// error if the profile does not exist
func (db *profilesDB) SetDefault(id string) error {
	if _, ok := db.Profiles[id]; !ok {
		return ErrProfileNotFound
	}
	db.Default = id
	return nil
}

// HCLFileProfileStorage is a hcl file backed profile store
type HCLFileProfileStorage struct {
	file string
	*profilesDB
}

// NewHCLFileProfileStorage returns a new store with the default profile loaded
func NewHCLFileProfileStorage(fpath string) *HCLFileProfileStorage {
	def := thrapb.DefaultProfile()

	return &HCLFileProfileStorage{
		profilesDB: &profilesDB{
			Default: def.ID,
			Profiles: map[string]*thrapb.Profile{
				def.ID: def,
			},
		},
		file: fpath,
	}
}

// LoadHCLFileProfileStorage loads all profiles in a directory containing .thrap dir
func LoadHCLFileProfileStorage(dir string) (*HCLFileProfileStorage, error) {
	fpath := filepath.Join(dir, consts.WorkDir, consts.ProfilesFile)
	db, err := parseProfiles(fpath)
	if err == nil {
		return &HCLFileProfileStorage{profilesDB: db, file: fpath}, nil
	}
	return nil, err
}

// Sync writes the db to disk
func (st *HCLFileProfileStorage) Sync() error {
	b, err := hclencoder.Encode(st.profilesDB)
	if err != nil {
		return err
	}

	return ioutil.WriteFile(st.file, b, 0644)
}

// ParseProfiles parses profiles at the given path
func parseProfiles(profPath string) (*profilesDB, error) {
	var db profilesDB
	b, err := ioutil.ReadFile(profPath)
	if err != nil {
		return nil, err
	}

	err = hcl.Unmarshal(b, &db)
	if err != nil {
		return nil, err
	}

	for k, v := range db.Profiles {
		v.ID = k
	}

	return &db, err
}
