package store

import (
	"errors"

	"github.com/dgraph-io/badger"
	"github.com/gogo/protobuf/proto"

	"github.com/euforia/thrap/thrapb"
)

const defaultIdentityPrefix = "/identity/"

var (
	// ErrIdentityNotFound is used when an identity is not found in the store
	ErrIdentityNotFound = errors.New("identity not found")
	// ErrIdentityExists is used when an identity exists
	ErrIdentityExists = errors.New("identity exists")
)

// BadgerIdentityStorage implements a badger backed IdentityStorage interface
type BadgerIdentityStorage struct {
	db *badger.DB
}

// NewBadgerIdentityStorage returns a new BadgerIdentityStorage
func NewBadgerIdentityStorage(db *badger.DB) *BadgerIdentityStorage {
	return &BadgerIdentityStorage{db: db}
}

func (store *BadgerIdentityStorage) getOpaqueKey(k string) []byte {
	return []byte(defaultIdentityPrefix + k)
}

// Get returns an identity by the id
func (store *BadgerIdentityStorage) Get(id string) (*thrapb.Identity, error) {
	var (
		key   = store.getOpaqueKey(id)
		ident *thrapb.Identity
	)

	err := store.db.View(func(txn *badger.Txn) error {
		var err error
		ident, err = store.getIdentity(txn, key)
		return err
	})

	return ident, err
}

// Create creates a new identity. It returns an error if it exists
func (store *BadgerIdentityStorage) Create(ident *thrapb.Identity) (*thrapb.Identity, error) {
	key := store.getOpaqueKey(ident.ID)
	val, err := proto.Marshal(ident)
	if err != nil {
		return nil, err
	}

	err = store.db.Update(func(txn *badger.Txn) error {
		_, err := txn.Get(key)
		if err == nil {
			return ErrIdentityExists
		}

		return txn.Set(key, val)
	})
	return ident, err
}

// Update an identity.  It returns an ErrIdentityNotFound
func (store *BadgerIdentityStorage) Update(ident *thrapb.Identity) (*thrapb.Identity, error) {
	key := store.getOpaqueKey(ident.ID)
	val, err := proto.Marshal(ident)
	if err != nil {
		return nil, err
	}

	err = store.db.Update(func(txn *badger.Txn) error {
		_, err := txn.Get(key)
		if err != nil {
			return ErrIdentityNotFound
		}

		return txn.Set(key, val)
	})

	return ident, err
}

// Iter iterates over each identity from the starting point
func (store *BadgerIdentityStorage) Iter(start string, callback func(*thrapb.Identity) error) error {
	prefix := store.getOpaqueKey(start)

	return store.db.View(func(txn *badger.Txn) error {
		iter := txn.NewIterator(badger.DefaultIteratorOptions)
		defer iter.Close()

		for iter.Seek(prefix); iter.ValidForPrefix(prefix); iter.Next() {
			item := iter.Item()
			ident, err := identityFromItem(item)
			if err != nil {
				return err
			}
			if err = callback(ident); err != nil {
				return err
			}
		}

		return nil
	})
}

// Delete deletes an identity by the id. It returns an ErrIdentityNotFound
func (store *BadgerIdentityStorage) Delete(id string) (*thrapb.Identity, error) {
	var (
		key   = store.getOpaqueKey(id)
		ident *thrapb.Identity
	)

	err := store.db.Update(func(txn *badger.Txn) error {
		var err error
		ident, err = store.getIdentity(txn, key)
		if err != nil {
			return ErrIdentityNotFound
		}

		return txn.Delete(key)
	})

	return ident, err
}

func (store *BadgerIdentityStorage) getIdentity(txn *badger.Txn, key []byte) (*thrapb.Identity, error) {
	item, err := txn.Get(key)
	if err != nil {
		return nil, err
	}

	return identityFromItem(item)
}

func identityFromItem(item *badger.Item) (*thrapb.Identity, error) {
	val, err := item.Value()
	if err != nil {
		return nil, err
	}

	var ident thrapb.Identity
	err = proto.Unmarshal(val, &ident)

	return &ident, err
}
