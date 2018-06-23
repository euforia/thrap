package vcs

import (
	"errors"
	"fmt"
)

var (
	errNotImplemented = errors.New("Not Implemented")
)

// Repository is a code repository with common fields
type Repository struct {
	Name        string
	Owner       string
	Description string
	Private     bool
}

// Config hold a VCS config
type Config struct {
	// Service providing the vcs functionality
	Provider string
	Conf     map[string]interface{}
}

// Option holds VCS call options
type Option struct {
	// Local path
	Path string
	// Remote path
	Remote string
}

// VCS implements a version control system interface such as git, svn etc.
type VCS interface {
	// Initialize the VCS interface
	Init(conf map[string]interface{}) error
	// Create version control repository
	Create(*Repository, Option) (interface{}, error)
	// Deletes a repository
	Delete(*Repository, Option) error
	// Name of ignores file used by this vcs
	IgnoresFile() string
	// Returns the default configured user.  This is used for commit info
	GlobalUser() string
	// Returns the default configured email. This is used for commit info
	GlobalEmail() string
	// Unique id for provider
	ID() string
}

// New returns a new VCS interface based on the given config
func New(conf *Config) (VCS, error) {
	var (
		v   VCS
		err error
	)

	switch conf.Provider {
	case "git":
		v = NewGitVCS()

	case "github":
		v = newGithubVCS(nil)

	default:
		err = fmt.Errorf("unsupported vcs: '%s'", conf.Provider)

	}

	if err == nil {
		err = v.Init(conf.Conf)
	}

	return v, err
}
