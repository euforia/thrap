package storage

import (
	"errors"
	"io/ioutil"
	"path/filepath"

	"github.com/euforia/kvdb"

	"github.com/euforia/hclencoder"
	"github.com/euforia/thrap/consts"
	"github.com/euforia/thrap/pkg/pb"
	"github.com/hashicorp/hcl"
)

var (
	// ErrProfileNotFound is returned when a profile could not be found
	ErrProfileNotFound = errors.New("profile not found")
)

// InmemProfileStorage is an in-memory profile store
type InmemProfileStorage struct {
	Default string
	// Needed to unmarshal.  Should not be used otherwise
	Profiles map[string]*pb.Profile
}

// Get returns the profile by id.  It returns nil if it is not found
func (db *InmemProfileStorage) Get(id string) (*pb.Profile, error) {
	if val, ok := db.Profiles[id]; ok {
		return val.Clone(), nil
	}

	return nil, ErrProfileNotFound
}

// List returns a list of all loaded profiles
func (db *InmemProfileStorage) List() []*pb.Profile {
	out := make([]*pb.Profile, 0, len(db.Profiles))

	for _, v := range db.Profiles {
		if v.ID == db.Default {
			// Default is the first in the list
			out = append([]*pb.Profile{v.Clone()}, out...)
		} else {
			out = append(out, v.Clone())
		}
	}

	return out
}

// Add adds a new profile to the store.  it returns an error if it already
// exists or required files are missing.
func (db *InmemProfileStorage) Add(profile *pb.Profile) error {
	err := profile.Validate()
	if err != nil {
		return err
	}

	if _, ok := db.Profiles[profile.ID]; ok {
		return kvdb.ErrExists
	}
	db.Profiles[profile.ID] = profile.Clone()
	return nil
}

// GetDefault returns the default profile.  It returns nil if one has not been
// previously declared
func (db *InmemProfileStorage) GetDefault() *pb.Profile {
	if val, ok := db.Profiles[db.Default]; ok {
		return val.Clone()
	}
	return nil
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
	def := pb.DefaultProfile()

	return &HCLFileProfileStorage{
		InmemProfileStorage: &InmemProfileStorage{
			Default: def.ID,
			Profiles: map[string]*pb.Profile{
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

// ReadHCLFileProfileStorage reads an hcl file
func ReadHCLFileProfileStorage(filename string) (*HCLFileProfileStorage, error) {
	db, err := parseProfiles(filename)
	if err == nil {
		return &HCLFileProfileStorage{InmemProfileStorage: db, file: filename}, nil
	}
	return nil, err
}

// ParseProfiles parses profiles at the given path
func parseProfiles(profPath string) (*InmemProfileStorage, error) {
	b, err := ioutil.ReadFile(profPath)
	if err != nil {
		return nil, err
	}

	return ParseProfiles(b)
}

// ParseProfiles parse profile hcl config bytes
func ParseProfiles(b []byte) (*InmemProfileStorage, error) {
	var db InmemProfileStorage
	err := hcl.Unmarshal(b, &db)
	if err != nil {
		return nil, err
	}

	for k, v := range db.Profiles {
		v.ID = k
		err = v.Validate()
		if err != nil {
			return nil, err
		}
	}

	return &db, err
}
