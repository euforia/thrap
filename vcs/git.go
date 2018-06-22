package vcs

import (
	"errors"
	"fmt"
	"net/url"
	"os"
	"path/filepath"

	"github.com/euforia/thrap/utils"
	"github.com/go-ini/ini"
	homedir "github.com/mitchellh/go-homedir"
	"gopkg.in/src-d/go-git.v4"
	"gopkg.in/src-d/go-git.v4/config"
	"gopkg.in/src-d/go-git.v4/plumbing"
	"gopkg.in/src-d/go-git.v4/plumbing/object"
)

const (
	gitUserConfigFile = ".gitconfig"
	gitIgnoresFile    = ".gitignore"
	defaultRemoteName = "origin"
)

// RepoVersion is the calculate repo / project / stack version
type RepoVersion struct {
	// Last known tag.  Defaults to v0.0.0
	Tag string
	// Count from the tag
	Count int
	// Last commit hash
	Hash plumbing.Hash
}

// GetRepoVersion returns the latest tag, the count from that tag and the hash
// of the last commit
func GetRepoVersion(path string) *RepoVersion {
	tagLabel := "v0.0.0"

	repo, _ := git.PlainOpen(path)

	tags, _ := repo.Tags()

	var rc *object.Commit
	lastTag, _ := tags.Next()
	if lastTag != nil {
		tagLabel = lastTag.Name().Short()
		to, _ := repo.TagObject(lastTag.Hash())
		if to != nil {
			rc, _ = to.Commit()
		}
	}

	cmtIter, _ := repo.CommitObjects()

	var c int

	if rc == nil {

		cmtIter.ForEach(func(arg1 *object.Commit) error {
			c++
			return nil
		})

	} else {
		// Calculate the total and mark where the tag was found
		var t int
		cmtIter.ForEach(func(arg1 *object.Commit) error {
			t++
			if arg1.ID().String() == rc.ID().String() {
				c = t
			}

			return nil
		})
		// Subtract tag count from total to get the count from the tag
		c = t - c
	}

	head, _ := repo.Head()

	return &RepoVersion{tagLabel, c, head.Hash()}
}

// DefaultGitRemoteURL returns the default url scheme scheme for remote access
func DefaultGitRemoteURL(addr, owner, name string) string {
	return "ssh://git@" + addr + "/" + owner + "/" + name
}

// GitVCS implements a git version controlled interface
type GitVCS struct {
	defaultUser  string
	defaultEmail string
}

func NewGitVCS() *GitVCS {
	return &GitVCS{}
}

// ID returns the vcs provider id
func (g *GitVCS) ID() string {
	return "git"
}

// Init initializes the git provider with the config info in the users home dir
func (g *GitVCS) Init(conf map[string]interface{}) error {
	fpath, err := homedir.Expand("~/" + gitUserConfigFile)
	if err != nil {
		return err
	}

	if utils.FileExists(fpath) {
		err = g.loadFromGitConfig(fpath)
	}

	return err
}

func (g *GitVCS) loadFromGitConfig(fpath string) error {
	f, err := ini.Load(fpath)
	if err != nil {
		return err
	}

	sec := f.Section("user")

	nkey, err := sec.GetKey("name")
	if err != nil {
		return err
	}
	g.defaultUser = nkey.Value()

	nkey, err = sec.GetKey("email")
	if err == nil {
		g.defaultEmail = nkey.Value()
	}

	return err
}

// DefaultUser satisfies the VCS interface
func (g *GitVCS) DefaultUser() string {
	return g.defaultUser
}

// DefaultEmail satisfies the VCS interface
func (g *GitVCS) DefaultEmail() string {
	return g.defaultEmail
}

// IgnoresFile satisfies the VCS interface
func (g *GitVCS) IgnoresFile() string {
	return gitIgnoresFile
}

// Create creates a new Repository returning a repo of the vcs' type. Each call
// only fills in missing pieces so multiple calls will not corrupt
func (g *GitVCS) Create(repo *Repository, opt Option) (interface{}, error) {
	var (
		gitRepo *git.Repository
		err     error
	)

	if len(opt.Path) == 0 {
		return nil, errors.New("path not specified")
	}

	gitRepo, err = git.PlainInit(opt.Path, false)

	// We still check remote
	if err != nil {
		if err != git.ErrRepositoryAlreadyExists {
			return nil, err
		}

		// Open so we can set remote
		gitRepo, err = git.PlainOpen(opt.Path)
		if err != nil {
			return nil, err
		}

	}

	err = g.setupRemote(gitRepo, opt.Remote)

	return gitRepo, err
}

// Delete removes the .git directory leaving all other files intact
func (g *GitVCS) Delete(repo *Repository, opt Option) error {
	abspath := filepath.Join(opt.Path, ".git")
	return os.RemoveAll(abspath)
}

// setup remote configuration on local instance
func (g *GitVCS) setupRemote(repo *git.Repository, remote string) (err error) {
	if len(remote) == 0 {
		return
	}

	_, err = parseURL(remote)
	if err != nil {
		return
	}

	rmt, err := repo.Remote(defaultRemoteName)
	if err == nil {
		conf := rmt.Config()
		if len(conf.URLs) > 0 {
			// Its has been externally set to we don't touch
			return
		}
		conf.URLs = []string{remote}
	}

	_, err = repo.CreateRemote(&config.RemoteConfig{
		Name: defaultRemoteName,
		URLs: []string{remote},
	})

	return err
}

func parseURL(uri string) (*url.URL, error) {
	u, err := url.Parse(uri)
	if err != nil {
		return nil, err
	}

	switch u.Scheme {
	case "http", "https", "ssh":
	default:
		err = fmt.Errorf("unsupported URL scheme: '%s'", u.Scheme)
	}

	return u, err
}
