package packs

import (
	"io/ioutil"
	"path/filepath"
)

const packManfiestFile = "manifest.hcl"

const (
	webPackID = "web"
	dsPackID  = "datastore"
	devPackID = "dev"
)

type Packs struct {
	dir string
}

func New(dir string) *Packs {
	return &Packs{dir}
}

func (packs *Packs) Web() *BasePacks {
	return NewBasePacks(filepath.Join(packs.dir, webPackID))
}
func (packs *Packs) Dev() *DevPacks {
	return NewDevPacks(filepath.Join(packs.dir, devPackID))
}
func (packs *Packs) Datastore() *BasePacks {
	return NewBasePacks(filepath.Join(packs.dir, dsPackID))
}

// packs of a singular kind
type basePackSet struct {
	typ string
	dir string
}

func (packs *basePackSet) Type() string {
	return packs.typ
}

func (packs *basePackSet) List() ([]string, error) {
	files, err := ioutil.ReadDir(packs.dir)
	if err == nil {
		packs := make([]string, 0, len(files))
		for _, f := range files {
			packs = append(packs, f.Name())
		}
		return packs, nil
	}

	return nil, err
}

//
// func readManifestBytes(dir, packID string) ([]byte, error) {
// 	ldir := filepath.Join(dir, packID)
// 	return ioutil.ReadFile(filepath.Join(ldir, packManfiestFile))
// }
