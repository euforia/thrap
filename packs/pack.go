package packs

import (
	"io/ioutil"
	"net/url"
	"os"
	"path/filepath"

	"github.com/euforia/thrap/utils"
	"github.com/pkg/errors"
	git "gopkg.in/src-d/go-git.v4"

	homedir "github.com/mitchellh/go-homedir"
)

const packManfiestFile = "manifest.hcl"

const (
	webPackID = "web"
	dsPackID  = "datastore"
	devPackID = "dev"
)

var (
	errPackDirRequired        = errors.New("pack directory required")
	errPackDirExists          = errors.New("pack directory exists")
	errPackSourceNotSupported = errors.New("pack source not supported")
)

// Packs is the interface to access various packs.
type Packs struct {
	// local pack path.  /{path}/{to}/packs
	dir string
}

// New returns a new packs manager
func New(dir string) (*Packs, error) {
	if dir == "" {
		return nil, errPackDirRequired
	}

	if !filepath.IsAbs(dir) {
		var (
			p   string
			err error
		)

		if dir[0] == '~' {

			p, err = homedir.Expand(dir)
			if err != nil {
				return nil, err
			}

		} else {
			p, err = filepath.Abs(dir)
			if err != nil {
				return nil, err
			}
		}
		dir = p
	}

	return &Packs{dir: dir}, nil
}

// Dir returns the packs data directory
func (packs *Packs) Dir() string {
	return packs.dir
}

// Load downloads packs from the given git remoteURL
func (packs *Packs) Load(remoteURL string) error {
	if utils.FileExists(packs.dir) {
		return errPackDirExists
	}

	u, err := url.Parse(remoteURL)
	if err != nil {
		return err
	}

	switch u.Scheme {
	case "http", "https", "ssh":
		_, err = git.PlainClone(packs.dir, false, &git.CloneOptions{
			URL:      remoteURL,
			Progress: os.Stdout,
		})

	default:
		return errPackSourceNotSupported

	}

	return err
}

// Update performs a git pull on the packs repo to get the latest updates.  This
// should be called after a successful call to load
func (packs *Packs) Update() error {
	repo, err := git.PlainOpen(packs.dir)
	if err != nil {
		return err
	}

	wt, err := repo.Worktree()
	if err != nil {
		return err
	}

	opt := &git.PullOptions{Progress: os.Stdout}
	return wt.Pull(opt)
}

// Web returns a web packs manager
func (packs *Packs) Web() *BasePacks {
	return NewBasePacks(filepath.Join(packs.dir, webPackID))
}

// Dev returns a dev packs manager
func (packs *Packs) Dev() *DevPacks {
	return NewDevPacks(filepath.Join(packs.dir, devPackID))
}

// Datastore returns a datastore packs manager
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
