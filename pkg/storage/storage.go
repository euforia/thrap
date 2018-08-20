package storage

import (
	"errors"

	"github.com/euforia/thrap/thrapb"
)

var (
	// ErrInvalidType is used when there is a type mismatch
	ErrInvalidType = errors.New("invalid type")
)

// IdentityStorage implements a storage interface for Identity
type IdentityStorage interface {
	Get(id string) (*thrapb.Identity, error)
	Update(proj *thrapb.Identity) error
	Create(proj *thrapb.Identity) error
	Delete(id string) error
	Iter(start string, cb func(*thrapb.Identity) error) error
}

// ProfileStorage implements a storage interface for profiles
type ProfileStorage interface {
	Get(id string) (*thrapb.Profile, error)
	List() []*thrapb.Profile
	GetDefault() *thrapb.Profile
	SetDefault(id string) error
}
