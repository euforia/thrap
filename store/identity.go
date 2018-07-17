package store

import (
	"errors"
	"time"

	"github.com/dgraph-io/badger"
	"github.com/gogo/protobuf/proto"

	"github.com/euforia/thrap/thrapb"
)

const defaultIdentityPrefix = "/identity/"

var (
	ErrIdentityNotFound = errors.New("identity not found")
	ErrIdentityExists   = errors.New("identity exists")
)

type BadgerIdentityStorage struct {
	db *badger.DB
}

func NewBadgerIdentityStorage(db *badger.DB) *BadgerIdentityStorage {
	return &BadgerIdentityStorage{db: db}
}

func (store *BadgerIdentityStorage) getOpaqueKey(k string) []byte {
	return []byte(defaultIdentityPrefix + k)
}

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

	if ident == nil {
		ident = &thrapb.Identity{}
	}

	return ident, err
}

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

// Iter iterates over each identity in the store
func (store *IdentityStore) Iter(start string, callback func(*thrapb.Identity) error) error {
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
