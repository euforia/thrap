package core

import "github.com/euforia/thrap/thrapb"

// IdentityStorage is a identity storage interface
type IdentityStorage interface {
	// Get returns an identity be the given id
	Get(id string) (*thrapb.Identity, error)
	Create(*thrapb.Identity) (*thrapb.Identity, error)
	Update(*thrapb.Identity) (*thrapb.Identity, error)
	Iter(string, func(*thrapb.Identity) error) error
}

// StackStorage is a stack storage interface
type StackStorage interface {
	Get(string) (*thrapb.Stack, error)
	Create(*thrapb.Stack) (*thrapb.Stack, error)
	Update(*thrapb.Stack) (*thrapb.Stack, error)
	Iter(string, func(*thrapb.Stack) error) error
}
