package packs

import (
	"io/ioutil"
	"path/filepath"

	"github.com/euforia/thrap/thrapb"
	"github.com/hashicorp/hcl"
)

type BasePack struct {
	dir string
	*thrapb.PackManifest
}

type BasePacks struct {
	*basePackSet
	packs map[string]*BasePack
}

func NewBasePacks(dir string) *BasePacks {
	return &BasePacks{
		basePackSet: &basePackSet{"web", dir},
		packs:       make(map[string]*BasePack),
	}
}

func (packs *BasePacks) Load(packID string) (*BasePack, error) {
	if val, ok := packs.packs[packID]; ok {
		return val, nil
	}
	pack, err := LoadBasePack(packID, packs.dir)
	if err == nil {
		packs.packs[packID] = pack
	}
	return pack, err
}

func LoadBasePack(packID, dir string) (*BasePack, error) {
	pdir := filepath.Join(dir, packID)
	b, err := ioutil.ReadFile(filepath.Join(pdir, packManfiestFile))
	if err != nil {
		return nil, err
	}

	wp := &BasePack{dir: pdir}
	var conf thrapb.PackManifest
	err = hcl.Unmarshal(b, &conf)
	if err == nil {
		wp.PackManifest = &conf
	}
	return wp, err
}
