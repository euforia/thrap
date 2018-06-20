package devpack

//
// Package devpack implements development packages for programming languages
//
import (
	"io/ioutil"
	"path/filepath"
	"strings"

	"github.com/euforia/pseudo"
	"github.com/euforia/pseudo/scope"
	"github.com/euforia/thrap/dockerfile"
	"github.com/euforia/thrap/thrapb"
	"github.com/euforia/thrap/vars"
	version "github.com/hashicorp/go-version"
	"github.com/hashicorp/hcl"
	"github.com/hashicorp/hil/ast"
	"github.com/pkg/errors"
)

const packManfiestFile = "manifest.hcl"

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

// DevPacks holds all dev packs that are available
type DevPacks struct {
	*basePackSet
	packs map[string]*DevPack
}

// NewDevPacks returns a new DevPack with the directory containing the pack data
func NewDevPacks(dir string) *DevPacks {
	return &DevPacks{
		basePackSet: &basePackSet{"dev", dir},
		packs:       make(map[string]*DevPack),
	}
}

// Load loads a pack by the id
func (packs *DevPacks) Load(packID string) (*DevPack, error) {
	if val, ok := packs.packs[packID]; ok {
		return val, nil
	}
	pack, err := LoadDevPack(packID, packs.dir)
	if err == nil {
		packs.packs[packID] = pack
	}
	return pack, err
}

type DevPack struct {
	// Language config object
	*thrapb.Language
	// directory containing pack files
	dir string
	// Variable eval vm
	vm *pseudo.VM
	// used to check version constraints
	vc []version.Constraints
}

// LoadDevPack loads a dev pack from the directory with the given id.  The id
// name folder must exist in the specified directory
func LoadDevPack(langID, dir string) (*DevPack, error) {
	ldir := filepath.Join(dir, langID)
	b, err := ioutil.ReadFile(filepath.Join(ldir, packManfiestFile))
	if err != nil {
		return nil, err
	}

	var conf thrapb.Language
	err = hcl.Unmarshal(b, &conf)
	if err != nil {
		return nil, err
	}

	lp := &DevPack{
		dir:      ldir,
		Language: &conf,
		vm:       pseudo.NewVM(),
	}
	lp.Name = langID

	err = lp.setContraints()

	return lp, err
}

func (lp *DevPack) ScopeVars() scope.Variables {
	return scope.Variables{
		vars.LangVersion: ast.Variable{
			Type:  ast.TypeString,
			Value: lp.DefaultVersion,
		},
		vars.LangID: ast.Variable{
			Type:  ast.TypeString,
			Value: lp.Name,
		},
	}
}

func (lp *DevPack) Files(variables scope.Variables) (map[string][]byte, error) {
	out := make(map[string][]byte, len(lp.ScaffoldFiles))
	for _, fpath := range lp.ScaffoldFiles {
		if isDockerfile(fpath) {
			continue
		}

		val, err := lp.parseEval(fpath, variables)
		if err != nil {
			return nil, errors.Wrap(err, fpath)
		}

		out[fpath] = []byte(val)
	}
	return out, nil
}

func (lp *DevPack) Dockerfile(svars scope.Variables) (*dockerfile.Dockerfile, error) {
	dfpath := lp.getDockerfile()
	if dfpath == "" {
		return nil, nil
	}

	dfc, err := lp.parseEval(dfpath, svars)
	if err != nil {
		return nil, err
	}

	raw, err := dockerfile.ParseBytes([]byte(dfc))
	if err != nil {
		return nil, err
	}

	return dockerfile.ParseRaw(raw), nil

}

// SetContraints sets the language version constraints from the provided
// input versions at time of initialization
func (lp *DevPack) setContraints() error {
	vc := make([]version.Constraints, len(lp.Versions))

	for i, ct := range lp.Versions {
		c, err := version.NewConstraint(ct)
		if err != nil {
			return err
		}
		vc[i] = c
	}

	lp.vc = vc

	return nil
}

func (lp *DevPack) parseEval(fpath string, svars scope.Variables) (string, error) {
	d, err := ioutil.ReadFile(filepath.Join(lp.dir, fpath))
	if err != nil {
		return "", err
	}

	result, err := lp.vm.ParseEval(string(d), svars)
	if err != nil {
		return "", err
	}

	return result.Value.(string), nil
}

func (lp *DevPack) getDockerfile() string {

	for i, fpath := range lp.ScaffoldFiles {
		if isDockerfile(fpath) {
			return lp.ScaffoldFiles[i]
		}
	}

	return ""
}

func isDockerfile(fpath string) bool {
	return strings.Contains(strings.ToLower(fpath), "dockerfile")
}
