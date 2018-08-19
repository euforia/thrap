package storage

import (
	"github.com/euforia/kvdb"

	"github.com/euforia/thrap/thrapb"
)

// KvdbIdentityStorage implements a badger backed IdentityStorage interface
type KvdbIdentityStorage struct {
	table kvdb.Table
}

// Get returns an identity by the id
func (store *KvdbIdentityStorage) Get(id string) (*thrapb.Identity, error) {
	obj, err := store.table.Get([]byte(id))
	if err == nil {
		if stack, ok := obj.(*thrapb.Identity); ok {
			return stack, nil
		}
		err = ErrInvalidType
	}
	return nil, err
}

// Create creates a new identity. It returns an error if it exists
func (store *KvdbIdentityStorage) Create(ident *thrapb.Identity) error {
	return store.table.Create([]byte(ident.ID), ident)
}

// Update an identity.  It returns an ErrIdentityNotFound
func (store *KvdbIdentityStorage) Update(ident *thrapb.Identity) error {
	return store.table.Update([]byte(ident.ID), ident)
}

// Iter iterates over each identity from the starting point
func (store *KvdbIdentityStorage) Iter(start string, callback func(*thrapb.Identity) error) error {
	var prefix []byte
	if start != "" {
		prefix = []byte(start)
	}
	return store.table.Iter(prefix, func(obj kvdb.Object) error {
		proj, ok := obj.(*thrapb.Identity)
		if ok {
			return callback(proj)
		}
		return ErrInvalidType
	})
}

// Delete deletes an identity by the id. It returns an ErrINotFound
func (store *KvdbIdentityStorage) Delete(id string) error {
	return store.table.Delete([]byte(id))
}
