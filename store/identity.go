package store

import (
	"errors"
	"time"

	"github.com/euforia/thrap/thrapb"
)

// IdentityStore is the identity storage system
type IdentityStore struct {
	st ObjectStorage
}

// NewIdentityStore returns a new IdentityStore instance
func NewIdentityStore(objs ObjectStorage) *IdentityStore {
	return &IdentityStore{st: objs}
}

// Get returns an identity by the id
func (store *IdentityStore) Get(id string) (*thrapb.Identity, error) {
	ref, _, err := store.st.GetRef(id, "latest")
	if err != nil {
		return nil, err
	}

	var ident thrapb.Identity
	err = store.st.Get(id, ref.DataDigest, &ident)

	return &ident, err
}

// Create cretes a new identity
func (store *IdentityStore) Create(ident *thrapb.Identity) (*thrapb.Identity, error) {
	if ident.ID == "" {
		return nil, errIDMissing
	}

	prev, _, err := store.st.CreateRef(ident.ID, "latest")
	if err != nil {
		return nil, err
	}

	_, err = store.setIdent(ident, "latest", prev)
	return ident, err
}

// Update updates an identity with the ident provided
func (store *IdentityStore) Update(ident *thrapb.Identity) (*thrapb.Identity, error) {
	_, prev, err := store.st.GetRef(ident.ID, "latest")
	if err != nil {
		return nil, err
	}

	_, err = store.setIdent(ident, "latest", prev)
	return ident, err
}

// Delete deletes an identity by the id
func (store *IdentityStore) Delete(id string) error {
	return errors.New("to be implemented")
}

func (store *IdentityStore) setIdent(ident *thrapb.Identity, ref string, prev []byte) (*thrapb.ChainHeader, error) {
	data, err := store.st.Set(ident.ID, ident)
	if err != nil {
		return nil, err
	}

	refobj := &thrapb.ChainHeader{
		Previous:   prev,
		DataDigest: data,
		Timestamp:  time.Now().UnixNano(),
	}

	_, err = store.st.SetRef(ident.ID, ref, refobj)
	return refobj, err
}
