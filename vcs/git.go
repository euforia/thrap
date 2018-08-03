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
)

const (
	gitUserConfigFile = ".gitconfig"
	gitIgnoresFile    = ".gitignore"
	defaultRemoteName = "origin"
)

var (
	errPathNotSpecified = errors.New("path not specified")
)

// DefaultGitRemoteURL returns the default url scheme scheme for remote access
func DefaultGitRemoteURL(addr, owner, name string) string {
	// User is required
	return "ssh://git@" + addr + "/" + owner + "/" + name
}

// SetupLocalGitRepo initializes a new git repo.  It returns an error
// if it already has been initialized or fails
func SetupLocalGitRepo(projName, repoOwner, projPath, remoteAddr string) (VCS, *git.Repository, error) {
	vcsp := NewGitVCS()

	rr := &Repository{Name: projName}
	opt := Option{
		Path:   projPath,
		Remote: DefaultGitRemoteURL(remoteAddr, repoOwner, projName),
	}

	resp, _, err := vcsp.Create(rr, opt)
	if err != nil {
		return vcsp, nil, err
	}

	return vcsp, resp.(*git.Repository), nil

}

// GitVCS implements a git version controlled interface
type GitVCS struct {
	globalUser  string
	globalEmail string
}

// NewGitVCS returns a new Git vcs
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

// loadFromGitConfig loads settings from the git config file in the users home
// directory.  It does not return an error is the name and email are missing
func (g *GitVCS) loadFromGitConfig(fpath string) error {
	f, err := ini.Load(fpath)
	if err != nil {
		return err
	}

	sec := f.Section("user")

	nkey, err := sec.GetKey("name")
	if err == nil {
		g.globalUser = nkey.Value()
	}

	nkey, err = sec.GetKey("email")
	if err == nil {
		g.globalEmail = nkey.Value()
	}

	return nil
}

// GlobalUser satisfies the VCS interface
func (g *GitVCS) GlobalUser() string {
	return g.globalUser
}

// GlobalEmail satisfies the VCS interface
func (g *GitVCS) GlobalEmail() string {
	return g.globalEmail
}

// IgnoresFile satisfies the VCS interface
func (g *GitVCS) IgnoresFile() string {
	return gitIgnoresFile
}

// Create creates a new Repository. Each call is idempotent.  It returns
// the object, whether it was created or an error
func (g *GitVCS) Create(repo *Repository, opt Option) (interface{}, bool, error) {
	var (
		gitRepo *git.Repository
		err     error
	)

	if len(opt.Path) == 0 {
		return nil, false, errPathNotSpecified
	}

	gitRepo, err = git.PlainInit(opt.Path, false)

	// We still check remote
	if err != nil {
		if err != git.ErrRepositoryAlreadyExists {
			return nil, false, err
		}

		// Open so we can set remote
		gitRepo, err = git.PlainOpen(opt.Path)
		if err != nil {
			return nil, false, err
		}

	}

	err = g.setupRemote(gitRepo, opt.Remote)

	return gitRepo, true, err
}

// Open opens a local repo
func (g *GitVCS) Open(repo *Repository, opt Option) (interface{}, error) {
	if len(opt.Path) == 0 {
		return nil, errPathNotSpecified
	}
	return git.PlainOpen(opt.Path)
}

// Status returns the repo status
func (g *GitVCS) Status(opt Option) (git.Status, error) {
	repo, err := git.PlainOpen(opt.Path)
	if err != nil {
		return git.Status{}, err
	}

	wt, err := repo.Worktree()
	if err != nil {
		return git.Status{}, err
	}

	status, err := wt.Status()
	if err != nil {
		return git.Status{}, err
	}

	// We can only publish if everything has been committed
	return status, nil
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
