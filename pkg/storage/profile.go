package storage

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
	// ErrProfileNotFound is returned when a profile could not be found
	ErrProfileNotFound = errors.New("profile not found")
)

// InmemProfileStorage is an in-memory profile store
type InmemProfileStorage struct {
	Default  string
	Profiles map[string]*thrapb.Profile
}

// Get returns the profile by id.  It returns nil if it is not found
func (db *InmemProfileStorage) Get(id string) (*thrapb.Profile, error) {
	if val, ok := db.Profiles[id]; ok {
		return val, nil
	}

	return nil, ErrProfileNotFound
}

// List returns a list of all loaded profiles
func (db *InmemProfileStorage) List() []*thrapb.Profile {
	out := make([]*thrapb.Profile, 0, len(db.Profiles))

	for _, v := range db.Profiles {
		out = append(out, v.Clone())
	}

	return out
}

// GetDefault returns the default profile.  It returns nil if one has not been
// previously declared
func (db *InmemProfileStorage) GetDefault() *thrapb.Profile {
	return db.Profiles[db.Default]
}

// SetDefault sets the id to the default profile.  It returns an errProfileNotFound
// error if the profile does not exist
func (db *InmemProfileStorage) SetDefault(id string) error {
	if _, ok := db.Profiles[id]; !ok {
		return ErrProfileNotFound
	}
	db.Default = id
	return nil
}

// HCLFileProfileStorage is a hcl file backed profile store
type HCLFileProfileStorage struct {
	file string
	*InmemProfileStorage
}

// NewHCLFileProfileStorage returns a new store with the default profile loaded
func NewHCLFileProfileStorage(fpath string) *HCLFileProfileStorage {
	def := thrapb.DefaultProfile()

	return &HCLFileProfileStorage{
		InmemProfileStorage: &InmemProfileStorage{
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
		return &HCLFileProfileStorage{InmemProfileStorage: db, file: fpath}, nil
	}
	return nil, err
}

// Sync writes the db to disk
func (st *HCLFileProfileStorage) Sync() error {
	b, err := hclencoder.Encode(st.InmemProfileStorage)
	if err != nil {
		return err
	}

	return ioutil.WriteFile(st.file, b, 0644)
}

// ParseProfiles parses profiles at the given path
func parseProfiles(profPath string) (*InmemProfileStorage, error) {
	var db InmemProfileStorage
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
